// Package qzone qq空间发说说
package qzone

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/guohuiyuan/qzone"
	"github.com/jinzhu/gorm"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() { // 插件主体
	engine := control.Register("qzone", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "QQ空间表白墙\n" +
			"- 登录QQ空间\n" +
			"- 发表白墙[xxx]\n" +
			"- [同意|拒绝]说说 [说说ID]\n" +
			"- 查看[等待|同意|拒绝]说说\n" +
			"- 查看说说 [说说ID]",
		PrivateDataFolder: "qzone",
	})
	go func() {
		qdb = initialize(engine.DataFolder() + "qzone.db")
	}()
	engine.OnFullMatch("登录QQ空间", zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var (
				uin     string
				cookies string
			)

			gCurCookieJar, _ := cookiejar.New(nil)
			client := &http.Client{
				Jar: gCurCookieJar,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			ptqrcodeReq, err := http.NewRequest("GET", qrcodeURL, nil)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			qrcodeResp, err := client.Do(ptqrcodeReq)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			defer qrcodeResp.Body.Close()
			var qrsig string
			for _, v := range qrcodeResp.Cookies() {
				if v.Name == "qrsig" {
					qrsig = v.Value
					break
				}
			}
			if qrsig == "" {
				ctx.SendChain(message.Text("ERROR: qrsig为空"))
				return
			}
			data, err := io.ReadAll(qrcodeResp.Body)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Text("请扫描二维码, 登录qq空间"))
			ctx.SendChain(message.ImageBytes(data))
			qrtoken := getPtqrtoken(qrsig)
			for {
				time.Sleep(2 * time.Second)
				checkReq, err := http.NewRequest("GET", fmt.Sprintf(loginCheckURL, qrtoken), nil)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				checkResp, err := client.Do(checkReq)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				defer checkResp.Body.Close()
				checkData, err := io.ReadAll(checkResp.Body)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				checkText := binary.BytesToString(checkData)
				switch {
				case strings.Contains(checkText, "二维码已失效"):
					ctx.SendChain(message.Text("二维码已失效, 登录失败"))
					return
				case strings.Contains(checkText, "登录成功"):
					dealedCheckText := strings.ReplaceAll(checkText, "'", "")
					redirectURL := strings.Split(dealedCheckText, ",")[2]
					u, err := url.Parse(redirectURL)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						return
					}
					values, err := url.ParseQuery(u.RawQuery)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						return
					}
					ptsigx := values["ptsigx"][0]
					uin = values["uin"][0]
					redirectReq, err := http.NewRequest("GET", fmt.Sprintf(checkSigURL, uin, ptsigx), nil)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						return
					}
					redirectResp, err := client.Do(redirectReq)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						return
					}
					defer redirectResp.Body.Close()
					for _, v := range redirectResp.Cookies() {
						if v.Value != "" {
							cookies += v.Name + "=" + v.Value + ";"
						}
					}
					qq, err := strconv.Atoi(uin)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						return
					}
					err = qdb.insertOrUpdate(int64(qq), cookies)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						return
					}
					ctx.SendChain(message.Text("登录成功"))
					return
				}
			}
		})
	engine.OnRegex(`^发表白墙.*?([\s\S]*)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			regexMatched := ctx.State["regex_matched"].([]string)
			qq := ctx.Event.UserID
			e := emotion{
				QQ:     qq,
				Msg:    message.UnescapeCQCodeText(regexMatched[1]),
				Status: waitStatus,
				Tag:    loveTag,
			}
			if zero.SuperUserPermission(ctx) {
				text, base64imgs, err := parseTextAndImg(e.Msg)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				err = publishEmotion(ctx.Event.SelfID, text, base64imgs)
				if err != nil {
					if gorm.IsRecordNotFoundError(err) {
						ctx.SendChain(message.Text(zero.BotConfig.NickName[0], "(", ctx.Event.SelfID, ")", "未登录QQ空间,请发送\"登录QQ空间\"初始化配置"))
						return
					}
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				ctx.SendChain(message.Text("发表成功"))
				return
			}
			_, err := qdb.saveEmotion(e)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Text("已收稿, 请耐心等待审核"))
		})
	engine.OnRegex(`^(同意|拒绝)说说\s?(\d{1,10})$`, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			regexMatched := ctx.State["regex_matched"].([]string)
			id, err := strconv.Atoi(regexMatched[2])
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			switch regexMatched[1] {
			case "同意":
				err = qdb.updateEmotionStatusByID(int64(id), agreeStatus)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				err = getAndPublishEmotion(ctx.Event.SelfID, int64(id))
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				ctx.SendChain(message.Text("同意说说", id, ", 发表成功"))
			case "拒绝":
				err = qdb.updateEmotionStatusByID(int64(id), disagreeStatus)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				ctx.SendChain(message.Text("拒绝说说", id))
			}
		})
	engine.OnRegex(`^查看(.{0,2})说说\s?(\d{0,10})$`, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			regexMatched := ctx.State["regex_matched"].([]string)
			id, _ := strconv.Atoi(regexMatched[2])
			if id != 0 {
				e, err := qdb.getEmotionByID(int64(id))
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				m := message.Message{ctxext.FakeSenderForwardNode(ctx, message.ParseMessageFromString(e.Msg)...)}
				if id := ctx.Send(m).ID(); id == 0 {
					ctx.SendChain(message.Text("ERROR:  可能被风控或下载图片用时过长，请耐心等待"))
				}
				return
			}
			var status int
			switch regexMatched[1] {
			case "等待":
				status = 1
			case "同意":
				status = 2
			case "拒绝":
				status = 3
			}
			el, err := qdb.getEmotionByStatus(status)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			t := "所有" + regexMatched[1] + "说说\n可发送\"查看说说[序号]\"查看详情\n"
			for _, v := range el {
				t += fmt.Sprintf("%v. %v: %v\n\n", v.ID, ctx.CardOrNickName(v.QQ), simpleMsg(v.Msg))
			}
			messageBytes, err := text.RenderToBase64(t, text.FontFile, 400, 20)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			if id := ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+binary.BytesToString(messageBytes))); id.ID() == 0 {
				ctx.SendChain(message.Text("ERROR: 可能被风控了"))
			}
		})
}

func getAndPublishEmotion(qq int64, id int64) (err error) {
	e, err := qdb.getEmotionByID(id)
	if err != nil {
		return
	}
	text, base64imgs, err := parseTextAndImg(e.Msg)
	if err != nil {
		return
	}
	return publishEmotion(qq, text, base64imgs)
}

func publishEmotion(qq int64, text string, base64imgs []string) (err error) {
	qc, err := qdb.getByUin(qq)
	if err != nil {
		return
	}
	m := qzone.NewManager(qc.Cookie)
	_ = m.RefreshToken()
	_, err = m.SendShuoShuoWithBase64Pic(text, base64imgs)
	return
}

func parseTextAndImg(raw string) (text string, base64imgs []string, err error) {
	base64imgs = make([]string, 0, 16)
	var imgdata []byte
	m := message.ParseMessageFromString(raw)
	for _, v := range m {
		if v.Type == "text" && v.Data["text"] != "" {
			text += v.Data["text"] + "\n"
		}
		if v.Type == "image" && v.Data["url"] != "" {
			imgdata, err = web.GetData(v.Data["url"])
			if err != nil {
				return
			}
			encodeStr := base64.StdEncoding.EncodeToString(imgdata)
			base64imgs = append(base64imgs, encodeStr)
		}
	}
	return
}

func simpleMsg(raw string) (simple string) {
	m := message.ParseMessageFromString(raw)
	for _, v := range m {
		if v.Type == "text" && v.Data["text"] != "" {
			simple += v.Data["text"]
		} else {
			simple += "[CQ:" + v.Type + "]"
		}
	}
	return
}

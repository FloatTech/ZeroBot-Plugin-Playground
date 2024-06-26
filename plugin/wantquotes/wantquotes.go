// Package wantquotes 据意查句
package wantquotes

import (
	"encoding/json"
	"fmt"
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
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	wantquotesURL = "https://wantquotes.net"
	getQrcode     = "/api/get_qrcode/"
	showqrcode    = "https://mp.weixin.qq.com/cgi-bin/showqrcode"
	loginCheck    = "/api/login_check/"
	semantic      = "/api/semantic/"
)

var (
	typeList = []string{"现", "现-名言", "现-佳句", "现-佳句-文学", "现-佳句-诗歌", "现-佳句-其他", "现-网络", "现-台词", "现-台词-影视剧", "现-台词-动漫", "现-台词-综艺",
		"古", "谚", "谚-谚语", "谚-俗语", "谚-惯用语", "歇"}
)

type getQrcodeRsp struct {
	Ticket  string `json:"ticket"`
	SceneID string `json:"scene_id"`
}

type loginCheckRsp struct {
	Login   int    `json:"login"`
	Secret  string `json:"secret"`
	Unionid string `json:"unionid"`
}

type semanticRsp struct {
	Quotes []Quotes `json:"quotes"`
	TopSim float64  `json:"top_sim"`
}

// Quotes 名句结构体
type Quotes struct {
	Quote      string `json:"quote"`
	SourceType string `json:"source_type"`
	Author     string `json:"author"`
	Work       string `json:"work"`
	ID         int    `json:"id"`
}

func init() {
	engine := control.Register("wantquotes", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "据意查句",
		Help: "- 据意查句 大海 (每次重启需, 登录据意查句)\n" +
			"- 登录据意查句",
		PrivateDataFolder: "wantquotes",
		Extra:             -1,
	})

	// 开启
	engine.OnRegex(`^据意查句\s?(.{1,25})$`, getPara).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		keyword := ctx.State["regex_matched"].([]string)[1]
		quotesType := ctx.State["quotesType"].(string)
		key := getAPIKey(ctx)
		unionid, secret, _ := strings.Cut(key, "|")
		apiURL := fmt.Sprintf("%s?query=%s&type=%s&unionid=%s&secret=%s", wantquotesURL+semantic, url.QueryEscape(keyword), url.QueryEscape(quotesType), url.QueryEscape(unionid), url.QueryEscape(secret))
		data, err := web.RequestDataWith(web.NewDefaultClient(), apiURL, "GET", wantquotesURL, web.RandUA(), nil)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		var rsp semanticRsp
		err = json.Unmarshal(data, &rsp)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		m := make(message.Message, 0, len(rsp.Quotes))
		text := strings.Builder{}
		for _, v := range rsp.Quotes {
			text.WriteString(v.Quote)
			text.WriteString("\n—— ")
			text.WriteString(v.Author)
			text.WriteString(" 《")
			text.WriteString(v.Work)
			text.WriteString("》")
			m = append(m, ctxext.FakeSenderForwardNode(ctx, message.Text(text.String())))
			text.Reset()
		}
		if id := ctx.Send(m).ID(); id == 0 {
			ctx.SendChain(message.Text("ERROR: 可能被风控或下载图片用时过长，请耐心等待"))
		}
	})
	engine.OnFullMatch(`登录据意查句`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getQrcodeData, err := web.RequestDataWith(web.NewDefaultClient(), wantquotesURL+getQrcode, "GET", wantquotesURL, web.RandUA(), nil)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		var qrRsp getQrcodeRsp
		err = json.Unmarshal(getQrcodeData, &qrRsp)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		showQrcodeData, err := web.RequestDataWith(web.NewDefaultClient(), showqrcode+"?ticket="+qrRsp.Ticket, "GET", "", web.RandUA(), nil)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(message.Text("WantQuotes\n微信扫码登录\n首次登录需关注公众号"))
		ctx.SendChain(message.ImageBytes(showQrcodeData))
		now := time.Now()
		for {
			if time.Since(now).Seconds() > 120 {
				ctx.SendChain(message.Text("据意查句登录超时,请重新登录"))
				return
			}
			time.Sleep(2 * time.Second)
			loginCheckData, err := web.RequestDataWith(web.NewDefaultClient(), wantquotesURL+loginCheck+"?scene_id="+qrRsp.SceneID, "GET", "", web.RandUA(), nil)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			var lcr loginCheckRsp
			err = json.Unmarshal(loginCheckData, &lcr)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			if lcr.Login == 1 {
				err := setAPIKey(ctx.State["manager"].(*ctrl.Control[*zero.Ctx]), lcr.Unionid+"|"+lcr.Secret)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				ctx.SendChain(message.Text("据意查句登录成功"))
				return
			}
		}
	})
}

func getPara(ctx *zero.Ctx) bool {
	next := zero.NewFutureEvent("message", 999, false, ctx.CheckSession())
	recv, cancel := next.Repeat()
	defer cancel()
	tex := strings.Builder{}
	tex.WriteString("请下列选择查询名句的类型\n")
	for i, v := range typeList {
		tex.WriteString(strconv.Itoa(i))
		tex.WriteString(". ")
		tex.WriteString(v)
		tex.WriteString("\n")
	}
	base64Str, err := text.RenderToBase64(tex.String(), text.FontFile, 400, 20)
	if err != nil {
		ctx.SendChain(message.Text("图片生成错误了, ", zero.BotConfig.NickName[0], "帮你选择查询名句类型"))
		ctx.State["quotesType"] = typeList[0]
		return true
	}
	ctx.SendChain(message.Image("base64://" + binary.BytesToString(base64Str)))
	for {
		select {
		case <-time.After(time.Second * 10):
			ctx.SendChain(message.Text("时间太久啦！", zero.BotConfig.NickName[0], "帮你选择查询名句类型"))
			ctx.State["quotesType"] = typeList[0]
			return true
		case c := <-recv:
			msg := c.Event.Message.ExtractPlainText()
			num, _ := strconv.Atoi(msg)
			if num < 0 || num >= len(typeList) {
				ctx.SendChain(message.Text("序号非法!"))
				continue
			}
			ctx.State["quotesType"] = typeList[num]
			return true
		}
	}
}

// setAPIKey 获取apikey
func getAPIKey(ctx *zero.Ctx) (apikey string) {
	m := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
	_ = m.GetExtra(&apikey)
	logrus.Debugln("[wantquotes] get api key:", apikey)
	return
}

// setAPIKey 设置apikey, 格式为unionid|secret
func setAPIKey(m *ctrl.Control[*zero.Ctx], apikey string) error {
	_ = m.Manager.Response(-1)
	return m.SetExtra(apikey)
}

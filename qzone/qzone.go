// Package qzone qq空间发说说
package qzone

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/binary"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() { // 插件主体
	engine := control.Register("qzone", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "qq空间发说说\n" +
			"- 登录qq空间\n" +
			"- 发说说 xxx",
		PrivateDataFolder: "qzone",
	})
	go func() {
		qdb = initialize(engine.DataFolder() + "qzone.db")
	}()
	engine.OnFullMatch("登录qq空间", zero.SuperUserPermission, zero.OnlyPrivate).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var (
				uin     string
				skey    string
				pskey   string
				cookies string
			)
			for i := 0; i < attempts; i++ {
				gCurCookieJar, _ := cookiejar.New(nil)
				client := &http.Client{
					Jar: gCurCookieJar,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				}
				ptqrcodeReq, err := http.NewRequest("GET", qrcodeURL, nil)
				if err != nil {
					continue
				}
				qrcodeResp, err := client.Do(ptqrcodeReq)
				if err != nil {
					continue
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
					continue
				}
				data, err := io.ReadAll(qrcodeResp.Body)
				if err != nil {
					continue
				}
				ctx.SendChain(message.Text("请扫描二维码, 登录qq空间"))
				ctx.SendChain(message.ImageBytes(data))
				qrtoken := getPtqrtoken(qrsig)
			LOOP:
				for {
					time.Sleep(2 * time.Second)
					checkReq, err := http.NewRequest("GET", fmt.Sprintf(loginCheckURL, qrtoken), nil)
					if err != nil {
						continue
					}
					checkResp, err := client.Do(checkReq)
					if err != nil {
						continue
					}
					defer checkResp.Body.Close()
					checkData, err := io.ReadAll(checkResp.Body)
					if err != nil {
						continue
					}
					checkText := binary.BytesToString(checkData)
					switch {
					case strings.Contains(checkText, "二维码已失效"):
						ctx.SendChain(message.Text("二维码已失效"))
						break LOOP
					case strings.Contains(checkText, "登录成功"):
						dealedCheckText := strings.ReplaceAll(checkText, "'", "")
						redirectURL := strings.Split(dealedCheckText, ",")[2]
						u, err := url.Parse(redirectURL)
						if err != nil {
							break LOOP
						}
						values, err := url.ParseQuery(u.RawQuery)
						if err != nil {
							break LOOP
						}
						ptsigx := values["ptsigx"][0]
						uin = values["uin"][0]
						redirectReq, err := http.NewRequest("GET", fmt.Sprintf(checkSigURL, uin, ptsigx), nil)
						if err != nil {
							break LOOP
						}
						redirectResp, err := client.Do(redirectReq)
						if err != nil {
							break LOOP
						}
						defer redirectResp.Body.Close()
						for _, v := range redirectResp.Cookies() {
							if v.Name == "skey" && skey == "" {
								skey = v.Value
							}
							if v.Name == "p_skey" && pskey == "" {
								pskey = v.Value
							}
							if v.Value != "" {
								cookies += v.Name + "=" + v.Value + ";"
							}
						}
						qq, err := strconv.Atoi(uin)
						if err != nil {
							break LOOP
						}
						err = qdb.insertOrUpdate(int64(qq), skey, pskey, cookies)
						if err != nil {
							break LOOP
						}
						ctx.SendChain(message.Text("登录成功"))
						return
					}
				}
			}
			ctx.SendChain(message.Text("登录失败"))
		})
	engine.OnRegex(`^发说说\s+(.*)$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			qc, err := qdb.getByUin(ctx.Event.SelfID)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			m := newManager(qc.Uin, qc.Skey, qc.Pskey, qc.Cookies)
			err = m.RefreshToken()
			for i := 0; i <= refreshTimes && err != nil; i++ {
				if i == refreshTimes {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				err = m.RefreshToken()
			}
			msg := ctx.State["regex_matched"].([]string)
			_, err = m.SendShuoShuo(msg[1])
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Text("发送成功"))
		})
}

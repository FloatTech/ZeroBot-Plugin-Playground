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
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/OPQBot/qzone"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	m qzone.Manager
)

func init() { // 插件主体
	engine := control.Register("qzone", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "qq空间发说说\n" +
			"- 登录qq空间\n" +
			"- 发说说 xxx",
		PrivateDataFolder: "qzone",
	})
	engine.OnFullMatch("登录qq空间", zero.SuperUserPermission, zero.OnlyPrivate).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			for i := 0; i < attempts; i++ {
				gCurCookieJar, _ := cookiejar.New(nil)
				client := &http.Client{
					Jar: gCurCookieJar,
				}
				ptqrcodeReq, err := http.NewRequest("GET", qrcodeURL, nil)
				if err != nil {
					logrus.Errorln("ERROR: ", err)
					continue
				}
				qrcodeResp, err := client.Do(ptqrcodeReq)
				if err != nil {
					logrus.Errorln("ERROR: ", err)
					continue
				}
				defer qrcodeResp.Body.Close()
				qrsig := ""
				ptqrcodeCookie := qrcodeResp.Header.Get("Set-Cookie")
				for _, v := range strings.Split(ptqrcodeCookie, ";") {
					if strings.HasPrefix(v, "qrsig=") {
						qrsig = strings.TrimPrefix(v, "qrsig=")
						break
					}
				}
				if qrsig == "" {
					logrus.Errorln("ERROR: qrsig is emtpy")
					continue
				}
				data, err := io.ReadAll(qrcodeResp.Body)
				if err != nil {
					logrus.Errorln("ERROR: ", err)
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
						logrus.Errorln("ERROR: ", err)
						continue
					}
					checkResp, err := client.Do(checkReq)
					if err != nil {
						logrus.Errorln("ERROR: ", err)
						continue
					}
					defer checkResp.Body.Close()
					checkData, err := io.ReadAll(checkResp.Body)
					if err != nil {
						logrus.Errorln("ERROR: ", err)
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
							logrus.Errorln("ERROR: ", err)
							break LOOP
						}
						values, err := url.ParseQuery(u.RawQuery)
						if err != nil {
							logrus.Errorln("ERROR: ", err)
							break LOOP
						}
						ptsigx := values["ptsigx"][0]
						uin := values["uin"][0]
						redirectReq, err := http.NewRequest("GET", fmt.Sprintf(checkSigURL, uin, ptsigx), nil)
						if err != nil {
							logrus.Errorln("ERROR: ", err)
							break LOOP
						}
						redirectResp, err := client.Do(redirectReq)
						if err != nil {
							logrus.Errorln("ERROR: ", err)
							break LOOP
						}
						defer redirectResp.Body.Close()
						oc := OPQBot.Cookie{}
						finalCookie := redirectResp.Header.Values("Set-Cookie")
						oc.Cookies = strings.Join(finalCookie, ";")
						for _, v := range finalCookie {
							for _, c := range strings.Split(v, ";") {
								l, b, f := strings.Cut(c, "=")
								if !f {
									logrus.Errorln("ERROR: cut ", c)
									continue
								}
								if l == "skey" && oc.Skey == "" {
									oc.Skey = b
								}
								if l == "p_skey" && oc.PSkey.Qzone == "" {
									oc.PSkey.Qzone = b
								}
							}
						}
						qq, err := strconv.Atoi(uin)
						if err != nil {
							logrus.Errorln("ERROR: ", err)
							break LOOP
						}
						m = qzone.NewQzoneManager(int64(qq), oc)
						fmt.Printf("m:%#v\n", m)
						ctx.SendChain(message.Text("登录成功"))
						return
					}
				}
			}
			ctx.SendChain(message.Text("登录失败"))
		})
	engine.OnRegex(`^发说说\s+(.*)$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			err := m.RefreshToken()
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			msg := ctx.State["regex_matched"].([]string)
			sssr, err := m.SendShuoShuo(msg[1])
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			fmt.Printf("sssr:%#v\n", sssr)
		})
}

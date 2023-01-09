// Package asoularticle asoul.icu小作文
package asoularticle

import (
	"fmt"
	"net/url"

	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"jaytaylor.com/html2text"
)

func init() { // 插件主体
	engine := control.Register("asoularticle", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "小作文",
		Help:              "- 小作文",
		PrivateDataFolder: "asoularticle",
	})
	go func() {
		adb = initialize(engine.DataFolder() + "asoularticle.db")
	}()
	engine.OnFullMatch(`小作文`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			a, err := adb.randomArticle()
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			u := fmt.Sprintf(detailURL, url.QueryEscape(a.Title))
			data, err := web.RequestDataWith(web.NewDefaultClient(), u, "GET", u, web.RandUA())
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			tex, err := html2text.FromString(gjson.ParseBytes(data).Get("htmlContent").String(), html2text.Options{PrettyTables: true})
			if err != nil {
				panic(err)
			}
			ctx.SendChain(message.Text(tex))
		})
	engine.OnFullMatch(`更新小作文`, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("少女祈祷中..."))
			err := adb.truncateAndInsert()
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Text("asoul.icu小作文更新完毕"))
		})
}

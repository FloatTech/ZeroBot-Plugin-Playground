// Package asoularticle asoul.icu小作文
package asoularticle

import (
	"net/url"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
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
			tex, err := url.QueryUnescape(a.Brief)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
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

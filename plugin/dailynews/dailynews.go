// Package dailynews 今日早报
package dailynews

import (
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() { // 插件主体
	engine := control.Register("dailynews", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "今日早报",
		Help:              "- 今日早报",
		PrivateDataFolder: "dailynews",
	})

	// 开启
	engine.OnKeyword(`今日早报`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			data, err := web.RequestDataWith(web.NewDefaultClient(), "http://dwz.2xb.cn/zaob", "GET", "", "")
			if err != nil {
				return
			}
			picURL := gjson.Get(string(data), "imageUrl").String()
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Image(picURL))
		})
}

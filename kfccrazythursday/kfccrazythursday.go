// Package kfccrazythursday 疯狂星期四
package kfccrazythursday

import (
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	crazyURL = "https://kfc-crazy-thursday.vercel.app/api/index"
)

func init() {
	engine := control.Register("kfccrazythursday", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "疯狂星期四\n疯狂星期四\n",
	})
	engine.OnFullMatch("疯狂星期四").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		data, err := web.GetData(crazyURL)
		if err != nil {
			ctx.SendChain(message.Text("Error:", err))
		}
		ctx.SendChain(message.Text(binary.BytesToString(data)))
	})
}

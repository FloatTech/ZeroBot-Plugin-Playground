// Package diablo 动感地带
package diablo

import (
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const api = "https://127.0.0.1/tz"

func init() {
	engine := control.AutoRegister(&ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "动感地带",
		Help:              "- 动感地带助手，每小时发送动感地带信息",
		PrivateDataFolder: "diablo",
	})

	engine.OnFullMatch(`动感地带`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			data, err := web.GetData(api)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			picURL := gjson.Get(binary.BytesToString(data), "imageUrl").String()
			dqqy := gjson.Get(binary.BytesToString(data), "dq").String()
			xgqy := gjson.Get(binary.BytesToString(data), "xg").String()
			sj := gjson.Get(binary.BytesToString(data), "sj").String()

			ctx.SendChain(message.Text(
			"当前区域：",dqqy,"\n",
			"下个区域：",xgqy,"\n",
			"下次更新时间",sj,
			))
		})
}

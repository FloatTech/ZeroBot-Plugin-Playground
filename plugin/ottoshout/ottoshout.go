// Package ottoshout otto鬼叫
package ottoshout

import (
	"net/url"
	"strings"

	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	websiteURL = "https://www.aolianfeiallin.top"
	makeURL    = websiteURL + "/make"
)

func init() { // 插件主体
	engine := control.Register("ottoshout", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief: "otto鬼叫" ,
		Help: 
			"- 电棍说 xxx",
		PrivateDataFolder: "ottoshout",
	})
	// 开启
	engine.OnPrefix(`电棍说`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			data, err := web.PostData(makeURL, "application/x-www-form-urlencoded", strings.NewReader("text="+url.QueryEscape(args)+"&inYsddMode=true&norm=true&reverse=false&speedMult=1&pitchMult=1"))
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
			}
			ctx.SendChain(message.Record(websiteURL + "/get/" + gjson.ParseBytes(data).Get("id").String() + ".mp3"))
		})
}

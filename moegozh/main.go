// Package moegoe 日韩 VITS 模型拟声
package moegoezh

import (
	"fmt"
	"net/url"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	gofyzhq "github.com/guohuiyuan/go-fyzhq"
)

const (
	jpapi = "https://moegoe.azurewebsites.net/api/speak?text=%s&id=%d"
	krapi = "https://moegoe.azurewebsites.net/api/speakkr?text=%s&id=%d"
)

var speakers = map[string]uint{
	"宁宁": 0, "爱瑠": 1, "芳乃": 2, "茉子": 3, "丛雨": 4, "小春": 5, "七海": 6,
	"수아": 0, "미미르": 1, "아린": 2, "연화": 3, "유화": 4, "선배": 5,
}

func init() {
	en := control.Register("moegoe", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "moegoe\n" +
			"- 让[宁宁|爱瑠|芳乃|茉子|丛雨|小春|七海]用中文说(中文)\n" +
			"- 让[수아|미미르|아린|연화|유화|선배]用中文说(韩语)",
	}).ApplySingle(ctxext.DefaultSingle)
	en.OnRegex("^让(宁宁|爱瑠|芳乃|茉子|丛雨|小春|七海)用中文说([A-Za-z\\s\\d\u4E00-\u9FA5.。,，、:：;；]+)$").Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			text := ctx.State["regex_matched"].([]string)[2]
			id := speakers[ctx.State["regex_matched"].([]string)[1]]
			text = gofyzhq.Transform(text, "jp")
			// ctx.SendChain(message.Text("转换结果: ", text))
			ctx.SendChain(message.Record(fmt.Sprintf(jpapi, url.QueryEscape(text), id)))
		})
	en.OnRegex("^让(수아|미미르|아린|연화|유화|선배)说([A-Za-z\\s\\d\u4E00-\u9FA5.。,，、:：;；]+)$").Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			text := ctx.State["regex_matched"].([]string)[2]
			id := speakers[ctx.State["regex_matched"].([]string)[1]]
			text = gofyzhq.Transform(text, "kr")
			ctx.SendChain(message.Record(fmt.Sprintf(krapi, url.QueryEscape(text), id)))
		})
}

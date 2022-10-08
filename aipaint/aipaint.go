// Package aipaint ai画图
package aipaint

import (
	"fmt"
	"net/url"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	aipaintURL = "http://91.216.169.75:5010/got_image?tag=%v"
)

func init() { // 插件主体
	engine := control.Register("aipaint", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "ai画图\n" +
			"- ai画图xxx",
		PrivateDataFolder: "aipaint",
	})

	// 开启
	engine.OnPrefix(`ai画图`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			ctx.SendChain(message.Image(fmt.Sprintf(aipaintURL, url.QueryEscape(strings.TrimSpace(strings.ReplaceAll(args, " ", "%20"))))).Add("cache", 0))
		})
}

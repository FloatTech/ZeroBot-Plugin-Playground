// Package magicprompt 魔力提词
package magicprompt

import (
	"fmt"
	"strings"

	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const magicpromptURL = "https://api-inference.huggingface.co/models/Gustavosta/MagicPrompt-Stable-Diffusion"

func init() { // 插件主体
	engine := control.Register("magicprompt", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "MagicPrompt-Stable-Diffusion\n" +
			"- 魔力提词 xxx",
		PrivateDataFolder: "magicprompt",
	})

	// 开启
	engine.OnPrefix(`魔力提词`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			data, err := web.PostData(magicpromptURL, "application/json", strings.NewReader(fmt.Sprintf(`{"inputs": "%v"}`, args)))
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(gjson.ParseBytes(data).Get("@this.0.generated_text").String()))
		})
}

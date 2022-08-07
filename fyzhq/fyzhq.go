// Package fyzhq 发音转换
package fyzhq

import (
	"flag"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	gofyzhq "github.com/guohuiyuan/go-fyzhq"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("fyzhq", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "发音转换器\n- /发音转换 -l jp -t 你好世界 (语言: jp,en,fr,gm,ru,kr,th)",
	})
	engine.OnCommand("发音转换").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		fset := flag.FlagSet{}
		var (
			language string
			text     string
		)
		fset.StringVar(&language, "l", "jp", "jp,en,fr,gm,ru,kr,th")
		fset.StringVar(&text, "t", "你好世界", "转换文本")
		arguments := shell.Parse(ctx.State["args"].(string))
		err := fset.Parse(arguments)
		if err != nil {
			ctx.SendChain(message.Text("Error:", err))
			return
		}
		ctx.Send(gofyzhq.Transform(text, language))
	})
}

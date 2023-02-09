package chatgpt

import (
	"os"

	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var apiKey string

func init() {
	engine := control.Register("chatgpt", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "chatgpt",
		Help:              "- chatgpt",
		PrivateDataFolder: "chatgpt",
	})
	apikeyfile := engine.DataFolder() + "apikey.txt"
	if file.IsExist(apikeyfile) {
		apikey, err := os.ReadFile(apikeyfile)
		if err != nil {
			panic(err)
		}
		apiKey = string(apikey)
	}
	engine.OnRegex(`^chatgpt\s*(.*)$`, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["regex_matched"].([]string)[1]
			ans, err := completions(args, apiKey)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Text(ans))
		})
	engine.OnRegex(`^设置\s*ChatGPT\s*apikey\s*(.*)$`, zero.OnlyPrivate, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		apiKey = ctx.State["regex_matched"].([]string)[1]
		f, err := os.Create(apikeyfile)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		defer f.Close()
		_, err = f.WriteString(apiKey)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(message.Text("设置成功"))
	})
}

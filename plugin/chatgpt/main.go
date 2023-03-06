// Package chatgpt 简易ChatGPT api聊天
package chatgpt

import (
	"os"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/ttl"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var apiKey string

type sessionKey struct {
	group int64
	user  int64
}

var cache = ttl.NewCache[sessionKey, []chatMessage](time.Minute * 15)

func init() {
	engine := control.Register("chatgpt", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "chatgpt",
		Help: "-@bot chatgpt [对话内容]\n" +
			"不支持上下文且响应较慢\n" +
			"(私聊发送)设置OpenAI apikey [apikey]",
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
			key := sessionKey{
				group: ctx.Event.GroupID,
				user:  ctx.Event.UserID,
			}
			if args == "reset" {
				cache.Delete(key)
				ctx.Send("已清除上下文！")
				return
			}
			messages := cache.Get(key)
			messages = append(messages, chatMessage{
				Role:    "user",
				Content: args,
			})
			ans, err := completions(messages, apiKey)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text(ans))
		})
	engine.OnRegex(`^设置\s*OpenAI\s*apikey\s*(.*)$`, zero.OnlyPrivate, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
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

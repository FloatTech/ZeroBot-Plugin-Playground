// Package chatgpt 简易ChatGPT api聊天
package chatgpt

import (
	"strings"
	"time"

	"github.com/FloatTech/ttl"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const chatgptapikeygid = -3

type sessionKey struct {
	group int64
	user  int64
}

var (
	apiKey string
	cache  = ttl.NewCache[sessionKey, []chatMessage](time.Minute * 15)
	engine = control.Register("chatgpt", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "chatgpt",
		Help: "-@bot chatgpt [对话内容]\n" +
			"添加预设xxx xxx\n" +
			"设置预设xxx\n" +
			"删除本群预设\n" +
			"查看预设列表\n" +
			"(私聊发送)设置OpenAI apikey [apikey]",
		PrivateDataFolder: "chatgpt",
	})
)

func init() {
	engine.OnRegex(`^(?:chatgpt|//)([\s\S]*)$`, zero.OnlyToMe, getdb).SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["regex_matched"].([]string)[1]
			key := sessionKey{
				group: ctx.Event.GroupID,
				user:  ctx.Event.UserID,
			}
			if args == "reset" || args == "重置记忆" {
				cache.Delete(key)
				ctx.SendChain(message.Text("已清除上下文！"))
				return
			}
			// 添加预设
			var messages []chatMessage
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			content, err := db.findgroupmode(gid)
			if err == nil {
				messages = append(messages, chatMessage{
					Role:    "system",
					Content: content,
				})
				if len(cache.Get(key)) > 1 {
					messages = append(messages, cache.Get(key)[1:]...)
				}
			} else {
				messages = append(messages, cache.Get(key)...)
			}
			messages = append(messages, chatMessage{
				Role:    "user",
				Content: args,
			})
			resp, err := completions(messages, apiKey)
			if err != nil {
				ctx.SendChain(message.Text("请求ChatGPT失败: ", err))
				return
			}
			reply := resp.Choices[0].Message
			reply.Content = strings.TrimSpace(reply.Content)
			messages = append(messages, reply)
			cache.Set(key, messages)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(reply.Content))
		})
	engine.OnRegex(`^设置\s*OpenAI\s*apikey\s*(.*)$`, zero.OnlyPrivate, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		apiKey = ctx.State["regex_matched"].([]string)[1]
		m := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
		_ = m.Manager.Response(chatgptapikeygid)
		err := m.Manager.SetExtra(chatgptapikeygid, apiKey)
		if err != nil {
			ctx.SendChain(message.Text("保存apikey失败"))
			return
		}
		ctx.SendChain(message.Text("保存apikey成功"))
	})
	engine.OnRegex(`^添加预设\s*(\S+)\s+(.*)$`, zero.SuperUserPermission, getdb).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			modename := ctx.State["regex_matched"].([]string)[1]
			content := ctx.State["regex_matched"].([]string)[2]
			err := db.insertmode(modename, content)
			if err != nil {
				ctx.SendChain(message.Text("添加失败: ", err))
				return
			}
			ctx.SendChain(message.Text("添加成功"))
		})
	engine.OnRegex(`^设置预设\s*(\S+)$`, getdb).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			modename := ctx.State["regex_matched"].([]string)[1]
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			err := db.changemode(gid, modename)
			if err != nil {
				ctx.SendChain(message.Text("设置失败: ", err))
				return
			}
			for _, v := range ctx.GetThisGroupMemberListNoCache().Array() {
				cache.Delete(
					sessionKey{
						group: ctx.Event.GroupID,
						user:  v.Get("user_id").Int(),
					})
			}
			ctx.SendChain(message.Text("设置成功"))
		})
	engine.OnFullMatch("删除本群预设", getdb).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			err := db.delgroupmode(gid)
			if err != nil {
				ctx.SendChain(message.Text("删除失败: ", err))
				return
			}
			for _, v := range ctx.GetThisGroupMemberListNoCache().Array() {
				cache.Delete(
					sessionKey{
						group: ctx.Event.GroupID,
						user:  v.Get("user_id").Int(),
					})
			}
			ctx.SendChain(message.Text("删除成功"))
		})
	engine.OnFullMatch("查看预设列表", getdb).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			pre, err := db.findformode()
			if err != nil {
				ctx.SendChain(message.Text("当前没有任何预设: ", err))
				return
			}
			ctx.SendChain(message.Text(pre))
		})
}

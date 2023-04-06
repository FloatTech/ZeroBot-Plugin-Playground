// Package chatgpt 简易ChatGPT api聊天
package chatgpt

import (
	//"encoding/json"
	//"strconv"
	"strings"
	"time"

	//"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/ttl"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type sessionKey struct {
	group int64
	user  int64
}

var (
	cache  = ttl.NewCache[sessionKey, []chatMessage](time.Minute * 15)
	engine = control.Register("chatgpt", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "chatgpt",
		Help: "-@bot chatgpt [对话内容]\n" +
			"添加预设xxx xxx\n" +
			"设置(默认)预设xxx\n" +
			"删除本群预设\n" +
			"查看预设列表\n" +
			//"余额查询\n" +
			"(私聊发送)设置OpenAI apikey [apikey]" +
			"(私聊发送)删除apikey" +
			"(群聊发送)授权||取消本群使用apikey" +
			"注:先私聊设置自己的key,再授权群聊使用,不会泄露key的",
		PrivateDataFolder: "chatgpt",
	})
)

func init() {
	engine.OnRegex(`^(?:chatgpt|//)([\s\S]*)$`, zero.OnlyToMe, getdb).SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			var (
				apiKey   string
				err      error
				messages []chatMessage
			)
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
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
				apiKey, err = db.findkey(gid)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
			} else {
				apiKey, err = db.findgkey(gid)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
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
				c, err := db.findgroupmode(-1)
				if err != nil {
					messages = append(messages, cache.Get(key)...)
				} else {
					messages = append(messages, chatMessage{
						Role:    "system",
						Content: c,
					})
					if len(cache.Get(key)) > 1 {
						messages = append(messages, cache.Get(key)[1:]...)
					}
				}
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
	engine.OnRegex(`^设置\s*OpenAI\s*apikey\s*(.*)$`, zero.OnlyPrivate, getdb).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		err := db.insertkey(-ctx.Event.UserID, ctx.State["regex_matched"].([]string)[1])
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		ctx.SendChain(message.Text("保存apikey成功"))
	})
	engine.OnFullMatch("删除apikey", getdb).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			err := db.delkey(-ctx.Event.UserID)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
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
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("添加成功"))
		})
	engine.OnRegex(`^设置(默认)?预设\s*(\S+)$`, getdb).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			modename := ctx.State["regex_matched"].([]string)[2]
			gid := ctx.Event.GroupID
			if ctx.State["regex_matched"].([]string)[1] == "" {
				if gid == 0 {
					gid = -ctx.Event.UserID
				}
			} else {
				gid = -1 //全局为-1的群号
			}
			err := db.changemode(gid, modename)
			if err != nil {
				ctx.SendChain(message.Text("设置失败: ", err))
				return
			}
			ctx.SendChain(message.Text("设置成功"))
			for _, v := range ctx.GetThisGroupMemberListNoCache().Array() {
				cache.Delete(
					sessionKey{
						group: ctx.Event.GroupID,
						user:  v.Get("user_id").Int(),
					})
			}
			ctx.SendChain(message.Text("本群记忆清除成功"))
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
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("删除成功"))
			for _, v := range ctx.GetThisGroupMemberListNoCache().Array() {
				cache.Delete(
					sessionKey{
						group: ctx.Event.GroupID,
						user:  v.Get("user_id").Int(),
					})
			}
			ctx.SendChain(message.Text("本群记忆清除成功"))
		})
	engine.OnFullMatch("查看预设列表", getdb).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			pre, err := db.findformode()
			if err != nil {
				ctx.SendChain(message.Text(message.Reply(ctx.Event.MessageID), "当前没有任何预设: ", err))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(pre))
		})
	/*engine.OnFullMatch("余额查询", getdb).SetBlock(true).
	Handle(func(ctx *zero.Ctx) {
		data, err := web.GetData("https://v1.apigpt.cn/key/?key=" + apiKey)
		if err != nil {
			ctx.SendChain(message.Text("请求网站失败,网站可能跑路惹"))
			return
		}
		var all chatkeymessage
		err = json.Unmarshal(data, &all)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		if all.Code != 200 {
			ctx.SendChain(message.Text("请求key错误", err))
			return
		}
		var msg strings.Builder
		msg.WriteString(all.Msg)
		msg.WriteString("\n总量：$")
		msg.WriteString(strconv.FormatFloat(all.TotalGranted, 'f', 2, 64))
		msg.WriteString("\n剩余：$")
		msg.WriteString(strconv.FormatFloat(all.TotalAvailable, 'f', 2, 64))
		msg.WriteString("\n注册时间：")
		tm := time.Unix(all.EffectiveAt, 0)
		msg.WriteString(tm.Format("2006-01-02 15:04:05")) // 格式化时间
		msg.WriteString("\n到期时间：")
		tm = time.Unix(all.ExpiresAt, 0)
		msg.WriteString(tm.Format("2006-01-02 15:04:05")) // 格式化时间
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(msg.String()))
	})*/
	engine.OnRegex(`^(取消|授权)本群使用apikey$`, getdb).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if ctx.State["regex_matched"].([]string)[1] == "授权" {
				err := db.insertgkey(-ctx.Event.UserID, ctx.Event.GroupID)
				if err != nil {
					ctx.SendChain(message.Text("授权失败: ", err))
					return
				}
				ctx.SendChain(message.Text("授权成功"))
				return
			}
			t, err := db.findgtoqq(ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text("取消失败: ", err))
				return
			}
			if t != ctx.Event.UserID {
				ctx.SendChain(message.Text("取消失败: 你不是授权用户"))
				return
			}
			err = db.delgkey(ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Text("取消失败: ", err))

				return
			}
			ctx.SendChain(message.Text("取消成功"))
		})
}

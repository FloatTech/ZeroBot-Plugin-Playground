// Package chatgpt 简易ChatGPT api聊天
package chatgpt

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
			"添加预设(0-9) xxx\n" +
			"设置预设(0-9)\n" +
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
	//初始化文件路径
	if file.IsNotExist(engine.DataFolder() + "system") {
		err := os.MkdirAll(engine.DataFolder()+"system", 0777)
		if err != nil {
			return
		}
	}
	if file.IsNotExist(engine.DataFolder() + "group") {
		err := os.MkdirAll(engine.DataFolder()+"group", 0777)
		if err != nil {
			return
		}
	}
	engine.OnRegex(`^(?:chatgpt|//)\s*(.*)$`, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if apiKey == "" {
				ctx.SendChain(message.Text("未设置OpenAI apikey"))
				return
			}
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
			//添加预设
			var messages []chatMessage
			num, err := os.ReadFile(engine.DataFolder() + "group/" + strconv.Itoa(int(key.group)) + ".txt")
			if err == nil {
				txt, err := os.ReadFile(engine.DataFolder() + "system/" + string(num) + ".txt")
				if err != nil {
					ctx.SendChain(message.Text("预设不存在!请重新设置预设"))
					return
				}
				messages = append(messages, chatMessage{
					Role:    "system",
					Content: string(txt),
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
	engine.OnRegex(`^添加预设\s*(\S+)\s+(.*)$`, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			num := ctx.State["regex_matched"].([]string)[1]
			word := ctx.State["regex_matched"].([]string)[2]
			if word == "" {
				return
			}
			file, err := os.OpenFile(engine.DataFolder()+"system/"+num+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			_, _ = file.WriteString(word)
			file.Close()
			ctx.SendChain(message.Text("设置成功"))
		})
	engine.OnRegex(`^设置预设\s*(\S+)$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			num := ctx.State["regex_matched"].([]string)[1]
			group := strconv.Itoa(int(ctx.Event.GroupID))
			file, err := os.OpenFile(engine.DataFolder()+"group/"+group+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			_, _ = file.WriteString(num)
			file.Close()
			ctx.SendChain(message.Text("设置成功"))
		})
	engine.OnRegex(`^删除预设$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			group := strconv.Itoa(int(ctx.Event.GroupID))
			err := os.Remove(engine.DataFolder() + "group/" + group + ".txt")
			if err != nil {
				//如果删除失败则输出 file remove Error!
				ctx.SendChain(message.Text("未设置预设"))
			} else {
				//如果删除成功则输出 file remove OK!
				ctx.SendChain(message.Text("删除成功"))
			}
		})
	engine.OnRegex(`^查看预设列表$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var buff strings.Builder
			buff.WriteString("当前拥有预设:")
			err := filepath.Walk(engine.DataFolder()+"system", func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					buff.WriteByte('\n')
					buff.WriteString(info.Name())
				}
				return nil
			})
			if err != nil {
				ctx.SendChain(message.Text("Error:", err))
				return
			}
			ctx.SendChain(message.Text(buff.String()))
		})
}

// Package chinesebqb 表情包
package chinesebqb

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/binary"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() { // 插件主体
	engine := control.Register("chinesebqb", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "chinesebqb表情包\n" +
			"- 表情包",
		PrivateDataFolder: "chinesebqb",
	})
	go func() {
		bdb = initialize(engine.DataFolder() + "chinesebqb.db")
	}()
	engine.OnPrefix(`表情包`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("少女祈祷中..."))
			args := ctx.State["args"].(string)
			blist, err := bdb.getByKey(args)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			imageList := make([]string, 0, 10)
			for _, v := range blist {
				imageList = append(imageList, v.URL)
			}
			rand.Shuffle(len(imageList), func(i, j int) {
				imageList[i], imageList[j] = imageList[j], imageList[i]
			})
			m := message.Message{}
			for _, v := range imageList[:10] {
				m = append(m, ctxext.FakeSenderForwardNode(ctx, message.Image(v)))
			}
			if id := ctx.Send(m).ID(); id == 0 {
				ctx.SendChain(message.Text("ERROR: 可能被风控或下载图片用时过长，请耐心等待"))
			}
		})
	engine.OnFullMatch(`系列表情包`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			next := zero.NewFutureEvent("message", 999, false, ctx.CheckSession())
			recv, cancel := next.Repeat()
			defer cancel()
			results, err := bdb.getAllCategory()
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			tex := "请输入系列表情包序号\n"
			for i, v := range results {
				tex += fmt.Sprintf("%d. %s\n", i, v.Category)
			}
			base64Str, err := text.RenderToBase64(tex, text.FontFile, 400, 20)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Image("base64://" + binary.BytesToString(base64Str)))
			for {
				select {
				case <-time.After(time.Second * 120):
					ctx.SendChain(message.Text("系列表情包指令过期"))
					return
				case c := <-recv:
					msg := c.Event.Message.ExtractPlainText()
					num, err := strconv.Atoi(msg)
					if err != nil {
						ctx.SendChain(message.Text("请输入数字!"))
						continue
					}
					if num < 0 || num >= len(results) {
						ctx.SendChain(message.Text("序号非法!"))
						continue
					}
					ctx.SendChain(message.Text("请欣赏系列表情包: ", results[num].Category))
					blist, err := bdb.getByCategory(results[num].Category)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						return
					}
					imageList := make([]string, 0, 50)
					for _, v := range blist {
						imageList = append(imageList, v.URL)
					}
					rand.Shuffle(len(imageList), func(i, j int) {
						imageList[i], imageList[j] = imageList[j], imageList[i]
					})
					m := message.Message{}
					for _, v := range imageList[:50] {
						m = append(m, ctxext.FakeSenderForwardNode(ctx, message.Image(v)))
					}
					if id := ctx.Send(m).ID(); id == 0 {
						ctx.SendChain(message.Text("ERROR: 可能被风控或下载图片用时过长，请耐心等待"))
					}
					return
				}
			}
		})
	engine.OnFullMatch(`更新表情包`, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("少女祈祷中..."))
			err := bdb.truncateAndInsert()
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Text("chinesebqb表情包更新完毕"))
		})

}

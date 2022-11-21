// Package chinesebqb 表情包
package chinesebqb

import (
	"math/rand"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
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
			imageList := make([]string, 0)
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

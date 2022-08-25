// Package tgyj 同归于尽
package tgyj

import (
	"math/rand"
	"strconv"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/math"
)

func init() {
	engine := control.Register("tgyj", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "同归于尽@xxx",
	})
	engine.OnRegex(`同归于尽.*?(\d+)`, zero.OnlyGroup).Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			rand.Seed(time.Now().Unix())                                   //种子
			群员禁言 := rand.Intn(3) + 1                                       //生成群员禁言随机数 避免是0
			我 := 群员禁言 + rand.Intn(2)                                       //生成自己禁言随机数
			群员 := math.Str2Int64(ctx.State["regex_matched"].([]string)[1]) //字符串转整数

			if 群员 == 1603798048 {
				ctx.SendChain(message.Text("达咩~")) //	自己白名单
				return
			}
			ctx.SetGroupBan(ctx.Event.GroupID, 群员, int64(群员禁言)*60) //设置禁言
			ctx.SetGroupBan(ctx.Event.GroupID, ctx.Event.UserID, int64(我)*60)
			ctx.SendChain(message.At(ctx.Event.UserID),
				message.Text("\n你向"),
				message.At(群员),
				message.Text("发动了 同归于尽 对方获得"+strconv.Itoa(int(群员禁言))+"分钟禁言"+"  爆炸伤害波及到自己  自己获得"+strconv.Itoa(int(我))+"分钟禁言"),
			)

		})

}

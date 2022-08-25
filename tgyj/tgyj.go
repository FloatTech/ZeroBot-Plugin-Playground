// Package tgyj 同归于尽
package tgyj

import (
	"math/rand"
	"strconv"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/floatbox/math"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
)

func init() {
	engine := control.Register("tgyj", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "同归于尽@xxx",
	})
	engine.OnRegex(`同归于尽.*?(\d+)`, zero.OnlyGroup).Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			userban := rand.Intn(3) + 1                                      //生成群员禁言随机数 避免是0
			me := userban + rand.Intn(2)                                     //生成自己禁言随机数
			user := math.Str2Int64(ctx.State["regex_matched"].([]string)[1]) //字符串转整数

			ctx.SetGroupBan(ctx.Event.GroupID, user, int64(userban)*60) //设置禁言
			ctx.SetGroupBan(ctx.Event.GroupID, ctx.Event.UserID, int64(me)*60)
			ctx.SendChain(message.At(ctx.Event.UserID),
				message.Text("\n你向"),
				message.At(user),
				message.Text("发动了 同归于尽 对方获得"+strconv.Itoa(userban)+"分钟禁言"+"  爆炸伤害波及到自己  自己获得"+strconv.Itoa(me)+"分钟禁言"),
			)

		})

}

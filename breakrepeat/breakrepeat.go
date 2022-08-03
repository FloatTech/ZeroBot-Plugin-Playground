// Package breakrepeat 打断复读
package breakrepeat

import (
	"math/rand"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	maxLimit = 3
)

var (
	limit  = make(map[int64]int64, 256)
	rawMsg = make(map[int64]string, 256)
)

func init() {
	engine := control.Register("breakrepeat", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "打断复读,打断3次以上复读\n",
	})
	engine.On(`message/group`, zero.OnlyGroup).SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			gid := ctx.Event.GroupID
			raw := ctx.Event.RawMessage
			if rawMsg[gid] != raw {
				rawMsg[gid] = raw
				limit[gid] = 0
				return
			}
			limit[gid]++
			if limit[gid] >= maxLimit {
				if len([]rune(rawMsg[gid])) > 2 {
					ctx.Send(shuffle(rawMsg[gid]))
				}
				limit[gid] = 0
			}
		})
}

func shuffle(s string) string {
	r := []rune(s)
	for i := len(r) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		r[i], r[num] = r[num], r[i]
	}
	return string(r)
}

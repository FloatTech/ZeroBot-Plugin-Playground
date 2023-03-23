package fgopickup

import (
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"
)

func listPickups(ctx *zero.Ctx) {
	service := service{}
	pickups := service.getPickups()

	msg := message.Message{}
	for _, pickup := range pickups {
		banner := message.Image(pickup.Banner)
		name := message.Text("\n" + pickup.Name)
		date := message.Text("\n" +
			time.Unix(pickup.StartTime/1000, 0).Format("2006-01-02") +
			"~" + time.Unix(pickup.EndTime/1000, 0).Format("2006-01-02"))
		msg = append(msg, ctxext.FakeSenderForwardNode(ctx, banner, name, date))
	}
	ctx.Send(msg)
}

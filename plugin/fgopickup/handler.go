package fgopickup

import (
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func listPickups(ctx *zero.Ctx) {
	service := service{}
	pickups := service.getPickups()

	msg := message.Message{}
	for _, pickup := range pickups {
		banner := message.Image(pickup.Banner)
		name := message.Text("\n" + pickup.Name)
		date := message.Text("\n" + pickup.StartTime + "~" + pickup.EndTime)
		msg = append(msg, ctxext.FakeSenderForwardNode(ctx, banner, name, date))
	}
	ctx.Send(msg)
}

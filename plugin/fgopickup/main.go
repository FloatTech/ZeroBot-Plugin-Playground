package fgopickup

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("fgopickup", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "例",
		Help:             "- example 插件的帮助",
		// Banner: "",
		PrivateDataFolder: "fgopickup",
	})

	_ = initialize(engine.DataFolder() + "fgopickup.db")

	engine.OnFullMatch("test").
		SetBlock(true).
		Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text())
		})

	engine.OnFullMatch(`测试`).
		SetBlock(true).
		Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			var banner message.Message
			pickups := test()
			for _, pickup := range pickups {
				banner = append(banner, message.Image(pickup.Banner))
			}
			ctx.Send(message.Message{ctxext.FakeSenderForwardNode(ctx, banner...)})
		})
}

func test() []pickup {
	dao := dao{DbEngine: getOrmEngine()}
	list := dao.List()
	return *list
}

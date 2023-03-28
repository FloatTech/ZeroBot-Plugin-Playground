// Package fgopickup FGO卡池相关插件
package fgopickup

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	engine := control.Register("fgopickup", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "Fate/Grand Order",
		Help: "- fgo未来视\n" +
			"- fgo卡池[id]",
		// Banner: "",
		PrivateDataFolder: "fgopickup",
	})

	_ = initialize(engine.DataFolder() + "fgopickup.db")

	engine.OnFullMatch(`fgo未来视`).
		SetBlock(true).
		Limit(ctxext.LimitByGroup).
		Handle(listPickups)

	engine.OnPrefix(`fgo卡池`).
		SetBlock(true).
		Limit(ctxext.LimitByGroup).
		Handle(pickupDetail)

	engine.OnPrefix(`fgo卡池天数`).
		SetBlock(true).
		Limit(ctxext.LimitByGroup).
		Handle(getPickupTimeGap)
}

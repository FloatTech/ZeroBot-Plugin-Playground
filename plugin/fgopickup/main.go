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
		Brief:            "例",
		Help:             "- fgo未来视 fgo未来卡池",
		// Banner: "",
		PrivateDataFolder: "fgopickup",
	})

	_ = initialize(engine.DataFolder() + "fgopickup.db")

	engine.OnFullMatch(`fgo未来视`).
		SetBlock(true).
		Limit(ctxext.LimitByGroup).
		Handle(listPickups)
}

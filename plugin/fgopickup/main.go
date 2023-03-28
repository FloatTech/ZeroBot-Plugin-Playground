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
		Help: "- fgo未来视 (查询未来卡池)\n" +
			"- fgo卡池[id] (以卡池id查询卡池详情)\n" +
			"- fgo从者[page] (分页查询从者)",
		Banner:            "https://wx2.sinaimg.cn/large/0083LFbYgy1hcfkreklmbj31e012w7i5.jpg",
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

	engine.OnPrefix(`fgo从者`).
		SetBlock(true).
		Limit(ctxext.LimitByGroup).
		Handle(getServantList)
}

// Package killall 一键禁用本群所有插件
package killall

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() { // 主函数
	en := control.Register("killall", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "一键禁用插件",
		Help: "- 禁用所有插件[除本插件]\n" +
			"- 启用所有插件[除本插件]",
	})
	en.OnRegex(`^(启用|禁用)(全部|所有)插件`, zero.UserOrGrpAdmin).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		key := ctx.State["regex_matched"].([]string)[1]
		grp := ctx.Event.GroupID
		if grp == 0 {
			// 个人用户
			grp = -ctx.Event.UserID
		}
		if key == "启用" {
			control.ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
				manager.Enable(grp)
				return true
			})
		} else {
			control.ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
				manager.Disable(grp)
				return true
			})
		}
		service, _ := control.Lookup("killall")
		service.Enable(grp) // 排除本插件
		ctx.SendChain(message.Text("已" + key + "所有插件喵~"))
	})
}

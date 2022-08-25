// Package qqclean 清理
package qqclean

import (
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("qqclean", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "qq号清理\n清理群聊@bot\n",
	})
	engine.OnFullMatch("清理群聊", zero.SuperUserPermission, zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		cleanGroupnameList := make([]string, 0, 64)
		ctx.GetGroupList().ForEach(func(_, value gjson.Result) bool {
			if ctx.GetGroupMemberInfo(value.Get("group_id").Int(), ctx.Event.SelfID, true).Get("role").String() == "member" {
				cleanGroupnameList = append(cleanGroupnameList, value.Get("group_name").String())
			}
			ctx.SetGroupLeave(value.Get("group_id").Int(), false)
			return true
		})
		ctx.SendChain(message.Text("已清理bot群聊: (", strings.Join(cleanGroupnameList, ", "), ")"))
	})
}

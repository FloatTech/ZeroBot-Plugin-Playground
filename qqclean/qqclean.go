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
		Help:             "qq号清理\n清理群聊@bot (清理bot不是管理员的群)\n(当群聊数小于20时自动同意加群)\n",
	})
	engine.OnFullMatch("清理群聊", zero.SuperUserPermission, zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		cleanGroupnameList := make([]string, 0, 64)
		ctx.GetGroupList().ForEach(func(_, value gjson.Result) bool {
			if ctx.GetGroupMemberInfo(value.Get("group_id").Int(), ctx.Event.SelfID, true).Get("role").String() == "member" {
				cleanGroupnameList = append(cleanGroupnameList, value.Get("group_name").String())
				ctx.SetGroupLeave(value.Get("group_id").Int(), false)
			}
			return true
		})
		ctx.SendPrivateMessage(ctx.Event.UserID, message.Text("已清理bot群聊: (", strings.Join(cleanGroupnameList, ", "), ")"))
	})
	engine.On("request/group/invite").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		freecnt := 0
		ctx.GetGroupList().ForEach(func(_, value gjson.Result) bool {
			if ctx.GetGroupMemberInfo(value.Get("group_id").Int(), ctx.Event.SelfID, true).Get("role").String() == "member" {
				freecnt++
			}
			return true
		})
		if freecnt < 20 {
			ctx.SetGroupAddRequest(ctx.Event.Flag, "invite", true, "爱你哟")
		} else {
			ctx.SetGroupAddRequest(ctx.Event.Flag, "invite", false, "游离群聊大于20")
		}
	})
}

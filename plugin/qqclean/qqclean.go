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
		Brief:            "qq号清理",
		Help:             "清理群聊@bot (清理bot不是管理员的群)\n(当群聊数小于20时自动同意加群)\n清理好友@bot (清理bot 5级以下的好友)\n(自动同意好友邀请)",
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
	engine.OnFullMatch("清理好友", zero.SuperUserPermission, zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		cleanFriendnameList := make([]string, 0, 64)
		ctx.GetFriendList().ForEach(func(_, value gjson.Result) bool {
			if ctx.GetStrangerInfo(value.Get("user_id").Int(), true).Get("level").Int() <= 5 {
				ctx.CallAction("delete_friend", zero.Params{
					"friend_id": value.Get("user_id").Int(),
				})
				cleanFriendnameList = append(cleanFriendnameList, value.Get("nickname").String())
			}
			return true
		})
		ctx.SendPrivateMessage(ctx.Event.UserID, message.Text("已清理bot好友: (", strings.Join(cleanFriendnameList, ", "), ")"))
	})
}

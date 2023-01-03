// Package recall 回复撤回
package recall

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("recall", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "Recall",
		Help: "Reply a message with 'recall' or 'Recall' to recall it.\n" +
			"- Recall | recall",
	})
	// 完全匹配和部分匹配都不能解决需要匹配回复这个问题, 只能用正则.
	// 同时兼容了普通QQ客户端(即回复包括reply以及at两个CQ码)和ICA客户端(即仅有reply一个CQ码)
	engine.OnRegex(`^\[CQ:reply,id=(\d+|-\d+).*(?:\](\s{0,1}))(?:Recall|recall)$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) { // 这个正则还是有缺陷的, 比如这样的消息: [CQ:reply,id=123]]]]]] 112233 依然会匹配. 脑细胞不太够用,想不出什么合适的解决方法
		regex := ctx.State["regex_matched"].([]string)
		curmsgid := ctx.Event.MessageID.(int64)                      // 发送消息的消息ID
		ctx.DeleteMessage(message.NewMessageIDFromString(regex[1]))  // 尝试撤回目标消息
		ctx.DeleteMessage(message.NewMessageIDFromInteger(curmsgid)) // 尝试撤回匹配消息(回复recall)
	})
}

// 没了, 很好奇为什么没有人做撤回相关的插件
// 在MiraiConsolLoader上遇到了CQ码不能匹配的问题, 估计是不兼容. Go-cqhttp是正常的. 别的平台没有测试.

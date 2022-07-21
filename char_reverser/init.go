// Package CharReverser 英文字符反转
package CharReverser

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	// 初始化engine
	engine := control.Register(
		"CharReverser",
		&ctrl.Options[*zero.Ctx]{
			DisableOnDefault: false,
			Help: "字符翻转\n -翻转 <英文字符串>",
		},
	)
	// 处理字符翻转指令
	engine.OnRegex(`翻转( )+[A-z]+`).SetBlock(true).Handle(HandleReverse)
}
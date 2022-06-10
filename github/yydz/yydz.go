// Package yydz 丁真表情包
package yydz

import (
	"math/rand"
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"

	"github.com/FloatTech/ZeroBot-Plugin-Playground/github"
)

var (
	dingzhenPath = github.Githubhost + "/BlinkDL/YYDZ/tree/main/imgs"
)

func init() {
	control.Register("yydz", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "一眼丁真\n- 一眼丁真",
	}).OnFullMatch("一眼丁真", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			fs, _, err := github.Files(dingzhenPath)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
			}
			imageURL := strings.ReplaceAll(github.Githubhost+fs[rand.Intn(len(fs))].Href, "blob", "raw")
			ctx.SendChain(message.Image(imageURL))
		})
}

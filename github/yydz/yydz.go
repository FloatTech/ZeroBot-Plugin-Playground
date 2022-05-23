// Package yydz 丁真表情包
package yydz

import (
	"math/rand"
	"strings"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin-Playground/github"
	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
)

var (
	dingzhenPath = github.Githubhost + "/BlinkDL/YYDZ/tree/main/imgs"
)

func init() {
	control.Register("yydz", &control.Options{
		DisableOnDefault: false,
		Help:             "一眼丁真\n- 一眼丁真",
	}).OnFullMatch("一眼丁真", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			fs, _, err := github.GithubFiles(dingzhenPath)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
			}
			r := rand.New(rand.NewSource(time.Now().Unix()))
			imageURL := strings.ReplaceAll(github.Githubhost+fs[r.Intn(len(fs))].Href, "blob", "raw")
			ctx.SendChain(message.Image(imageURL))
		})
}

// Package fadian 发癫
package fadian

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"

	fbctxext "github.com/FloatTech/floatbox/ctxext"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	exampleName = "阿咪"
)

type fadian struct {
	Post []string `json:"post"`
}

var (
	fd fadian
)

func init() { // 插件主体
	engine := control.Register("fadian", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "每日发癫",
		Help:             "- 每日发癫 嘉然",
		PublicDataFolder: "Fadian",
	})

	// 开启
	engine.OnPrefix(`每日发癫`, fbctxext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		_, err := engine.GetLazyData("post.json", false)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		data, err := os.ReadFile(engine.DataFolder() + "post.json")
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		err = json.Unmarshal(data, &fd)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		return true
	})).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		name := ctx.NickName()
		text := fd.Post[rand.Intn(len(fd.Post))]
		ctx.SendChain(message.Text(strings.ReplaceAll(text, exampleName, name)))
	})
}

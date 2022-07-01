// Package gif 制图
package gif

import (
	"reflect"
	"strconv"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/file"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	cmds     = make([]string, 0)
	datapath string
	cmdMap   = map[string]string{
		"搓":  "Cuo",
		"冲":  "Xqe",
		"摸":  "Mo",
		"拍":  "Pai",
		"丢":  "Diu",
		"吃":  "Chi",
		"敲":  "Qiao",
		"啃":  "Ken",
		"蹭":  "Ceng",
		"爬":  "Pa",
		"撕":  "Si",
		"灰度": "Grayscale",
		"上翻": "FlipV",
		"下翻": "FlipV",
		"左翻": "FlipH",
		"右翻": "FlipH",
		"反色": "Invert",
		"浮雕": "Convolve3x3",
		"打码": "Blur",
		"负片": "InvertAndGrayscale",
		"亲":  "Kiss",
	}
)

func init() { // 插件主体
	en := control.Register("gif", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "制图\n- " + strings.Join(cmds, "\n- "),
		PrivateDataFolder: "gif",
	})
	datapath = file.BOTPATH + "/" + en.DataFolder()
	for k := range cmdMap {
		cmds = append(cmds, k)
	}
	en.OnRegex(`^(` + strings.Join(cmds, "|") + `)\D*?(\[CQ:(image\,file=([0-9a-zA-Z]{32}).*|at.+?(\d{5,11}))\].*|(\d+))$`).
		SetBlock(true).Handle(func(ctx *zero.Ctx) {
		c := newContext(ctx.Event.UserID)
		list := ctx.State["regex_matched"].([]string)
		err := c.prepareLogos(list[4]+list[5]+list[6], strconv.FormatInt(ctx.Event.UserID, 10))
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		r := reflect.ValueOf(c).MethodByName(cmdMap[list[1]]).Call(nil)
		picurl := r[0].String()
		if !r[1].IsNil() {
			err = r[1].Interface().(error)
		}
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		ctx.SendChain(message.Image(picurl))
	})
}

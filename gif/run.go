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
	cmd        = make([]string, 0)
	cmdNoimage = make([]string, 0)
	datapath   string
	cmdMap     = map[string]string{
		"搓":    "Cuo",
		"冲":    "Xqe",
		"摸":    "Mo",
		"拍":    "Pai",
		"丢":    "Diu",
		"吃":    "Chi",
		"敲":    "Qiao",
		"啃":    "Ken",
		"蹭":    "Ceng",
		"爬":    "Pa",
		"撕":    "Si",
		"灰度":   "Grayscale",
		"上翻":   "FlipV",
		"下翻":   "FlipV",
		"左翻":   "FlipH",
		"右翻":   "FlipH",
		"反色":   "Invert",
		"浮雕":   "Convolve3x3",
		"打码":   "Blur",
		"负片":   "InvertAndGrayscale",
		"旋转":   "Rotate",
		"变形":   "Deformation",
		"亲":    "Kiss",
		"娶":    "Marriage",
		"结婚申请": "Marriage",
		"结婚登记": "Marriage",
	}
	cmdNoimageMap = map[string]string{
		"像个": "Alike",
	}
)

func init() { // 插件主体
	for k := range cmdMap {
		cmd = append(cmd, k)
	}
	for k := range cmdNoimageMap {
		cmdNoimage = append(cmdNoimage, k)
	}
	en := control.Register("petpet", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "制图\n- " + strings.Join(cmd, "\n- ") + "\n- " + strings.Join(cmdNoimage, "\n- "),
		PrivateDataFolder: "petpet",
	})
	datapath = file.BOTPATH + "/" + en.DataFolder()
	en.OnRegex(`^(` + strings.Join(cmd, "|") + `)[\s\S]*?(\[CQ:(image\,file=([0-9a-zA-Z]{32}).*|at.+?(\d{5,11}))\].*|(\d+))$`).
		SetBlock(true).Handle(func(ctx *zero.Ctx) {
		c := newContext(ctx.Event.UserID)
		list := ctx.State["regex_matched"].([]string)
		err := c.prepareLogos(list[4]+list[5]+list[6], strconv.FormatInt(ctx.Event.UserID, 10))
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		argslist := strings.Split(strings.TrimSuffix(strings.TrimPrefix(list[0], list[1]), list[2]), " ")
		args := make([]reflect.Value, len(argslist))
		for i := 0; i < len(argslist); i++ {
			args[i] = reflect.ValueOf(argslist[i])
		}
		r := reflect.ValueOf(c).MethodByName(cmdMap[list[1]]).Call(args)
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
	en.OnRegex(`^(` + strings.Join(cmdNoimage, "|") + `)[\s\S]*`).
		SetBlock(true).Handle(func(ctx *zero.Ctx) {
		c := newContext(ctx.Event.UserID)
		list := ctx.State["regex_matched"].([]string)
		if len(list) < 2 {
			ctx.SendChain(message.Text("ERROR:", list[1], "后面要带参数"))
			return
		}
		argslist := strings.Split(strings.TrimPrefix(list[0], list[1]), " ")
		args := make([]reflect.Value, len(argslist))
		for i := 0; i < len(argslist); i++ {
			args[i] = reflect.ValueOf(argslist[i])
		}
		r := reflect.ValueOf(c).MethodByName(cmdNoimageMap[list[1]]).Call(args)
		picurl := r[0].String()
		var err error
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

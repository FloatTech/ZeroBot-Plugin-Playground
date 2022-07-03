// Package gif 制图
package gif

import (
	"reflect"
	"strconv"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/file"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	cmd      = make([]string, 0)
	datapath string
	cmdMap   = map[string]string{
		"搓":     "Cuo",
		"冲":     "Xqe",
		"摸":     "Mo",
		"拍":     "Pai",
		"丢":     "Diu",
		"吃":     "Chi",
		"敲":     "Qiao",
		"啃":     "Ken",
		"蹭":     "Ceng",
		"爬":     "Pa",
		"撕":     "Si",
		"灰度":    "Grayscale",
		"上翻":    "FlipV",
		"下翻":    "FlipV",
		"左翻":    "FlipH",
		"右翻":    "FlipH",
		"反色":    "Invert",
		"浮雕":    "Convolve3x3",
		"打码":    "Blur",
		"负片":    "InvertAndGrayscale",
		"旋转":    "Rotate",
		"变形":    "Deformation",
		"亲":     "Kiss",
		"娶":     "Marriage",
		"结婚申请":  "Marriage",
		"结婚登记":  "Marriage",
		"阿尼亚喜欢": "Anyasuki",
		"像只":    "Alike",
		"我永远喜欢": "AlwaysLike",
		"永远喜欢":  "AlwaysLike",
		"像样的亲亲": "DecentKiss",
		"国旗":    "ChinaFlag",
		"不要靠近":  "DontTouch",
		"万能表情":  "Universal",
		"空白表情":  "Universal",
		"采访":    "Interview",
		"需要":    "Need",
		"你可能需要": "Need",
		"这像画吗":  "Paint",
		"小画家":   "Painter",
		"完美":    "Perfect",
		"玩游戏":   "PlayGame",
		"出警":    "Police",
		"警察":    "Police1",
	}
)

func init() { // 插件主体
	for k := range cmdMap {
		cmd = append(cmd, k)
	}
	en := control.Register("petpet", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "制图\n- 搓\n- 冲\n- 摸\n- 拍\n- 丢\n- 吃\n- 敲\n- 啃\n- 蹭\n- 爬\n- 撕\n- 灰度\n- 上翻|下翻\n" +
			"- 左翻|右翻\n- 反色\n- 浮雕\n- 打码\n- 负片\n- 旋转 45\n- 变形 100 100\n- 亲\n- 娶|结婚申请|结婚登记\n- 阿尼亚喜欢\n- 像只\n" +
			"- 我永远喜欢|永远喜欢\n- 像样的亲亲\n- 国旗\n- 不要靠近\n- 万能表情|空白表情\n- 采访\n- 需要|你可能需要\n- 这像画吗\n- 小画家\n" +
			"- 完美\n- 玩游戏\n- 出警\n- 警察\n",
		PrivateDataFolder: "petpet",
	}).ApplySingle(ctxext.DefaultSingle)
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
}

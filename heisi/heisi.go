// Package heisi 每日新闻
package heisi

import (
	"math/rand"
	"os"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/ctxext"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	heisiPic []string
)

func init() { // 插件主体
	engine := control.Register("heisi", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "黑丝\n" +
			"- 随机黑丝",
		PublicDataFolder: "Heisi",
	})

	// 开启
	engine.OnFullMatch("随机黑丝", ctxext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		_, err := engine.GetLazyData("heisi.txt", false)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		data, err := os.ReadFile(engine.DataFolder() + "heisi.txt")
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		heisiPic = strings.Split(strings.ReplaceAll(binary.BytesToString(data), "\r", ""), "\n")
		logrus.Infoln("[黑丝]加载", len(heisiPic), "条黑丝图")
		return true
	})).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			pic := heisiPic[rand.Intn(len(heisiPic))]
			ctx.SendChain(message.Image(pic))
		})
}

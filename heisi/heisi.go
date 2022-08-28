// Package heisi 每日新闻
package heisi

import (
	"math/rand"
	"os"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	fbctxext "github.com/FloatTech/floatbox/ctxext"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	heisiPic []string
	baisiPic []string
	jkPic    []string
	jurPic   []string
	zukPic   []string
	mcnPic   []string
	fileList = []string{"heisi.txt", "baisi.txt", "jk.txt", "jur.txt", "zuk.txt", "mcn.txt"}
)

func init() { // 插件主体
	engine := control.Register("heisi", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "黑丝\n" +
			"- 来点黑丝\n- 来点白丝\n- 来点jk\n- 来点巨乳\n- 来点足控\n- 来点网红",
		PublicDataFolder: "Heisi",
	})

	// 开启
	engine.OnFullMatchGroup([]string{"来点黑丝", "来点白丝", "来点jk", "来点巨乳", "来点足控", "来点网红"}, zero.OnlyGroup, fbctxext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		for i, filePath := range fileList {
			_, err := engine.GetLazyData(filePath, false)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return false
			}
			data, err := os.ReadFile(engine.DataFolder() + filePath)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return false
			}
			switch i {
			case 0:
				heisiPic = strings.Split(strings.ReplaceAll(binary.BytesToString(data), "\r", ""), "\n")
			case 1:
				baisiPic = strings.Split(strings.ReplaceAll(binary.BytesToString(data), "\r", ""), "\n")
			case 2:
				jkPic = strings.Split(strings.ReplaceAll(binary.BytesToString(data), "\r", ""), "\n")
			case 3:
				jurPic = strings.Split(strings.ReplaceAll(binary.BytesToString(data), "\r", ""), "\n")
			case 4:
				zukPic = strings.Split(strings.ReplaceAll(binary.BytesToString(data), "\r", ""), "\n")
			case 5:
				mcnPic = strings.Split(strings.ReplaceAll(binary.BytesToString(data), "\r", ""), "\n")
			}
		}
		return true
	})).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			matched := ctx.State["matched"].(string)
			var pic string
			switch matched {
			case "来点黑丝":
				pic = heisiPic[rand.Intn(len(heisiPic))]
			case "来点白丝":
				pic = baisiPic[rand.Intn(len(baisiPic))]
			case "来点jk":
				pic = jkPic[rand.Intn(len(jkPic))]
			case "来点巨乳":
				pic = jurPic[rand.Intn(len(jurPic))]
			case "来点足控":
				pic = zukPic[rand.Intn(len(zukPic))]
			case "来点网红":
				pic = mcnPic[rand.Intn(len(mcnPic))]
			}
			m := message.Message{ctxext.FakeSenderForwardNode(ctx, message.Image(pic))}
			if id := ctx.SendGroupForwardMessage(
				ctx.Event.GroupID,
				m).Get("message_id").Int(); id == 0 {
				ctx.SendChain(message.Text("ERROR: 可能被风控或下载图片用时过长，请耐心等待"))
			}
		})
}

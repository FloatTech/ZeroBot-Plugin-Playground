// Package aipaint ai画图
package aipaint

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/img/writer"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	mytoken           = "06LhgOew9PJDFQKnfdcSI3BtXz84AGoM"
	aipaintHost       = "http://91.216.169.75:5010"
	aipaintTxt2ImgURL = aipaintHost + "/got_image?token=%v&tags=%v"
	aipaintImg2ImgURL = aipaintHost + "/got_image2image?token=%v&tags=%v"
	// aipaintURL  = "https://22229.gradio.app/api/predict/"
	// sessionHash = "zerobot"
)

var (
	datapath  string
	predictRe = regexp.MustCompile(`{"steps".+?}`)
)

// type request struct {
// 	FnIndex     int           `json:"fn_index"`
// 	Data        []interface{} `json:"data"`
// 	SessionHash string        `json:"session_hash"`
// }

func init() { // 插件主体
	engine := control.Register("aipaint", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "ai绘图\n" +
			"[ai绘图|生成色图|生成涩图|ai画图] xxx\n" +
			"[以图绘图|以图生图|以图画图] xxx [图片]|@xxx|[qq号]",
		// "- 画1张a photo of sks toy riding a bicycle",
		PrivateDataFolder: "aipaint",
	})
	datapath = file.BOTPATH + "/" + engine.DataFolder()
	engine.OnPrefixGroup([]string{`ai绘图`, `生成色图`, `生成涩图`, `ai画图`}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("少女祈祷中..."))
			args := ctx.State["args"].(string)
			data, err := web.GetData(aipaintTxt2ImgURL)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			var loadData string
			if predictRe.MatchString(binary.BytesToString(data)) {
				loadData = predictRe.FindStringSubmatch(binary.BytesToString(data))[0]
			}
			m := message.Message{ctxext.FakeSenderForwardNode(ctx, message.Image(fmt.Sprintf(aipaintTxt2ImgURL, mytoken, url.QueryEscape(strings.TrimSpace(strings.ReplaceAll(args, " ", "%20"))))).Add("cache", 0))}
			m = append(m, ctxext.FakeSenderForwardNode(ctx, message.Text("seed: ", gjson.Get(loadData, "seed").Int(), "\n", "scale: ", gjson.Get(loadData, "scale").Float())))
			if id := ctx.Send(m).ID(); id == 0 {
				ctx.SendChain(message.Text("ERROR: 可能被风控或下载图片用时过长，请耐心等待"))
			}
		})
	engine.OnRegex(`^(以图绘图|以图生图|以图画图)[\s\S]*?(\[CQ:(image\,file=([0-9a-zA-Z]{32}).*|at.+?(\d{5,11}))\].*|(\d+))$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			c := newContext(ctx.Event.UserID)
			list := ctx.State["regex_matched"].([]string)
			err := c.prepareLogos(list[4]+list[5]+list[6], strconv.FormatInt(ctx.Event.UserID, 10))
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			args := strings.TrimSuffix(strings.TrimPrefix(list[0], list[1]), list[2])
			if args == "" {
				ctx.SendChain(message.Text("ERROR: 以图绘图必须添加tag"))
				return
			}
			ctx.SendChain(message.Text("少女祈祷中..."))
			postURL := fmt.Sprintf(aipaintImg2ImgURL, mytoken, url.QueryEscape(strings.TrimSpace(strings.ReplaceAll(args, " ", "%20"))))

			f, err := os.Open(c.headimgsdir[0])
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			imageShape := ""
			switch {
			case img.Bounds().Dx() > img.Bounds().Dy():
				imageShape = "Landscape"
			case img.Bounds().Dx() == img.Bounds().Dy():
				imageShape = "Square"
			default:
				imageShape = "Portrait"
			}

			// 图片转base64
			base64Bytes, err := writer.ToBase64(img)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			data, err := web.PostData(postURL+"&shape="+imageShape, "text/plain", bytes.NewReader(base64Bytes))
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			var loadData string
			if predictRe.MatchString(binary.BytesToString(data)) {
				loadData = predictRe.FindStringSubmatch(binary.BytesToString(data))[0]
			}
			encodeStr := base64.StdEncoding.EncodeToString(data)
			m := message.Message{ctxext.FakeSenderForwardNode(ctx, message.Image("base64://"+encodeStr))}
			m = append(m, ctxext.FakeSenderForwardNode(ctx, message.Text("seed: ", gjson.Get(loadData, "seed").Int(), "\n", "scale: ", gjson.Get(loadData, "scale").Float())))
			if id := ctx.Send(m).ID(); id == 0 {
				ctx.SendChain(message.Text("ERROR: 可能被风控或下载图片用时过长，请耐心等待"))
			}
		})

	// engine.OnRegex(`^画(\d{0,3})张([\s\S]*)$`).SetBlock(true).
	// 	Handle(func(ctx *zero.Ctx) {
	// 		regexMatched := ctx.State["regex_matched"].([]string)
	// 		count, err := strconv.Atoi(regexMatched[1])
	// 		if err != nil {
	// 			ctx.SendChain(message.Text("Error:", err))
	// 			return
	// 		}
	// 		if count > 100 {
	// 			count = 100
	// 		}
	// 		ctx.SendChain(message.Text("少女祈祷中..."))
	// 		r := request{
	// 			FnIndex:     0,
	// 			Data:        []interface{}{regexMatched[2], count},
	// 			SessionHash: sessionHash,
	// 		}
	// 		b, err := json.Marshal(r)
	// 		if err != nil {
	// 			ctx.SendChain(message.Text("Error:", err))
	// 			return
	// 		}
	// 		data, err := web.PostData(aipaintURL, "application/json", bytes.NewReader(b))
	// 		if err != nil {
	// 			ctx.SendChain(message.Text("Error:", err))
	// 			return
	// 		}
	// 		fmt.Println(string(data))
	// 		m := message.Message{}
	// 		gjson.ParseBytes(data).Get("data.0").ForEach(func(_, value gjson.Result) bool {
	// 			m = append(m, ctxext.FakeSenderForwardNode(ctx, message.Image(strings.ReplaceAll(value.String(), "data:image/png;base64,", "base64://"))))
	// 			return true
	// 		})
	// 		if id := ctx.Send(m).ID(); id == 0 {
	// 			ctx.SendChain(message.Text("ERROR: 可能被风控或下载图片用时过长，请耐心等待"))
	// 		}
	// 	})
}

// Package aipaint ai画图
package aipaint

import (
	"fmt"
	"net/url"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	mytoken    = "06LhgOew9PJDFQKnfdcSI3BtXz84AGoM"
	aipaintURL = "http://91.216.169.75:5010/got_image?token=%v&tags=%v"
	// aipaintURL  = "https://22229.gradio.app/api/predict/"
	// sessionHash = "zerobot"
)

// type request struct {
// 	FnIndex     int           `json:"fn_index"`
// 	Data        []interface{} `json:"data"`
// 	SessionHash string        `json:"session_hash"`
// }

func init() { // 插件主体
	engine := control.Register("aipaint", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "ai画图\n" +
			"ai画图xxx",
		// "- 画1张a photo of sks toy riding a bicycle",
		PrivateDataFolder: "aipaint",
	})

	engine.OnPrefix(`ai画图`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			m := message.Message{ctxext.FakeSenderForwardNode(ctx, message.Image(fmt.Sprintf(aipaintURL, mytoken, url.QueryEscape(strings.TrimSpace(strings.ReplaceAll(args, " ", "%20"))))).Add("cache", 0))}
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

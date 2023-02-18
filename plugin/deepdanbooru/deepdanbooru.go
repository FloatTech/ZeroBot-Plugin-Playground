// Package deepdanbooru 二次元标签分析
package deepdanbooru

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"strings"
	"time"

	hf "github.com/FloatTech/AnimeAPI/huggingface"
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	magicpromptURL = "wss://hysts-deepdanbooru.hf.space/queue/join"
)

func init() { // 插件主体
	engine := control.Register("deepdanbooru", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "二次元标签分析",
		Help:              "- 鉴赏图片 [xxx]",
		PrivateDataFolder: "deepdanbooru",
	})

	cachefolder := engine.DataFolder()
	// 开启
	engine.OnKeywordGroup([]string{"鉴赏图片"}, zero.MustProvidePicture).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("少女祈祷中..."))
			for _, url := range ctx.State["image_url"].([]string) {
				t, tseq, err := handleImage(url)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
				}
				digest := md5.Sum(binary.StringToBytes(url))
				f := cachefolder + hex.EncodeToString(digest[:])
				if file.IsNotExist(f) {
					_ = imgfactory.SavePNG2Path(f, t)
				}
				m := message.Message{ctxext.FakeSenderForwardNode(ctx, message.Image("file:///"+file.BOTPATH+"/"+f)),
					ctxext.FakeSenderForwardNode(ctx, message.Text("tags: ", strings.Join(tseq, ",")))}
				if id := ctx.Send(m).ID(); id == 0 {
					ctx.SendChain(message.Text("ERROR: 可能被风控或下载图片用时过长，请耐心等待"))
				}
			}
		})
}

func handleImage(url string) (im image.Image, tseq []string, err error) {
	_ctx, _cancel := context.WithTimeout(context.Background(), hf.TimeoutMax*time.Second)
	defer _cancel()
	var (
		data []byte
		img  image.Image
	)
	c, _, err := websocket.DefaultDialer.Dial(magicpromptURL, nil)
	if err != nil {
		return
	}
	defer c.Close()

	imageData, err := web.GetData(url)
	if err != nil {
		return
	}
	r := hf.PushRequest{
		FnIndex:     0,
		SessionHash: "zerobot",
	}
	b, err := json.Marshal(r)
	if err != nil {
		return
	}

	encodeStr := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imageData)
	r1 := hf.PushRequest{
		FnIndex:     0,
		Data:        []interface{}{encodeStr, 0.5},
		SessionHash: "zerobot",
	}
	b1, err := json.Marshal(r1)
	if err != nil {
		return
	}
	t := time.NewTicker(time.Second * 1)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			_, data, err = c.ReadMessage()
			if err != nil {
				return
			}
			j := gjson.ParseBytes(data)
			switch j.Get("msg").String() {
			case "send_hash":
				err = c.WriteMessage(websocket.TextMessage, b)
				if err != nil {
					return
				}
			case "send_data":
				err = c.WriteMessage(websocket.TextMessage, b1)
				if err != nil {
					return
				}
			case hf.WssCompleteStatus:
				arr := j.Get("output.data.0.confidences").Array()

				img, _, err = image.Decode(bytes.NewReader(imageData))
				if err != nil {
					return
				}

				img = imgfactory.Limit(img, 1280, 720)

				canvas := gg.NewContext(img.Bounds().Size().X, img.Bounds().Size().Y+int(float64(img.Bounds().Size().X)*0.2)+len(arr)*img.Bounds().Size().X/25)
				canvas.SetRGB(1, 1, 1)
				canvas.Clear()
				canvas.DrawImage(img, 0, 0)
				if err = canvas.LoadFontFace(text.BoldFontFile, float64(img.Bounds().Size().X)*0.1); err != nil {
					return
				}
				canvas.SetRGB(0, 0, 0)
				canvas.DrawString("", float64(img.Bounds().Size().X)*0.02, float64(img.Bounds().Size().Y)+float64(img.Bounds().Size().X)*0.1)
				if err = canvas.LoadFontFace(text.ConsolasFontFile, float64(img.Bounds().Size().X)*0.04); err != nil {
					return
				}

				longestlen := 0
				for _, v := range arr {
					if len(v.Get("label").String()) > longestlen {
						longestlen = len(v.Get("label").String())
					}
					tseq = append(tseq, v.Get("label").String())
				}
				longestlen++

				i := float64(img.Bounds().Size().Y) + float64(img.Bounds().Size().X)*0.2
				rate := float64(img.Bounds().Size().X) * 0.04
				for _, v := range arr {
					canvas.DrawString(fmt.Sprintf("* %-*s -%.3f-", longestlen, v.Get("label").String(), v.Get("confidence").Float()), float64(img.Bounds().Size().X)*0.04, i)
					i += rate
				}
				im = canvas.Image()
				return
			}

		case <-_ctx.Done():
			err = errors.New("ERROR: 吟唱提示指令超时")
			return
		}
	}
}

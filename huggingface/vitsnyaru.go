package huggingface

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	vitsnyaruRepo = "/innnky/vits-nyaru"
)

func init() { // 插件主体
	engine := control.Register("vitsnyaru", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "猫雷vits\n" +
			"- 让猫雷说 xxx",
		PrivateDataFolder: "vitsnyaru",
	})

	// 开启
	engine.OnPrefix(`让猫雷说`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			pushURL := embed + vitsnyaruRepo + pushPath
			statusURL := embed + vitsnyaruRepo + statusPath
			ctx.SendChain(message.Text("少女祈祷中..."))
			var (
				pushReq   pushRequest
				pushRes   pushResponse
				statusReq statusRequest
				statusRes statusResponse
				data      []byte
			)
			ch := make(chan []byte, 1)
			_ctx, _cancel := context.WithTimeout(context.Background(), timeoutMax*time.Second)
			defer _cancel()
			// 获取clean后的文本
			pushReq = pushRequest{
				Action:      defaultAction,
				Data:        []interface{}{args},
				FnIndex:     1,
				SessionHash: defaultSessionHash,
			}
			pushRes, err := push(pushURL, pushReq)
			if err != nil {
				ctx.SendChain(message.Text("Error:", err))
				return
			}
			statusReq = statusRequest{
				Hash: pushRes.Hash,
			}
			go func(c context.Context) {
				t := time.NewTicker(time.Second * 1)
				defer t.Stop()
			LOOP:
				for {
					select {
					case <-t.C:
						data, err = status(statusURL, statusReq)
						if err != nil {
							ctx.SendChain(message.Text("Error:", err))
							break LOOP
						}
						if gjson.ParseBytes(data).Get("status").String() == completeStatus {
							ch <- data
							break LOOP
						}
					case <-c.Done():
						break LOOP
					}
				}
			}(_ctx)
			data = <-ch
			err = json.Unmarshal(data, &statusRes)
			if err != nil {
				ctx.SendChain(message.Text("Error:", err))
				return
			}
			// 用clean的文本预测语音
			pushReq = pushRequest{
				Action:      defaultAction,
				Data:        statusRes.Data.Data,
				FnIndex:     2,
				SessionHash: defaultSessionHash,
			}
			pushRes, err = push(pushURL, pushReq)
			if err != nil {
				ctx.SendChain(message.Text("Error:", err))
				return
			}
			statusReq = statusRequest{
				Hash: pushRes.Hash,
			}
			go func(c context.Context) {
				t := time.NewTicker(time.Second * 1)
				defer t.Stop()
			LOOP:
				for {
					select {
					case <-t.C:
						data, err = status(statusURL, statusReq)
						if err != nil {
							ctx.SendChain(message.Text("Error:", err))
							return
						}
						if gjson.ParseBytes(data).Get("status").String() == completeStatus {
							ch <- data
							break LOOP
						}
					case <-c.Done():
						break LOOP
					}
				}
			}(_ctx)
			data = <-ch
			err = json.Unmarshal(data, &statusRes)
			if err != nil {
				ctx.SendChain(message.Text("Error:", err))
				return
			}
			fmt.Printf("%#v\n", statusRes)
			// 发送语音
			if len(statusRes.Data.Data) < 2 {
				ctx.SendChain(message.Text("Error: 未能获取语音"))
				return
			}
			ctx.SendChain(message.Record("base64://" + strings.TrimPrefix(statusRes.Data.Data[1].(string), "data:audio/wav;base64,")))
		})
}

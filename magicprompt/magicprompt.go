// Package magicprompt 魔力提词
package magicprompt

import (
	"encoding/json"
	"strings"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const magicpromptURL = "wss://spaces.huggingface.tech/Gustavosta/MagicPrompt-Stable-Diffusion/queue/join"

type hfRequest struct {
	Action      string        `json:"action,omitempty"`
	FnIndex     int           `json:"fn_index"`
	Data        []interface{} `json:"data"`
	SessionHash string        `json:"session_hash"`
}

func init() { // 插件主体
	engine := control.Register("magicprompt", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "MagicPrompt-Stable-Diffusion\n" +
			"- 魔力提词 xxx",
		PrivateDataFolder: "magicprompt",
	})

	// 开启
	engine.OnPrefix(`魔力提词`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := ctx.State["args"].(string)
			c, _, err := websocket.DefaultDialer.Dial(magicpromptURL, nil)
			if err != nil {
				ctx.SendChain(message.Text("Error:", err))
				return
			}
			defer c.Close()

			done := make(chan struct{})

			go func() {
				defer close(done)
				for {
					_, m, err := c.ReadMessage()
					if err != nil {
						logrus.Println("read:", err)
						return
					}
					text := gjson.ParseBytes(m).Get("output.data.0").String()
					if text != "" {
						m := message.Message{}
						for _, v := range strings.Split(text, "\n\n") {
							m = append(m, ctxext.FakeSenderForwardNode(ctx, message.Text(v)))
						}
						if id := ctx.Send(m).ID(); id == 0 {
							ctx.SendChain(message.Text("ERROR: 可能被风控或下载图片用时过长，请耐心等待"))
						}
					}
					logrus.Printf("recv: %s", m)
				}
			}()

			r := hfRequest{
				FnIndex:     0,
				Data:        []interface{}{args},
				SessionHash: "zerobot",
			}
			b, err := json.Marshal(r)
			if err != nil {
				ctx.SendChain(message.Text("Error:", err))
				return
			}

			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					err := c.WriteMessage(websocket.TextMessage, b)
					if err != nil {
						ctx.SendChain(message.Text("Error:", err))
						return
					}
				}
			}
		})
}

package huggingface

import (
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
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

			// 获取clean后的文本
			pushReq := pushRequest{
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
			statusReq := statusRequest{
				Hash: pushRes.Hash,
			}
			statusRes, err := status(statusURL, statusReq)
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
			statusRes, err = status(statusURL, statusReq)
			if err != nil {
				ctx.SendChain(message.Text("Error:", err))
				return
			}

			// 发送语音
			if len(statusRes.Data.Data) < 2 {
				ctx.SendChain(message.Text("Error: 未能获取语音"))
				return
			}
			ctx.SendChain(message.Record("base64://" + strings.TrimPrefix(statusRes.Data.Data[1].(string), "data:audio/wav;base64,")))
		})
}

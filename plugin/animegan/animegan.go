// Package animegan 二次元画风
package animegan

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	envURL        = "https://hf.space/embed/akhaliq/AnimeGANv2/api/queue"
	pushURL       = envURL + "/push/"
	statusURL     = envURL + "/status/"
	sessionHash   = "zerobot"
	predictAction = "predict"
	version1      = "version 1 (🔺 stylization, 🔻 robustness)"
	version2      = "version 2 (🔺 robustness,🔻 stylization)"
)

var (
	datapath string
)

type hfRequest struct {
	Action      string        `json:"action"`
	FnIndex     int           `json:"fn_index"`
	Data        []interface{} `json:"data"`
	SessionHash string        `json:"session_hash"`
}

func init() { // 插件主体
	engine := control.Register("animegan", &ctrl.Options[*zero.Ctx]{
		Brief:            "二次元画风",
		DisableOnDefault: false,
		Help: "- 二次元画风1 [图片]\n" +
			"- 二次元画风2 [图片]",
		PrivateDataFolder: "animegan",
	})
	datapath = file.BOTPATH + "/" + engine.DataFolder()
	// 开启
	engine.OnRegex(`^(二次元画风)[\s\S]*?(\[CQ:(image\,file=([0-9a-zA-Z]{32}).*|at.+?(\d{5,11}))\].*|(\d+))$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			// server, token, err := cfg.load()
			// if err != nil {
			// 	ctx.SendChain(message.Text("ERROR: ", err))
			// 	return
			// }
			c := newContext(ctx.Event.UserID)
			list := ctx.State["regex_matched"].([]string)
			err := c.prepareLogos(list[4]+list[5]+list[6], strconv.FormatInt(ctx.Event.UserID, 10))
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			args := strings.TrimSuffix(strings.TrimPrefix(list[0], list[1]), list[2])
			if args == "" {
				args = version1
			}
			if strings.TrimSpace(args) == "1" {
				args = version1
			} else {
				args = version2
			}
			ctx.SendChain(message.Text("少女祈祷中..."))
			data, err := os.ReadFile(c.headimgsdir[0])
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			hs, err := pushData(data, args)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			imgBase64, err := statusData(hs)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			m := message.Message{ctxext.FakeSenderForwardNode(ctx, message.Image("base64://"+strings.TrimPrefix(imgBase64, "data:image/png;base64,")))}
			if id := ctx.Send(m).ID(); id == 0 {
				ctx.SendChain(message.Text("ERROR: 可能被风控或下载图片用时过长，请耐心等待"))
			}
		})
}

func pushData(data []byte, version string) (hash string, err error) {
	encodeStr := base64.StdEncoding.EncodeToString(data)
	encodeStr = "data:image/jpeg;base64," + encodeStr
	r := hfRequest{
		Action:      predictAction,
		FnIndex:     0,
		Data:        []interface{}{encodeStr, version},
		SessionHash: sessionHash,
	}
	b, err := json.Marshal(r)
	if err != nil {
		return
	}
	data, err = web.PostData(pushURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return
	}
	time.Sleep(2 * time.Second)
	hash = gjson.ParseBytes(data).Get("hash").String()
	return
}

func statusData(hash string) (imgBase64 string, err error) {
	data, err := web.PostData(statusURL, "application/json", strings.NewReader(fmt.Sprintf(`{"hash": "%v"}`, hash)))
	if err != nil {
		return
	}
	imgBase64 = gjson.ParseBytes(data).Get("data.data.0").String()
	return
}

// Package youdaotranslate 有道翻译
package youdaotranslate

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	youdaotlURL = "https://fanyi.youdao.com/translate?doctype=json&type=%v&i=%v"
)

type youdao struct {
	LType           string `json:"type"`
	ErrorCode       int    `json:"errorCode"`
	ElapsedTime     int    `json:"elapsedTime"`
	TranslateResult [][]struct {
		Src string `json:"src"`
		Tgt string `json:"tgt"`
	} `json:"translateResult"`
}

func init() {
	en := control.Register("youdaotranslate", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "youdaotranslate\n" +
			"- 有道(英语|日语|韩语)翻译\n",
	})
	en.OnMessage(func(ctx *zero.Ctx) bool {
		msg := ctx.Event.Message
		if msg[0].Type != "reply" {
			return false
		}
		for _, elem := range msg {
			if elem.Type == "text" {
				text := elem.Data["text"]
				text = strings.ReplaceAll(text, " ", "")
				text = strings.ReplaceAll(text, "\r", "")
				text = strings.ReplaceAll(text, "\n", "")
				switch text {
				case "有道翻译":
					ctx.State["tl_type"] = "AUTO"
					return true
				case "有道英语翻译":
					ctx.State["tl_type"] = "ZH_CN2EN"
					return true
				case "有道日语翻译":
					ctx.State["tl_type"] = "ZH_CN2JA"
					return true
				case "有道韩语翻译":
					ctx.State["tl_type"] = "ZH_CN2KR"
					return true
				}
			}
		}
		return false
	}).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			text := ctx.GetMessage(message.NewMessageIDFromString(ctx.Event.Message[0].Data["id"])).Elements[0].Data["text"]
			tlType := ctx.State["tl_type"].(string)
			tgt, err := getYoudao(text, tlType)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.ReplyWithMessage(ctx.Event.Message[0].Data["id"], message.Text(tgt))...)
		})
	en.OnRegex(`^有道(.{0,2})翻译\s?(.{1,200})`).
		Handle(func(ctx *zero.Ctx) {
			var (
				tlType string
				text   string
			)
			switch ctx.State["regex_matched"].([]string)[1] {
			case "英语":
				tlType = "ZH_CN2EN"
			case "日语":
				tlType = "ZH_CN2JA"
			case "韩语":
				tlType = "ZH_CN2KR"
			default:
				tlType = "AUTO"
			}
			text = ctx.State["regex_matched"].([]string)[2]
			tgt, err := getYoudao(text, tlType)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.ReplyWithMessage(ctx.Event.MessageID, message.Text(tgt))...)
		})
}

func getYoudao(text string, tlType string) (tgt string, err error) {
	var yd youdao
	data, err := web.RequestDataWith(web.NewDefaultClient(), fmt.Sprintf(youdaotlURL, tlType, text), "GET", "", web.RandUA())
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &yd)
	if err != nil {
		return
	}
	if len(yd.TranslateResult) > 0 && len(yd.TranslateResult[0]) > 0 {
		tgt = yd.TranslateResult[0][0].Tgt
	} else {
		err = errors.New("数据为空")
	}
	return
}

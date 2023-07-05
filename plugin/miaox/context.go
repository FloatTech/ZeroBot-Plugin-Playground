package aimgr

import (
	"context"
	"errors"
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/repo"
	pTypes "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/types"
	"github.com/bincooo/MiaoX/types"
	"github.com/bincooo/MiaoX/vars"
	wapi "github.com/bincooo/openai-wapi"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"strconv"
	"strings"
	"time"
)

var (
	contextStore = make(map[string]types.ConversationContext)
)

func deleteConversationContext(ctx *zero.Ctx) {
	var id int64 = 0
	if ctx.Event.GroupID == 0 {
		id = ctx.Event.UserID
	} else {
		id = ctx.Event.GroupID
	}
	key := strconv.FormatInt(id, 10)
	delete(contextStore, key)
}

func updateConversationContext(cctx types.ConversationContext) {
	contextStore[cctx.Id] = cctx
}

func createConversationContext(ctx *zero.Ctx, bot string) (types.ConversationContext, error) {
	var id int64 = 0
	if ctx.Event.GroupID == 0 {
		id = ctx.Event.UserID
	} else {
		id = ctx.Event.GroupID
	}

	key := strconv.FormatInt(id, 10)

	if cctx, ok := contextStore[key]; ok {
		return cctx, nil
	}

	global := repo.GetGlobal()
	if bot == "" {
		bot = global.Bot
	}

	model := ""
	if strings.HasPrefix(bot, vars.Bing) {
		expr := bot[len(bot)-1:]
		switch expr {
		case "b":
			model = "Balanced"
		case "p":
			model = "Precise"
		case "s":
			model = "Sydney"
		default:
			model = "Creative"
		}
		bot = vars.Bing
	}

	cctx := types.ConversationContext{
		Id:        key,
		Bot:       bot,
		MaxTokens: global.MaxTokens,
		Chain:     "replace,cache,args,",
		Model:     model,
		Data:      pTypes.ConversationContextArgs{},
	}

	tokens, err := repo.FindTokens(bot)
	if err != nil {
		return cctx, errors.New("查询凭证失败, 请先添加`" + bot + "`凭证")
	}
	if len(tokens) == 0 {
		return cctx, errors.New("无可用的凭证")
	}

	if bot == vars.OpenAIAPI {
		// 检查余额
		if e := checkApiOpenai(*tokens[0], global.Proxy); e != nil {
			return cctx, e
		}
	}

	if bot == vars.OpenAIWeb {
		// 检查失效
		cctx.BaseURL = "https://ai.fakeopen.com/api"
		if err := checkWebOpenai(tokens, global.Proxy); err != nil {
			return cctx, err
		}

		// 为空，尝试登陆
		if tokens[0].Token == "" {
			if err := loginWebOpenai(tokens, global); err != nil {
				return cctx, err
			}
		}
	}

	if bot == vars.Claude {
		cctx.AppId = tokens[0].AppId
	}

	if bot == vars.Bing {
		cctx.BaseURL = global.NbServ
	}

	// 默认预设
	if global.Preset != "" {
		preset := repo.GetPresetScene(global.Preset)
		if preset == nil {
			logrus.Warn("预设`", global.Preset, "`不存在")
		} else if preset.Type != bot {
			logrus.Warn("预设`", global.Preset, "`类型不匹配, 需要（", bot, "）实际为（", preset.Type, "）")
		} else {
			cctx.Preset = preset.Content
			cctx.Format = preset.Message
			if preset.Chain != "" {
				cctx.Chain += preset.Chain
			}
		}
	}

	cctx.Token = tokens[0].Token
	contextStore[key] = cctx
	return cctx, nil
}

// 登陆网页版
func loginWebOpenai(tokens []*repo.Token, global repo.Global) error {
	token, err := wapi.WebLogin(tokens[0].Email, tokens[0].Passwd, global.Proxy)
	if err != nil {
		return errors.New("OpenAI WEB `" + tokens[0].Key + "`登陆失败: " + err.Error())
	}
	tokens[0].Token = token
	tokens[0].Expire = time.Now().Add(15 * 24 * time.Hour).Format("2006-01-02 15:04:05")
	repo.UpdateToken(*tokens[0])
	return nil
}

// 检查余额
func checkApiOpenai(token repo.Token, proxy string) error {
	if billing, _ := wapi.Query(context.Background(), token.Token, proxy); billing == nil || billing.System-billing.Soft < 0 {
		return errors.New("Err: `" + token.Key + "`凭证余额为0")
	}
	return nil
}

// 检查过期时间
func checkWebOpenai(tokens []*repo.Token, proxy string) error {
	if tokens[0].Expire != "" && tokens[0].Expire != "-1" {
		expire, err := time.Parse("2006-01-02 15:04:05", tokens[0].Expire)
		if err != nil {
			return errors.New("warning： `" + tokens[0].Expire + "`过期日期解析有误")
		}

		if expire.Before(time.Now()) {
			// 已过期
			token, err := wapi.WebLogin(tokens[0].Email, tokens[0].Passwd, proxy)
			if err != nil {
				return errors.New("OpenAI WEB `" + tokens[0].Key + "`登陆失败: " + err.Error())
			}
			tokens[0].Token = token
			tokens[0].Expire = time.Now().Add(15 * 24 * time.Hour).Format("2006-01-02 15:04:05")
			repo.UpdateToken(*tokens[0])
		}
	}
	return nil
}

func parseMessage(ctx *zero.Ctx) string {
	// and more...
	return ctx.ExtractPlainText()
}

package chain

import (
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/repo/store"
	pTypes "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/types"
	"github.com/bincooo/MiaoX/types"
	"strings"
)

const MaxOnlineCount = 30

type OnlineInterceptor struct {
	types.BaseInterceptor
}

func (c *OnlineInterceptor) Before(bot *types.Bot, ctx *types.ConversationContext) bool {
	cacheOnline(ctx)
	if strings.Contains(ctx.Prompt, "[online]") {
		online := ""
		for _, o := range store.GetOnline(ctx.Id) {
			online += fmt.Sprintf(`{"qq":"%s", "name": "%s"}`, o["id"], o["name"])
		}
		ctx.Prompt = strings.Replace(ctx.Prompt, "[online]", "["+online+"]", -1)
	}

	return true
}

func cacheOnline(ctx *types.ConversationContext) {
	online := store.GetOnline(ctx.Id)
	args := ctx.Data.(pTypes.ConversationContextArgs)
	// 如果已在线列表中，先删除后加入到结尾
	for i, ol := range online {
		if ol["id"] != args.Current {
			continue
		}

		if len(online) == 1 {
			online = make([]map[string]string, 0)
		} else {
			online = append(online[:i], online[i+1:]...)
		}

		break
	}

	// 加入在线列表
	online = append(online, map[string]string{
		"id":   args.Current,
		"name": args.Nickname,
	})

	// 控制最大在线人数
	if len(online) > MaxOnlineCount {
		online = online[len(online)-MaxOnlineCount:]
	}
	store.CacheOnline(ctx.Id, online)
}

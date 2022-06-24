package main

import (
	"math/rand"
	"time"

	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/example/JiangRed"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/github/yydz"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/nihongo"

	// 程序主体
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	rand.Seed(time.Now().Unix()) // 全局 seed, 插件无需再 seed

	zero.OnCommand("hello").
		Handle(func(ctx *zero.Ctx) {
			ctx.Send("world")
		})

	zero.RunAndBlock(zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{123456},
		Driver: []zero.Driver{
			driver.NewWebSocketClient("ws://127.0.0.1:6700", "access_token"),
		},
	}, nil)
}

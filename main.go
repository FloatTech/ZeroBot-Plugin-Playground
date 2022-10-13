// Package main ZeroBot-Plugin-Plugin main file
package main

import (
	"encoding/json"
	"flag"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/process"
	"github.com/sirupsen/logrus"

	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/aipaint"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/bottle"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/dailynews"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/deepdanbooru"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/example/JiangRed"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/example/xiaoguofan"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/fadian"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/fyzhq"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/jptingroom"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/kfccrazythursday"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/magicprompt"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/moegozh"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/ottoshout"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/partygame"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/playwright"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/qqci"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/qqclean"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/recall" // SGK2401
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/tgyj"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/vote"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/youdaotranslate"

	// 程序主体
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

type zbpcfg struct {
	Z zero.Config        `json:"zero"`
	W []*driver.WSClient `json:"ws"`
}

var config zbpcfg

func init() {
	sus := make([]int64, 0, 16)
	// 直接写死 AccessToken 时，请更改下面第二个参数
	token := flag.String("t", "", "Set AccessToken of WSClient.")
	// 直接写死 URL 时，请更改下面第二个参数
	url := flag.String("u", "ws://127.0.0.1:6700", "Set Url of WSClient.")
	// 默认昵称
	adana := flag.String("n", "椛椛", "Set default nickname.")
	prefix := flag.String("p", "/", "Set command prefix.")
	runcfg := flag.String("c", "", "Run from config file.")

	flag.Parse()

	for _, s := range flag.Args() {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		sus = append(sus, i)
	}
	// sus = append(sus, 123456)
	if *runcfg != "" {
		f, err := os.Open(*runcfg)
		if err != nil {
			panic(err)
		}
		config.W = make([]*driver.WSClient, 0, 2)
		err = json.NewDecoder(f).Decode(&config)
		f.Close()
		if err != nil {
			panic(err)
		}
		config.Z.Driver = make([]zero.Driver, len(config.W))
		for i, w := range config.W {
			config.Z.Driver[i] = w
		}
		logrus.Infoln("[main] 从", *runcfg, "读取配置文件")
		return
	}

	config.W = []*driver.WSClient{driver.NewWebSocketClient(*url, *token)}
	config.Z = zero.Config{
		NickName:      append([]string{*adana}, "ATRI", "atri", "亚托莉", "アトリ"),
		CommandPrefix: *prefix,
		SuperUsers:    sus,
		Driver:        []zero.Driver{config.W[0]},
	}
}
func main() {
	rand.Seed(time.Now().Unix()) // 全局 seed, 插件无需再 seed

	zero.OnCommand("hello").
		Handle(func(ctx *zero.Ctx) {
			ctx.Send("world")
		})

	zero.RunAndBlock(config.Z, process.GlobalInitMutex.Unlock)
}

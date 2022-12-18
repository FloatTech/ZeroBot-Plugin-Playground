// Package main ZeroBot-Plugin-Plugin main file
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/process"
	"github.com/sirupsen/logrus"

	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/animegan"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/chinesebqb"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/dailynews"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/deepdanbooru"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/example/JiangRed"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/example/xiaoguofan"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/fadian"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/fyzhq"

	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/games"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/hitokoto"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/kfccrazythursday"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/moegozh"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/ottoshout"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/partygame"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/playwright"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/qqci"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/qqclean"
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/recall" // SGK2401
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/subweibo"
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
	// 解析命令行参数
	d := flag.Bool("d", false, "Enable debug level log and higher.")
	w := flag.Bool("w", false, "Enable warning level log and higher.")
	h := flag.Bool("h", false, "Display this help.")
	// 直接写死 AccessToken 时，请更改下面第二个参数
	token := flag.String("t", "", "Set AccessToken of WSClient.")
	// 直接写死 URL 时，请更改下面第二个参数
	url := flag.String("u", "ws://127.0.0.1:6700", "Set Url of WSClient.")
	// 默认昵称
	adana := flag.String("n", "椛椛", "Set default nickname.")
	prefix := flag.String("p", "/", "Set command prefix.")
	runcfg := flag.String("c", "", "Run from config file.")
	save := flag.String("s", "", "Save default config to file and exit.")
	late := flag.Uint("l", 233, "Response latency (ms).")
	rsz := flag.Uint("r", 4096, "Receiving buffer ring size.")
	maxpt := flag.Uint("x", 4, "Max process time (min).")

	flag.Parse()

	if *h {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	} else {
		if *d && !*w {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if *w {
			logrus.SetLevel(logrus.WarnLevel)
		}
	}

	for _, s := range flag.Args() {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		sus = append(sus, i)
	}

	// 通过代码写死的方式添加主人账号
	// sus = append(sus, 12345678)
	// sus = append(sus, 87654321)

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
		NickName:       append([]string{*adana}, "ATRI", "atri", "亚托莉", "アトリ"),
		CommandPrefix:  *prefix,
		SuperUsers:     sus,
		RingLen:        *rsz,
		Latency:        time.Duration(*late) * time.Millisecond,
		MaxProcessTime: time.Duration(*maxpt) * time.Minute,
		Driver:         []zero.Driver{config.W[0]},
	}

	if *save != "" {
		f, err := os.Create(*save)
		if err != nil {
			panic(err)
		}
		err = json.NewEncoder(f).Encode(&config)
		f.Close()
		if err != nil {
			panic(err)
		}
		logrus.Infoln("[main] 配置文件已保存到", *save)
		os.Exit(0)
	}
}
func main() {
	rand.Seed(time.Now().Unix()) // 全局 seed, 插件无需再 seed

	zero.OnCommand("hello").
		Handle(func(ctx *zero.Ctx) {
			ctx.Send("world")
		})

	zero.RunAndBlock(&config.Z, process.GlobalInitMutex.Unlock)
}

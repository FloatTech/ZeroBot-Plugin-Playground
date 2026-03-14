// Package main ZeroBot-Plugin-Playground main file
package main

//go:generate go run github.com/FloatTech/ZeroBot-Plugin-Playground/abineundo/ref -r .

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/process"
	"github.com/sirupsen/logrus"

	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/abineundo" // 设置插件优先级&更改控制台属性ssss

	// ---------以下插件均可通过前面加 // 注释，注释后停用并不加载插件--------- //
	// ----------------------插件优先级按顺序从高到低---------------------- //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	// ----------------------------高优先级区---------------------------- //
	// vvvvvvvvvvvvvvvvvvvvvvvvvvvv高优先级区vvvvvvvvvvvvvvvvvvvvvvvvvvvv //
	//               vvvvvvvvvvvvvv高优先级区vvvvvvvvvvvvvv               //
	//                      vvvvvvv高优先级区vvvvvvv                      //
	//                          vvvvvvvvvvvvvv                          //
	//                               vvvv                               //

	//                               ^^^^                               //
	//                          ^^^^^^^^^^^^^^                          //
	//                      ^^^^^^^高优先级区^^^^^^^                      //
	//               ^^^^^^^^^^^^^^高优先级区^^^^^^^^^^^^^^               //
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^高优先级区^^^^^^^^^^^^^^^^^^^^^^^^^^^^ //
	// ----------------------------高优先级区---------------------------- //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	// ----------------------------中优先级区---------------------------- //
	// vvvvvvvvvvvvvvvvvvvvvvvvvvvv中优先级区vvvvvvvvvvvvvvvvvvvvvvvvvvvv //
	//               vvvvvvvvvvvvvv中优先级区vvvvvvvvvvvvvv               //
	//                      vvvvvvv中优先级区vvvvvvv                      //
	//                          vvvvvvvvvvvvvv                          //
	//                               vvvv                               //

	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/anotherfortune"  // 另一个今日人品
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/chatgpt"         // ChatGPT对话
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/chinesebqb"      // 表情包
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/cybercat"        // 云养猫
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/draw"            // 服务详情
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/exJiangRed"      // 江林的插件编写教学
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/exXiaoGuoFan"    // 小锅饭的示例
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/fadian"          // 每日发癫
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/fgopickup"       // FGO卡池查询
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/fyzhq"           // 发音转换器
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/games"           // 游戏系统
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/klala"           // 星穹铁道面板/图鉴查询
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/kokomi"          // 原神面板查询
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/moegozh"         // VITS模型拟声迫真中文
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/movies"          // 电影查询
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/ottoshout"       // otto鬼叫
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/partygame"       // 轮盘赌
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/playwright"      // 网页截图
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/qqci"            // 简易CI/CD
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/qqclean"         // QQ号清理
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/recall"          // 回复撤回
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/rsshub"          // RssHub订阅
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/slash"           // Slash自交互
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/subweibo"        // 订阅微博消息
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/tgyj"            // 同归于尽
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/vote"            // 投票
	_ "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/youdaotranslate" // 有道翻译

	// _ "github.com/FloatTech/ZeroBot-Plugin/plugin/wtf"           // 鬼东西

	//                               ^^^^                               //
	//                          ^^^^^^^^^^^^^^                          //
	//                      ^^^^^^^中优先级区^^^^^^^                      //
	//               ^^^^^^^^^^^^^^中优先级区^^^^^^^^^^^^^^               //
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^中优先级区^^^^^^^^^^^^^^^^^^^^^^^^^^^^ //
	// ----------------------------中优先级区---------------------------- //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	// ----------------------------低优先级区---------------------------- //
	// vvvvvvvvvvvvvvvvvvvvvvvvvvvv低优先级区vvvvvvvvvvvvvvvvvvvvvvvvvvvv //
	//               vvvvvvvvvvvvvv低优先级区vvvvvvvvvvvvvv               //
	//                      vvvvvvv低优先级区vvvvvvv                      //
	//                          vvvvvvvvvvvvvv                          //
	//                               vvvv                               //

	//                               ^^^^                               //
	//                          ^^^^^^^^^^^^^^                          //
	//                      ^^^^^^^低优先级区^^^^^^^                      //
	//               ^^^^^^^^^^^^^^低优先级区^^^^^^^^^^^^^^               //
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^低优先级区^^^^^^^^^^^^^^^^^^^^^^^^^^^^ //
	// ----------------------------低优先级区---------------------------- //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	// -----------------------以下为内置依赖，勿动------------------------ //
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	// -----------------------以上为内置依赖，勿动------------------------ //
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
	}
	if *d && !*w {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if *w {
		logrus.SetLevel(logrus.WarnLevel)
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
	if !strings.Contains(runtime.Version(), "go1.2") { // go1.20之前版本需要全局 seed，其他插件无需再 seed
		rand.Seed(time.Now().UnixNano()) //nolint:staticcheck
	}

	zero.OnCommand("hello").
		Handle(func(ctx *zero.Ctx) {
			ctx.Send("world")
		})

	zero.RunAndBlock(&config.Z, process.GlobalInitMutex.Unlock)
}

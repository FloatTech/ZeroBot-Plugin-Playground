// Package anotherfortune 简单的测人品 仿照的是鱼子酱的www
package anotherfortune

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"math/rand"
	"strconv"
	"time"

	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type card struct {
	Name     string `json:"name"`
	Cardtype string `json:"cardtype"`
	Info     struct {
		Description        string `json:"description"`
		ReverseDescription string `json:"reverseDescription"`
	} `json:"info"`
}

type cardset = map[string]card

var (
	jrrpbk   string
	info     string
	uptime   string
	vme50    string
	cardMap  = make(cardset, 256)
	reasons  = []string{" | "}
	position = []string{"正位", "逆位"}
	result   map[int64](int)
	egg      map[string](int)
	signTF   map[string](int)
)

func init() {
	engine := control.Register("anotherfortune", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "另一个今日人品",
		Help:              "仿照鱼子酱的今日人品w",
		PrivateDataFolder: "anotherfortune",
	})
	signTF = make(map[string](int))
	egg = make(map[string](int))
	result = make(map[int64](int))
	getTarot := fcext.DoOnceOnSuccess(
		func(ctx *zero.Ctx) bool { // 检查 塔罗牌文件是否存在
			data, err := os.ReadFile(engine.DataFolder() + "tarots.json")
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return false
			}
			err = json.Unmarshal(data, &cardMap)
			if err != nil {
				panic(err)
			}
			return true
		},
	)

	engine.OnFullMatch("今日人品", getTarot).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var mutex sync.RWMutex // 添加读写锁以保证稳定性
			mutex.Lock()
			p := rand.Intn(2)
			i := rand.Intn(78)
			card := cardMap[(strconv.Itoa(i))]
			name := card.Name
			cardtype := card.Cardtype
			if p == 0 {
				info = card.Info.Description
			} else {
				info = card.Info.ReverseDescription
			} // 塔罗牌生成 (随机的)
			// 写的非常恶心 建议有时间赶紧重构x awa
			user := ctx.Event.UserID
			userS := strconv.FormatInt(user, 10)
			now := time.Now().Format("20060102")
			randEveryone := fcext.RandSenderPerDayN(ctx.Event.UserID, 100)
			si := now + userS // 合成
			dyn := time.Now().Hour()
			weeks := time.Now().Weekday()
			switch {
			case dyn <= 6 && dyn >= 0:
				uptime = "凌晨好~还没有睡觉呢~再不睡觉的话咱把你敲晕~" // 计算是早上还是晚上
			case dyn <= 11 && dyn > 6:
				uptime = "上午好~~是个笨蛋(bushi)~"
			case dyn <= 14 && dyn > 11:
				uptime = "中午好~吃饭了嘛w 如果没有快去吃饭哦w"
			case dyn <= 18 && dyn > 14:
				uptime = "下午好ww~咱很高兴看到你精力充沛的样子w"
			case dyn <= 24 && dyn > 18:
				uptime = "晚上好吖w~今天过的开心嘛ww"
			}
			if weeks.String() == "Thursday" {
				vme50 = "今天是疯狂星期四 v我50好嘛 www"
			} else {
				vme50 = ""
			}
			// CTRL C + CTRL V
			if signTF[si] == 0 {
				signTF[si] = (1)
				result[user] = (randEveryone)
				botName := zero.BotConfig.NickName[0]
				switch {
				case result[user] <= 20:
					jrrpbk = "[小凶]\n#" + botName + "抱了抱你~"
				case result[user] > 20 && result[user] < 50:
					jrrpbk = "[小吉]\n#" + botName + "偷瞄瞄~w"
				case result[user] >= 50 && result[user] < 90:
					jrrpbk = "[中吉]\n#" + botName + "捏了捏你的脸"
				case result[user] >= 90 && result[user] < 100:
					jrrpbk = "[吉]\n#" + botName + "摸了摸你的脸"
				case result[user] == 100:
					jrrpbk = "[大吉]\n#好诶~" + botName + "给你递了张彩票"
				}
				ctx.SendChain(message.At(user),
					message.Text(fmt.Sprintf("\n%s\n"+botName+"正在帮你整理~\n", uptime)),
					message.Text("今日的人品值为", result[user]),
					message.Text(jrrpbk),
					message.Text("今日塔罗牌是: \n归类于", cardtype, reasons[rand.Intn(len(reasons))], position[p], " 的 ", name, "\n"),
					message.Text("\n其意义为：\n", info, "\n", vme50))
			} else {
				ctx.SendChain(message.At(user), message.Text(" 今天已经测过了哦~今日的人品值为", result[user], "呢~"))
			}
			mutex.Unlock()
			// special time !
			if result[user] >= 90 && result[user] < 100 && egg[si] == 0 {
				egg[si] = (1)
				img, err := web.GetData("https://api.lolicon.app/setu/v2?r18=1&num=1")
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				picURL := gjson.Get(string(img), "data.0.urls.original").String()
				time.Sleep(time.Second * 3)
				deleteme := ctx.SendChain(message.At(user), message.Text("\n这是今日奖励哦"), message.Image(picURL))
				time.Sleep(time.Second * 20)
				ctx.DeleteMessage(deleteme)
			}
		})
}

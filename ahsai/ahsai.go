// Package ahsai AH Soft フリーテキスト音声合成 demo API
package ahsai

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/file"
	ahsaitts "github.com/fumiama/ahsai"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	namelist = []string{"伊織弓鶴", "紲星あかり", "結月ゆかり", "京町セイカ", "東北きりたん", "東北イタコ", "ついなちゃん標準語", "ついなちゃん関西弁", "音街ウナ", "琴葉茜", "吉田くん", "民安ともえ", "桜乃そら", "月読アイ", "琴葉葵", "東北ずん子", "月読ショウタ", "水奈瀬コウ"}
)

func init() {
	engine := control.Register("ahsai", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "ahsai tts\n- /ahsai -n 琴葉葵 -t にーはぉすーじぇ",
		PrivateDataFolder: "ahsai",
	})
	cachePath := engine.DataFolder() + "cache/"
	_ = os.RemoveAll(cachePath)
	_ = os.MkdirAll(cachePath, 0755)
	engine.OnCommand("ahsai", selectName).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("少女祈祷中..."))
		uid := ctx.Event.UserID
		now := time.Now()
		today := now.Format("20060102")
		ahsaiFile := cachePath + strconv.FormatInt(uid, 10) + today + "ahsai.wav"
		s := ahsaitts.NewSpeaker()
		err := s.SetName(ctx.State["ahsainame"].(string))
		if err != nil {
			ctx.SendChain(message.Text("Error:", err))
			return
		}
		u, err := s.Speak(ctx.State["ahsaitext"].(string))
		if err != nil {
			ctx.SendChain(message.Text("Error:", err))
			return
		}
		err = ahsaitts.SaveOggToFile(u, ahsaiFile)
		if err != nil {
			ctx.SendChain(message.Text("Error:", err))
			return
		}
		ctx.SendChain(message.Record("file:///" + file.BOTPATH + "/" + ahsaiFile))
	})
}

func selectName(ctx *zero.Ctx) bool {
	fset := flag.FlagSet{}
	var (
		name string
		text string
	)
	fset.StringVar(&name, "n", "", "speaker name")
	fset.StringVar(&text, "t", "にーはぉすーじぇ", "转换文本")
	ctx.State["ahsaitext"] = text
	arguments := shell.Parse(ctx.State["args"].(string))
	err := fset.Parse(arguments)
	if err != nil {
		ctx.SendChain(message.Text("Error:", err))
		return false
	}
	if name != "" {
		return true
	}
	speaktext := ""
	for i, v := range namelist {
		speaktext += fmt.Sprintf("%d. %s\n", i, v)
	}
	ctx.SendChain(message.Text("输入的音源为空, 请输入音源序号\n", speaktext))
	next, cancel := zero.NewFutureEvent("message", 999, false, ctx.CheckSession(), zero.RegexRule(`\d{0,2}`)).Repeat()
	defer cancel()
	for {
		select {
		case <-time.After(time.Second * 10):
			ctx.SendChain(message.Text("时间太久啦！", zero.BotConfig.NickName[0], "帮你选择"))
			ctx.State["ahsainame"] = namelist[rand.Intn(len(namelist))]
			return true
		case c := <-next:
			msg := c.Event.Message.ExtractPlainText()
			num, _ := strconv.Atoi(msg)
			if num < 0 || num >= len(namelist) {
				ctx.SendChain(message.Text("序号非法!"))
				continue
			}
			ctx.State["ahsainame"] = namelist[num]
			return true
		}
	}
}

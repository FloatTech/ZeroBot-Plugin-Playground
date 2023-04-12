// Package arcaea Arcaea类，不包含B30渲染
package arcaea

import (
	"encoding/json"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	aua "github.com/MoYoez/Arcaea_auaAPI"
	"github.com/fumiama/jieba/util/helper"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

import (
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"time"
)

type randsong struct {
	Status  int `json:"status"`
	Content struct {
		Id          string `json:"id"`
		RatingClass int    `json:"ratingClass"`
		Songinfo    struct {
			NameEn         string  `json:"name_en"`
			NameJp         string  `json:"name_jp"`
			Artist         string  `json:"artist"`
			Bpm            string  `json:"bpm"`
			BpmBase        float64 `json:"bpm_base"`
			Set            string  `json:"set"`
			SetFriendly    string  `json:"set_friendly"`
			Time           int     `json:"time"`
			Side           int     `json:"side"`
			WorldUnlock    bool    `json:"world_unlock"`
			RemoteDownload bool    `json:"remote_download"`
			Bg             string  `json:"bg"`
			Date           int     `json:"date"`
			Version        string  `json:"version"`
			Difficulty     int     `json:"difficulty"`
			Rating         int     `json:"rating"`
			Note           int     `json:"note"`
			ChartDesigner  string  `json:"chart_designer"`
			JacketDesigner string  `json:"jacket_designer"`
			JacketOverride bool    `json:"jacket_override"`
			AudioOverride  bool    `json:"audio_override"`
		} `json:"songinfo"`
	} `json:"content"`
}

type arcGPT struct {
	En     []string `json:"en"`
	Ja     []string `json:"ja"`
	ZhHans []string `json:"zh-Hans"`
}

var (
	randGPT        arcGPT
	arcRandSong    randsong
	getArcSongName string
	randLimit      = rate.NewManager[int64](time.Minute*5, 45) // 模仿 ArcGPT 的回复
)

func init() {
	engine := control.Register("arcaea", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "Arcaea相关，不包括B30",
		Help: "- 设置 AUA Key (必须)" +
			"- 设置 Server (必须)" +
			"- arcGPT",
		PublicDataFolder: "Arcaea",
	})
	// 使用这些类请确保你已经向 AUA (Awbugul 负责) 申请到您的 Key 和 API Sever 地址
	engine.OnRegex(`^设置\s?AUAkey\s*(.*)$`, zero.OnlyPrivate, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		file, _ := os.OpenFile(engine.DataFolder()+"apikey.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, err := file.WriteString(ctx.State["regex_matched"].([]string)[1])
		file.Close()
		if err != nil {
			ctx.SendChain(message.Text("设置apikey失败"))
			return
		}
		ctx.SendChain(message.Text("设置apikey成功"))
	})
	engine.OnRegex(`^设置\s?Server\s*(.*)$`, zero.OnlyPrivate, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		file, _ := os.OpenFile(engine.DataFolder()+"Server.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, err := file.WriteString(ctx.State["regex_matched"].([]string)[1])
		file.Close()
		if err != nil {
			ctx.SendChain(message.Text("设置Server失败"))
			return
		}
		ctx.SendChain(message.Text("设置Server成功"))
	})
	engine.OnFullMatch("arcGPT").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		list, err := engine.GetLazyData("list.json", false)
		if err != nil {
			ctx.SendChain(message.Text("ERR：", err))
			return
		}
		auaLink, err := os.ReadFile(engine.DataFolder() + "Server.txt")
		if err != nil {
			ctx.SendChain(message.Text("ERR：", err))
			return
		}
		auaKey, err := os.ReadFile(engine.DataFolder() + "apikey.txt")
		if err != nil {
			ctx.SendChain(message.Text("ERR：", err))
			return
		}
		switch {
		case randLimit.Load(ctx.Event.GroupID).AcquireN(6):
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Sending response to ArcGPT...Hold on("))
			auaRandSongBytes, err := aua.GetSongRandom(helper.BytesToString(auaLink), helper.BytesToString(auaKey), strconv.Itoa(0), strconv.Itoa(12))
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Cannot get message from arcgpt(", err))
				return
			}
			err = json.Unmarshal(list, &randGPT)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Cannot get message from arcgpt(", err))
				return
			}
			err = json.Unmarshal(auaRandSongBytes, &arcRandSong)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Cannot get message from arcgpt(", err))
				return
			}
			// get length.
			getGPTNums := rand.Intn(len(randGPT.ZhHans))
			// first check the jp name, if not, use eng name.
			getArcSongName = arcRandSong.Content.Songinfo.NameJp
			getSongArtist := arcRandSong.Content.Songinfo.Artist
			if getArcSongName == "" {
				getArcSongName = arcRandSong.Content.Songinfo.NameEn
			}
			// handle texts.
			handledSongGPTtextDone := strings.ReplaceAll(randGPT.ZhHans[getGPTNums], "歌曲名称", getArcSongName)
			handledSongGPTtext := strings.ReplaceAll(handledSongGPTtextDone, "作曲家", getSongArtist)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(handledSongGPTtext))
		case randLimit.Load(ctx.Event.GroupID).Acquire():
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("暂时无法处理更多请求。\n"))
		default:
		}
	})
}

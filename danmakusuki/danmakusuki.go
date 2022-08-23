// Package danmakusuki 弹幕网
package danmakusuki

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/img/writer"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/text"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type danmakusuki struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Data []struct {
			Channel struct {
				Name      string `json:"name"`
				IsLiving  bool   `json:"isLiving"`
				UID       int64  `json:"uId"`
				RoomID    int64  `json:"roomId"`
				FaceURL   string `json:"faceUrl"`
				LiveCount int64  `json:"liveCount"`
			} `json:"channel"`
			Live struct {
				LiveID        string  `json:"liveId"`
				Title         string  `json:"title"`
				IsFinish      bool    `json:"isFinish"`
				CoverURL      string  `json:"coverUrl"`
				StartDate     int64   `json:"startDate"`
				StopDate      int64   `json:"stopDate"`
				DanmakusCount int64   `json:"danmakusCount"`
				TotalIncome   float64 `json:"totalIncome"`
				WatchCount    int64   `json:"watchCount"`
			} `json:"live"`
			Danmakus []struct {
				Name     string  `json:"name"`
				Type     int64   `json:"type"`
				UID      int64   `json:"uId"`
				SendDate int64   `json:"sendDate"`
				Price    float64 `json:"price"`
				Message  string  `json:"message"`
			} `json:"danmakus"`
		} `json:"data"`
		Total    int64 `json:"total"`
		PageNum  int64 `json:"pageNum"`
		PageSize int64 `json:"pageSize"`
		HasMore  bool  `json:"hasMore"`
	} `json:"data"`
}

var (
	re             = regexp.MustCompile(`^\d+$`)
	danmakuTypeMap = map[int64]string{
		0: "普通消息",
		1: "礼物",
		2: "上舰",
		3: "Superchat",
		4: "进入直播间",
		5: "标题变动",
	}
)

const (
	searchUserURL = "http://api.bilibili.com/x/web-interface/search/type?search_type=bili_user&keyword=%v"
	danmakuAPI    = "https://danmaku.suki.club/api/search/user/detail?uid=%v&pagenum=%v&pagesize=5"
	danmakuURL    = "https://danmaku.suki.club/user/%v"
	memberCardURL = "https://account.bilibili.com/api/member/getCardByMid?mid=%v"
)

func init() {
	engine := control.Register("danmakusuki", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "https://danmaku.suki.club/\n查弹幕 嘉然 0 | 查弹幕 2 0 (第一个参数是name或uid, 第二个参数是页面)",
		PrivateDataFolder: "danmakusuki",
	})
	cachePath := engine.DataFolder() + "cache/"
	_ = os.RemoveAll(cachePath)
	_ = os.MkdirAll(cachePath, 0755)
	engine.OnRegex(`^查弹幕\s?(\S{1,25})\s?(\d*)$`, getPara).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		id := ctx.State["uid"].(string)
		pagenum := ctx.State["regex_matched"].([]string)[2]
		if pagenum == "" {
			pagenum = "0"
		}
		u, err := getMemberCard(id)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		var danmaku danmakusuki
		tr := &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}
		data, err := web.RequestDataWith(client, fmt.Sprintf(danmakuAPI, id, pagenum), "GET", "", web.RandUA())
		if err != nil {
			ctx.SendChain(message.Text("Error:", err))
			return
		}
		err = json.Unmarshal(data, &danmaku)
		if err != nil {
			ctx.SendChain(message.Text("Error:", err))
			return
		}
		today := time.Now().Format("20060102150415")
		drawedFile := cachePath + id + today + "vupLike.png"
		facePath := cachePath + id + "vupFace" + path.Ext(u.Face)
		backX := 500
		backY := 500
		var back image.Image
		if path.Ext(u.Face) != ".webp" {
			err = initFacePic(facePath, u.Face)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			back, err = gg.LoadImage(facePath)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			back = img.Size(back, backX, backY).Im
		}
		canvas := gg.NewContext(100, 100)
		fontSize := 50.0
		_, err = file.GetLazyData(text.BoldFontFile, true)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
		}
		if err = canvas.LoadFontFace(text.BoldFontFile, fontSize); err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		dz, h := canvas.MeasureString("好")
		danmuH := h * 2
		faceH := float64(510)

		totalDanmuku := 0
		for i := 0; i < len(danmaku.Data.Data); i++ {
			totalDanmuku += len(danmaku.Data.Data[i].Danmakus) + 1
		}
		cw := 10000
		mcw := float64(2000)
		ch := 550 + len(danmaku.Data.Data)*int(faceH) + totalDanmuku*int(danmuH)
		canvas = gg.NewContext(cw, ch)
		canvas.SetColor(color.White)
		canvas.Clear()
		canvas.SetColor(color.Black)
		if err = canvas.LoadFontFace(text.BoldFontFile, fontSize); err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		facestart := 100
		fontH := h * 1.6
		startWidth := float64(700)
		startWidth2 := float64(20)

		if back != nil {
			canvas.DrawImage(back, facestart, 0)
		}
		length, _ := canvas.MeasureString(u.Mid)
		n, _ := canvas.MeasureString(u.Name)
		canvas.DrawString(u.Name, startWidth, 122.5)
		canvas.DrawRoundedRectangle(900+n-length*0.1, 66, length*1.2, 75, fontSize*0.2)
		canvas.SetRGB255(221, 221, 221)
		canvas.Fill()
		canvas.SetColor(color.Black)
		canvas.DrawString(u.Mid, 900+n, 122.5)
		canvas.DrawString(fmt.Sprintf("粉丝：%d   关注：%d", u.Fans, u.Attention), startWidth, 222.5)
		canvas.DrawString(fmt.Sprintf("页码：[%d/%d]", danmaku.Data.PageNum, (danmaku.Data.Total-1)/5), startWidth, 322.5)
		canvas.DrawString("网页链接: "+fmt.Sprintf(danmakuURL, u.Mid), startWidth, 422.5)
		var channelStart float64
		channelStart = float64(550)
		for i := 0; i < len(danmaku.Data.Data); i++ {
			item := danmaku.Data.Data[i]
			facePath = cachePath + strconv.Itoa(int(item.Channel.UID)) + "vupFace" + path.Ext(item.Channel.FaceURL)
			if path.Ext(item.Channel.FaceURL) != ".webp" {
				err = initFacePic(facePath, item.Channel.FaceURL)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				back, err = gg.LoadImage(facePath)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				back = img.Size(back, backX, backY).Im
			}
			if back != nil {
				canvas.DrawImage(back, facestart, int(channelStart))
			}
			canvas.SetRGB255(24, 144, 255)
			canvas.DrawString("标题: "+item.Live.Title, startWidth, channelStart+fontH)
			canvas.DrawString("主播: "+item.Channel.Name, startWidth, channelStart+fontH*2)
			canvas.SetColor(color.Black)
			canvas.DrawString("开始时间: "+time.UnixMilli(item.Live.StartDate).Format("2006-01-02 15:04:05"), startWidth, channelStart+fontH*3)
			if item.Live.IsFinish {
				canvas.DrawString("结束时间: "+time.UnixMilli(item.Live.StopDate).Format("2006-01-02 15:04:05"), startWidth, channelStart+fontH*4)
				canvas.DrawString("直播时长: "+strconv.FormatFloat(float64(item.Live.StopDate-item.Live.StartDate)/3600000.0, 'f', 1, 64)+"小时", startWidth, channelStart+fontH*5)
			} else {
				t := "结束时间:"
				l, _ := canvas.MeasureString(t)
				canvas.DrawString(t, startWidth, channelStart+fontH*4)

				canvas.SetRGB255(0, 128, 0)
				t = "正在直播"
				canvas.DrawString(t, startWidth+l*1.1, channelStart+fontH*4)
				canvas.SetColor(color.Black)

				canvas.DrawString("直播时长: "+strconv.FormatFloat(float64(time.Now().UnixMilli()-item.Live.StartDate)/3600000.0, 'f', 1, 64)+"小时", startWidth, channelStart+fontH*5)
			}
			canvas.DrawString("弹幕数量: "+strconv.Itoa(int(item.Live.DanmakusCount)), startWidth, channelStart+fontH*6)
			canvas.DrawString("观看次数: "+strconv.Itoa(int(item.Live.WatchCount)), startWidth, channelStart+fontH*7)

			t := "收益:"
			l, _ := canvas.MeasureString(t)
			canvas.DrawString(t, startWidth, channelStart+fontH*8)

			t = "￥" + strconv.Itoa(int(item.Live.TotalIncome))
			canvas.SetRGB255(255, 0, 0)
			canvas.DrawString(t, startWidth+l*1.1, channelStart+fontH*8)
			canvas.SetColor(color.Black)

			DanmakuStart := channelStart + faceH
			for i := 0; i < len(item.Danmakus); i++ {
				moveW := startWidth2
				danmuNow := DanmakuStart + danmuH*float64(i+1)
				danItem := item.Danmakus[i]

				t := time.UnixMilli(danItem.SendDate).Format("15:04:05")
				l, _ := canvas.MeasureString(t)
				canvas.DrawString(t, moveW, danmuNow)
				moveW += l + dz

				t = danItem.Name
				l, _ = canvas.MeasureString(t)
				canvas.SetRGB255(24, 144, 255)
				canvas.DrawString(t, moveW, danmuNow)
				canvas.SetColor(color.Black)
				moveW += l + dz

				switch danItem.Type {
				case 0:
					t = danItem.Message
					l, _ = canvas.MeasureString(t)
					canvas.DrawString(t, moveW, danmuNow)
					moveW += l + dz
				case 1:
					t = danmakuTypeMap[danItem.Type]
					l, _ = canvas.MeasureString(t)
					canvas.SetRGB255(255, 0, 0)
					canvas.DrawString(t, moveW, danmuNow)
					moveW += l + dz

					t = danItem.Message
					l, _ = canvas.MeasureString(t)
					canvas.DrawString(t, moveW, danmuNow)
					canvas.SetColor(color.Black)
					moveW += l + dz
				case 2, 3:
					t = danmakuTypeMap[danItem.Type]
					l, _ = canvas.MeasureString(t)
					if danItem.Type == 3 {
						canvas.SetRGB255(0, 85, 255)
					} else {
						canvas.SetRGB255(128, 0, 128)
					}

					canvas.DrawString(t, moveW, danmuNow)
					moveW += l + dz

					t = danItem.Message
					l, _ = canvas.MeasureString(t)
					canvas.DrawString(t, moveW, danmuNow)
					moveW += l

					t = "["
					l, _ = canvas.MeasureString(t)
					canvas.DrawString(t, moveW, danmuNow)
					moveW += l

					t = "￥" + strconv.FormatFloat(danItem.Price, 'f', 1, 64)
					l, _ = canvas.MeasureString(t)
					canvas.SetRGB255(255, 0, 0)
					canvas.DrawString(t, moveW, danmuNow)
					if danItem.Type == 3 {
						canvas.SetRGB255(0, 85, 255)
					} else {
						canvas.SetRGB255(128, 0, 128)
					}
					moveW += l

					t = "]"
					l, _ = canvas.MeasureString(t)
					canvas.DrawString(t, moveW, danmuNow)
					canvas.SetColor(color.Black)
					moveW += l + dz
				case 4, 5:
					t = danmakuTypeMap[danItem.Type]
					canvas.SetRGB255(0, 128, 0)
					l, _ = canvas.MeasureString(t)
					canvas.DrawString(t, moveW, danmuNow)
					canvas.SetColor(color.Black)
					moveW += l + dz
				}
				if moveW > mcw {
					mcw = moveW
				}
			}
			channelStart = DanmakuStart + float64(len(item.Danmakus)+1)*danmuH
		}
		im := canvas.Image().(*image.RGBA)
		nim := im.SubImage(image.Rect(0, 0, int(mcw), ch))
		f, err := os.Create(drawedFile)
		if err != nil {
			log.Errorln("[danmakusuki]", err)
			data, cl := writer.ToBytes(nim)
			ctx.SendChain(message.ImageBytes(data))
			cl()
			return
		}
		_, err = writer.WriteTo(nim, f)
		_ = f.Close()
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
	})
}

func getPara(ctx *zero.Ctx) bool {
	keyword := ctx.State["regex_matched"].([]string)[1]
	if !re.MatchString(keyword) {
		searchRes, err := searchUser(keyword)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return false
		}
		ctx.State["uid"] = strconv.FormatInt(searchRes[0].Mid, 10)
		return true
	}
	next := zero.NewFutureEvent("message", 999, false, ctx.CheckSession())
	recv, cancel := next.Repeat()
	defer cancel()
	ctx.SendChain(message.Text("输入为纯数字，请选择查询uid还是用户名，输入对应序号：\n0. 查询uid\n1. 查询用户名"))
	for {
		select {
		case <-time.After(time.Second * 10):
			ctx.SendChain(message.Text("时间太久啦！", zero.BotConfig.NickName[0], "帮你选择查询uid"))
			ctx.State["uid"] = keyword
			return true
		case c := <-recv:
			msg := c.Event.Message.ExtractPlainText()
			num, err := strconv.Atoi(msg)
			if err != nil {
				ctx.SendChain(message.Text("请输入数字!"))
				continue
			}
			if num < 0 || num > 1 {
				ctx.SendChain(message.Text("序号非法!"))
				continue
			}
			if num == 0 {
				ctx.State["uid"] = keyword
				return true
			} else if num == 1 {
				searchRes, err := searchUser(keyword)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return false
				}
				ctx.State["uid"] = strconv.FormatInt(searchRes[0].Mid, 10)
				return true
			}
		}
	}
}

type searchResult struct {
	Mid    int64  `json:"mid"`
	Uname  string `json:"uname"`
	Gender int64  `json:"gender"`
	Usign  string `json:"usign"`
	Level  int64  `json:"level"`
}

// searchUser 查找b站用户
func searchUser(keyword string) (r []searchResult, err error) {
	data, err := web.GetData(fmt.Sprintf(searchUserURL, keyword))
	if err != nil {
		return
	}
	j := gjson.ParseBytes(data)
	if j.Get("data.numResults").Int() == 0 {
		err = errors.New("查无此人")
		return
	}
	err = json.Unmarshal(binary.StringToBytes(j.Get("data.result").Raw), &r)
	if err != nil {
		return
	}
	return
}

// memberCard 个人信息卡片
type memberCard struct {
	Mid        string  `json:"mid"`
	Name       string  `json:"name"`
	Sex        string  `json:"sex"`
	Face       string  `json:"face"`
	Coins      float64 `json:"coins"`
	Regtime    int64   `json:"regtime"`
	Birthday   string  `json:"birthday"`
	Sign       string  `json:"sign"`
	Attentions []int64 `json:"attentions"`
	Fans       int     `json:"fans"`
	Friend     int     `json:"friend"`
	Attention  int     `json:"attention"`
	LevelInfo  struct {
		CurrentLevel int `json:"current_level"`
	} `json:"level_info"`
}

// getMemberCard 获取b站个人详情
func getMemberCard(uid interface{}) (result memberCard, err error) {
	data, err := web.GetData(fmt.Sprintf(memberCardURL, uid))
	if err != nil {
		return
	}
	err = json.Unmarshal(binary.StringToBytes(gjson.ParseBytes(data).Get("card").Raw), &result)
	if err != nil {
		return
	}
	return
}

func initFacePic(filename, faceURL string) error {
	if file.IsNotExist(filename) {
		data, err := web.GetData(faceURL)
		if err != nil {
			return err
		}
		err = os.WriteFile(filename, data, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

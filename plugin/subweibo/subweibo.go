// Package subweibo 微博订阅
package subweibo

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type channelItem struct {
	ChannelKey  string `json:"channelKey"`
	ChannelName string `json:"channelName"`
	TestURL     string `json:"testApi"`
	ContURL     string `json:"contApi"`
}

var (
	testAPI         = "https://m.weibo.cn/api/container/getIndex?containerid=100505"
	contAPI         = "https://m.weibo.cn/api/container/getIndex?containerid=107603"
	channelItemData []*channelItem
	//  15天清除一次缓存清除订阅信息,  flush cache 后会产生会重复发message
	cacheMap       = cache.New(5*time.Minute, 360*time.Hour)
	messageSwitch  = false
	db             = &sql.Sqlite{}
	weiboMsgLocker sync.RWMutex
)

type weiboWebData struct {
	id                int64
	msgText           string
	msgPic            []gjson.Result
	scheme            string
	username          string
	createdAt         time.Time
	retweetedID       int64
	retweetedUserName string
	retweetedText     string
	retweetedPic      []gjson.Result
}

type weiboDBData struct {
	ID        int64     `db:"id"`
	Scheme    string    `db:"scheme"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
}

func (d *weiboDBData) throw(db *sql.Sqlite) error {
	weiboMsgLocker.Lock()
	defer weiboMsgLocker.Unlock()
	return db.Insert("weiboMsg", d)
}

func getWeiboMessageBox(url string) (weiboWebData, error) {
	var weiboWebData weiboWebData
	cont, err := getRequest(url)
	if err != nil {
		logrus.Error(cont)
		return weiboWebData, err
	}
	cards := gjson.Get(cont.String(), "data.cards").Array()
	for _, card := range cards {
		// 排除置顶微博和其它类型微博
		isTop := gjson.Get(card.String(), "mblog.mblogtype").Int()
		cardType := gjson.Get(card.String(), "card_type").Int()
		if isTop == 2 || cardType != 9 {
			continue
		}
		weiboWebData.id = gjson.Get(card.String(), "mblog.id").Int()
		weiboWebData.msgText = gjson.Get(card.String(), "mblog.text").String()
		weiboWebData.msgPic = gjson.Get(card.String(), "mblog.pics.#.large.url").Array()
		weiboWebData.scheme = gjson.Get(card.String(), "scheme").String()
		weiboWebData.username = gjson.Get(card.String(), "mblog.user.screen_name").String()
		weiboWebData.createdAt, _ = time.Parse(time.RubyDate, gjson.Get(card.String(), "mblog.created_at").String())
		isRetweeted := gjson.Get(card.String(), "mblog.retweeted_status").Exists()
		if isRetweeted {
			weiboWebData.retweetedID = gjson.Get(card.String(), "mblog.retweeted_status.id").Int()
			weiboWebData.retweetedUserName = gjson.Get(card.String(), "mblog.retweeted_status.user.screen_name").String()
			weiboWebData.retweetedText = gjson.Get(card.String(), "mblog.retweeted_status.text").String()
			weiboWebData.retweetedPic = gjson.Get(card.String(), "mblog.retweeted_status.pics.#.large.url").Array()
		}
		break
	}
	return weiboWebData, nil
}

func dataBuild(id int64, scheme string, username string, createdAt time.Time) *weiboDBData {
	return &weiboDBData{
		ID:        id,
		Scheme:    scheme,
		Username:  username,
		CreatedAt: createdAt,
	}
}

func getWeiboLink(url string) (str string, err error) {
	conn, err := getRequest(url)
	if err != nil {
		return "", err
	}
	value := gjson.Get(conn.String(), "data.userInfo.screen_name")
	return value.String(), nil
}
func getChannels(arg string) string {
	value, ok := cacheMap.Get(arg)
	if value == false || !ok {
		weiboName, err := getWeiboLink(testAPI + arg)
		if err != nil {
			return ""
		}
		if weiboName != "" {
			channelItemData = append(channelItemData, &channelItem{
				ChannelKey:  arg,
				ChannelName: weiboName,
				TestURL:     testAPI + arg,
				ContURL:     contAPI + arg,
			})
			cacheMap.Set(arg, true, cache.NoExpiration)
			return "已经成功订阅: " + weiboName + ",  UID: " + arg
		}
		return "未查询到订阅用户信息, 确认uid信息是否正确～"
	}
	return "请勿重复添加订阅:  " + arg
}
func delChannels(arg string) string {
	value, ok := cacheMap.Get(arg)
	if value == true && ok {
		for i, item := range channelItemData {
			if item.ChannelKey == arg {
				delName := item.ChannelName
				channelItemData = append(channelItemData[:i], channelItemData[i+1:]...)
				cacheMap.Set(arg, false, cache.NoExpiration)
				return "取消订阅: " + delName + ", UID:  " + arg + " 成功~"
			}
		}
	}
	return "还没有订阅: " + arg + ", 无法取消哦"
}
func running(ctx *zero.Ctx) {
	if messageSwitch {
		ctx.Send(message.Message{
			message.Text("已经开启了订阅，请勿重复开启"),
		})
	} else {
		ctx.Send(message.Message{
			message.Text("开启订阅成功, 准备开始接收消息"),
		})
		ticker := time.NewTicker(60 * time.Second)
		messageSwitch = true
		for range ticker.C {
			for _, item := range channelItemData {
				cURL := item.ContURL
				weiboMsgBoxData, err := getWeiboMessageBox(cURL)
				if err != nil || weiboMsgBoxData.msgText == "" {
					logrus.Error(err)
					continue
				}
				ok := db.CanFind("weiboMsg", "WHERE id="+strconv.FormatInt(weiboMsgBoxData.id, 10))
				if !ok {
					_ = dataBuild(weiboMsgBoxData.id, weiboMsgBoxData.scheme, weiboMsgBoxData.username, weiboMsgBoxData.createdAt).throw(db)
					if weiboMsgBoxData.retweetedText == "" {
						ctx.Send(message.Message{
							message.Text(weiboMsgBoxData.createdAt.String() + "\n" + weiboMsgBoxData.username + "发布了微博:\n" + trimHTML(weiboMsgBoxData.msgText) + "\n\nURL:" + weiboMsgBoxData.scheme),
						})
						for _, picURL := range weiboMsgBoxData.msgPic {
							picData, err := getRequest(picURL.String())
							if err != nil {
								logrus.Error("pic, Error: ", err)
								continue
							}
							ctx.Send(message.Message{
								message.ImageBytes(picData.Body()),
							})
						}
					} else {
						ok = db.CanFind("weiboMsg", "WHERE id="+strconv.FormatInt(weiboMsgBoxData.retweetedID, 10))
						if !ok {
							_ = dataBuild(weiboMsgBoxData.retweetedID, weiboMsgBoxData.scheme, weiboMsgBoxData.username, weiboMsgBoxData.createdAt).throw(db)
							ctx.Send(message.Message{
								message.Text(weiboMsgBoxData.createdAt.String() + "\n" + weiboMsgBoxData.username + "  转发了  " + weiboMsgBoxData.retweetedUserName + "  的微博:\n" + trimHTML(weiboMsgBoxData.retweetedText) + "\n评论了:\n" + trimHTML(weiboMsgBoxData.msgText) + "\nURL:" + weiboMsgBoxData.scheme),
							})
							for _, retURL := range weiboMsgBoxData.retweetedPic {
								retPicData, err := getRequest(retURL.String())
								if err != nil {
									logrus.Error("pic, Error: ", err)
									continue
								}
								ctx.Send(message.Message{
									message.ImageBytes(retPicData.Body()),
								})
							}
						}
					}
				}
			}
			if !messageSwitch {
				break
			}
		}
	}
}
func stop(ctx *zero.Ctx) {
	if messageSwitch {
		messageSwitch = false
		ctx.Send(message.Message{
			message.Text("关闭订阅成功！停止开始接收消息"),
		})
		// 清空订阅消息 清空缓存
		channelItemData = nil
		cacheMap.Flush()
	} else {
		ctx.Send(message.Message{
			message.Text("还未开启消息订阅哦～"),
		})
	}
}
func selectAllSubChannelsInfo(ctx *zero.Ctx) {
	var allChannelsInfo string
	for _, channel := range channelItemData {
		allChannelsInfo = allChannelsInfo + "\n" + channel.ChannelName + ",    UID:  " + channel.ChannelKey
	}
	if allChannelsInfo != "" {
		ctx.Send(message.Message{
			message.Text("当前已经订阅: \n", allChannelsInfo),
		})
	} else {
		ctx.Send(message.Message{
			message.Text("当前还为订阅任何内容哦～"),
		})
	}
}

func init() {
	engine := control.Register("weiboMessage", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "订阅微博消息",
		Help: "- 开启订阅 [UID] 开始接收订阅消息\n" +
			"- 关闭订阅 停止接收消息，并且清空缓存\n" +
			"- 订阅微博 [UID] 订阅xxx的微博消息\n" +
			"- 退订微博 [UID] 停止订阅xxx的微博消息\n" +
			"- 查看订阅 查看当前所有订阅",
		PrivateDataFolder: "subweibo",
	})
	go func() {
		db.DBPath = engine.DataFolder() + "weibo.db"
		err := db.Open(time.Hour * 24)
		if err != nil {
			panic(err)
		}
		err = db.Create("weiboMsg", &weiboDBData{})
		if err != nil {
			panic(err)
		}
	}()

	engine.OnFullMatch("开启订阅", zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		running(ctx)
	})
	engine.OnFullMatch("关闭订阅", zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		stop(ctx)
	})
	engine.OnFullMatch("查看订阅", zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		selectAllSubChannelsInfo(ctx)
	})
	engine.OnPrefix("订阅微博", zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		args := ctx.State["args"].(string)
		msg := getChannels(args)
		ctx.SendChain(message.Text(msg))
	})
	engine.OnPrefix("退订微博", zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		arg := ctx.State["args"].(string)
		msg := delChannels(arg)
		ctx.SendChain(message.Text(msg))
	})
	initChannel()
}

func initChannel() {
	var args []string
	// "7791102134",  "3273865405",  "5533669771",  "2732469654",  "3223557554",  "2339808364"
	args = append(args, "")
	for _, arg := range args {
		msg := getChannels(arg)
		logrus.Info("初始化订阅信息", msg)
	}
}

func trimHTML(src string) string {
	// 将HTML标签全转换成小写
	re := regexp.MustCompile(`<[\S\s]+?>`)
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	// 去除STYLE
	re = regexp.MustCompile(`<style[\S\s]+?</style>`)
	src = re.ReplaceAllString(src, "")
	// 去除SCRIPT
	re = regexp.MustCompile(`<script[\S\s]+?</script>`)
	src = re.ReplaceAllString(src, "")
	// 去除所有尖括号内的HTML代码，并换成换行符
	re = regexp.MustCompile(`<[\S\s]+?>`)
	src = re.ReplaceAllString(src, "\n")
	// 去除连续的换行符
	re = regexp.MustCompile(`\s{2,}`)
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}
func getRequest(url string) (resp *resty.Response, err error) {
	client := resty.New()
	resp, err = client.R().Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

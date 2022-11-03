package subweibo

import (
	"fmt"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type channelItem struct {
	ChannelKey  string `json:"channelKey"`
	ChannelName string `json:"channelName"`
	TestUrl     string `json:"testApi"`
	ContUri     string `json:"contApi"`
}
type wbFunc struct{}

var (
	testApi         = "https://m.weibo.cn/api/container/getIndex?containerid=100505"
	contApi         = "https://m.weibo.cn/api/container/getIndex?containerid=107603"
	channelItemData []*channelItem
	//  15天清除一次缓存, flush cache 后会产生会重复发message
	cacheMap      = cache.New(5*time.Minute, 360*time.Hour)
	wb            = wbFunc{}
	messageSwitch = false
)

func (w *wbFunc) TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("<[\\S\\s]+?>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("<style[\\S\\s]+?</style>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("<script[\\S\\s]+?</script>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("<[\\S\\s]+?>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

func (w *wbFunc) getRequest(url string) string {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	result, _ := io.ReadAll(resp.Body)
	return string(result)
}

func (w *wbFunc) getWeiboMessageBox(url string) (string, string, []gjson.Result, string, string, string) {

	cont := w.getRequest(url)
	cards := gjson.Get(cont, "data.cards").Array()
	for _, card := range cards {
		// 排除置顶微博
		isTop := gjson.Get(card.String(), "mblog.title").String()
		if isTop != "" {
			continue
		} else {
			profileId := gjson.Get(card.String(), "profile_type_id").String()
			msgText := gjson.Get(card.String(), "mblog.text").String()
			msgPic := gjson.Get(card.String(), "mblog.pics.#.large.url").Array()
			scheme := gjson.Get(card.String(), "scheme").String()
			username := gjson.Get(card.String(), "mblog.user.screen_name").String()
			createdAt, _ := time.Parse(time.RubyDate, gjson.Get(card.String(), "mblog.created_at").String())
			return profileId, msgText, msgPic, scheme, username, createdAt.String()
		}
	}
	return "", "", nil, "", "", ""
}

func (w *wbFunc) getImageByUrl(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
	}
	body, _ := io.ReadAll(resp.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	return body
}

func (w *wbFunc) getWeiboLink(url string) string {
	conn := w.getRequest(url)
	value := gjson.Get(conn, "data.userInfo.screen_name")
	return value.String()
}
func (w *wbFunc) getChannels(arg string, ctx *zero.Ctx) {
	value, ok := cacheMap.Get(arg)
	if value == false || !ok {
		weiboName := wb.getWeiboLink(testApi + arg)
		if weiboName != "" {
			ctx.Send(message.Message{
				message.Text("已经成功订阅:  " + weiboName + ",uid:(" + arg + ")"),
			})
			channelItemData = append(channelItemData, &channelItem{
				ChannelKey:  arg,
				ChannelName: weiboName,
				TestUrl:     testApi + arg,
				ContUri:     contApi + arg,
			})
			cacheMap.Set(arg, true, cache.NoExpiration)
		}
	} else {
		ctx.Send(message.Message{
			message.Text("请勿重复添加订阅:  ", arg),
		})
	}
	return
}
func (w *wbFunc) delChannels(arg string, ctx *zero.Ctx) {
	value, ok := cacheMap.Get(arg)
	if value == true && ok {
		for i, item := range channelItemData {
			if item.ChannelKey == arg {
				delName := item.ChannelName
				channelItemData = append(channelItemData[:i], channelItemData[i+1:]...)
				ctx.Send(message.Message{
					message.Text("取消订阅: " + delName + ",UID:  " + arg + " 成功~"),
				})
				cacheMap.Set(arg, false, cache.NoExpiration)
			}
		}
	} else {
		ctx.Send(message.Message{
			message.Text("还没有订阅：", arg, ",无法取消哦"),
		})
	}
}
func (w *wbFunc) running(ctx *zero.Ctx) {
	if messageSwitch == true {
		ctx.Send(message.Message{
			message.Text("已经开启了订阅，请勿重复开启"),
		})
	} else {
		ctx.Send(message.Message{
			message.Text("开启订阅成功,准备开始接收消息"),
		})
		ticker := time.NewTicker(60 * time.Second)
		messageSwitch = true
		for {
			select {
			case <-ticker.C:
				for _, item := range channelItemData {
					cUrl := item.ContUri
					pId, mText, mPic, scheme, username, creatAt := w.getWeiboMessageBox(cUrl)
					_, ok := cacheMap.Get(pId)
					if ok == false {
						cacheMap.Set(pId, true, cache.NoExpiration)
						ctx.Send(message.Message{
							message.Text(creatAt + "\n" + username + "发布了微博:\n" + w.TrimHtml(mText) + "\n\nURL:" + scheme),
						})
						for _, picUrl := range mPic {
							ctx.Send(message.Message{
								message.ImageBytes(w.getImageByUrl(picUrl.String())),
							})
						}
					}
				}
			}
			if messageSwitch == false {
				break
			}
		}
	}

}
func (w *wbFunc) stop(ctx *zero.Ctx) {
	if messageSwitch == true {
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

func (w *wbFunc) selectAllSubChannelsInfo(ctx *zero.Ctx) {
	var allChannelsInfo string
	for _, channel := range channelItemData {
		fmt.Println(channel.ChannelKey)
		fmt.Println(channel.ChannelName)
		allChannelsInfo = allChannelsInfo + "\n" + channel.ChannelName + ",   UID:  " + channel.ChannelKey
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
		Help: "订阅微博消息\n" +
			"- 开启订阅 [UID] 开始接收订阅消息\n" +
			"- 关闭订阅 停止接收消息，并且清空缓存\n" +
			"- 订阅微博 [UID] 订阅xxx的微博消息\n" +
			"- 退订微博 [UID] 停止订阅xxx的微博消息\n" +
			"- 查看订阅 查看当前所有订阅",
	})
	engine.OnFullMatch("开启订阅", zero.AdminPermission, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		go wb.running(ctx)
	})
	engine.OnFullMatch("关闭订阅", zero.AdminPermission, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		wb.stop(ctx)
	})
	engine.OnFullMatch("查看订阅", zero.AdminPermission, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		wb.selectAllSubChannelsInfo(ctx)
	})
	engine.OnPrefix("订阅微博", zero.AdminPermission, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		args := ctx.State["args"].(string)
		wb.getChannels(args, ctx)
	})
	engine.OnPrefix("退订微博", zero.AdminPermission, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		arg := ctx.State["args"].(string)
		wb.delChannels(arg, ctx)
	})
}

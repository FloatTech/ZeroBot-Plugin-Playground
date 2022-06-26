// Package dynamic b站动态解析库
package dynamic

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	referer = "https://www.bilibili.com/"
	// TURL bilibili动态前缀
	TURL = "https://t.bilibili.com/"
	// SpaceHistoryURL 历史动态信息,一共12个card
	SpaceHistoryURL = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/space_history?host_uid=%v&offset_dynamic_id=%v&need_top=0"
	// DetailURL 当前动态信息,一个card
	DetailURL = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail?dynamic_id=%v"
)

var (
	typeMsg = map[int]string{
		1:    "转发了动态",
		2:    "有图营业",
		4:    "无图营业",
		8:    "投稿了视频",
		16:   "投稿了短视频",
		64:   "投稿了文章",
		256:  "投稿了音频",
		2048: "发布了简报",
		4200: "发布了直播",
		4308: "发布了直播",
	}
)

// DynCard 总动态结构体,包括desc,card
type DynCard struct {
	Desc Desc   `json:"desc"`
	Card string `json:"card"`
}

// Desc 描述结构体
type Desc struct {
	Type         int    `json:"type"`
	DynamicIDStr string `json:"dynamic_id_str"`
	OrigType     int    `json:"orig_type"`
	Timestamp    int    `json:"timestamp"`
	Origin       struct {
		DynamicIDStr string `json:"dynamic_id_str"`
	} `json:"origin"`
	UserProfile struct {
		Info struct {
			Uname string `json:"uname"`
		} `json:"info"`
	} `json:"user_profile"`
}

// Card 卡片结构体
type Card struct {
	Item struct {
		Content     string `json:"content"`
		UploadTime  int    `json:"upload_time"`
		Description string `json:"description"`
		Pictures    []struct {
			ImgSrc string `json:"img_src"`
		} `json:"pictures"`
		Timestamp int `json:"timestamp"`
		Cover     struct {
			Default string `json:"default"`
		} `json:"cover"`
		OrigType int `json:"orig_type"`
	} `json:"item"`
	AID       interface{} `json:"aid"`
	Bvid      interface{} `json:"bvid"`
	Dynamic   interface{} `json:"dynamic"`
	Pic       string      `json:"pic"`
	Title     string      `json:"title"`
	ID        int         `json:"id"`
	Summary   string      `json:"summary"`
	ImageUrls []string    `json:"image_urls"`
	Sketch    struct {
		Title     string `json:"title"`
		DescText  string `json:"desc_text"`
		CoverURL  string `json:"cover_url"`
		TargetURL string `json:"target_url"`
	} `json:"sketch"`
	Stat struct {
		Aid      int `json:"aid"`
		View     int `json:"view"`
		Danmaku  int `json:"danmaku"`
		Reply    int `json:"reply"`
		Favorite int `json:"favorite"`
		Coin     int `json:"coin"`
		Share    int `json:"share"`
		Like     int `json:"like"`
	} `json:"stat"`
	Owner struct {
		Name    string `json:"name"`
		Pubdate int    `json:"pubdate"`
	} `json:"owner"`
	Cover        string      `json:"cover"`
	ShortID      interface{} `json:"short_id"`
	LivePlayInfo struct {
		ParentAreaName string `json:"parent_area_name"`
		AreaName       string `json:"area_name"`
		Cover          string `json:"cover"`
		Link           string `json:"link"`
		Online         int    `json:"online"`
		RoomID         int    `json:"room_id"`
		LiveStatus     int    `json:"live_status"`
		WatchedShow    string `json:"watched_show"`
		Title          string `json:"title"`
	} `json:"live_play_info"`
	Intro    string      `json:"intro"`
	Schema   string      `json:"schema"`
	Author   interface{} `json:"author"`
	PlayCnt  int         `json:"play_cnt"`
	ReplyCnt int         `json:"reply_cnt"`
	TypeInfo string      `json:"type_info"`
	User     struct {
		Name  string `json:"name"`
		Uname string `json:"uname"`
	} `json:"user"`
	Desc          string `json:"desc"`
	ShareSubtitle string `json:"share_subtitle"`
	ShortLink     string `json:"short_link"`
	PublishTime   int    `json:"publish_time"`
	BannerURL     string `json:"banner_url"`
	Ctime         int    `json:"ctime"`
	Vest          struct {
		Content string `json:"content"`
	} `json:"vest"`
	Upper   string `json:"upper"`
	Origin  string `json:"origin"`
	Pubdate int    `json:"pubdate"`
}

// Card2msg cType=0时,处理DynCard字符串,cType=1, 2, 4, 8, 16, 64, 256, 2048, 4200, 4308时,处理Card字符串,cType为card类型
func Card2msg(str string, cType int) (msg []message.MessageSegment, err error) {
	var (
		DynCard DynCard
		card    Card
	)
	// 初始化结构体
	switch cType {
	case 0:
		err = json.Unmarshal([]byte(str), &DynCard)
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(DynCard.Card), &card)
		if err != nil {
			return
		}
		cType = DynCard.Desc.Type
	case 1, 2, 4, 8, 16, 64, 256, 2048, 4200, 4308:
		err = json.Unmarshal([]byte(str), &card)
		if err != nil {
			return
		}
	default:
		err = errors.New("只有0, 1, 2, 4, 8, 16, 64, 256, 2048, 4200, 4308模式")
		return
	}
	// 生成消息
	switch cType {
	case 1:
		msg = append(msg, message.Text(card.User.Uname, typeMsg[cType], "\n"))
		msg = append(msg, message.Text(card.Item.Content, "\n"))
		msg = append(msg, message.Text("转发的内容：\n"))
		var originMsg []message.MessageSegment
		originMsg, err = Card2msg(card.Origin, card.Item.OrigType)
		if err != nil {
			return
		}
		msg = append(msg, originMsg...)
	case 2:
		msg = append(msg, message.Text(card.User.Name, "在", time.Unix(int64(card.Item.UploadTime), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n"))
		msg = append(msg, message.Text(card.Item.Description))
		for i := 0; i < len(card.Item.Pictures); i++ {
			msg = append(msg, message.Image(card.Item.Pictures[i].ImgSrc))
		}
	case 4:
		msg = append(msg, message.Text(card.User.Name, "在", time.Unix(int64(card.Item.Timestamp), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n"))
		msg = append(msg, message.Text(card.Item.Content, "\n"))
	case 8:
		msg = append(msg, message.Text(card.Owner.Name, "在", time.Unix(int64(card.Pubdate), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n"))
		msg = append(msg, message.Text(card.Title))
		msg = append(msg, message.Image(card.Pic))
		msg = append(msg, message.Text(card.Desc, "\n"))
		msg = append(msg, message.Text(card.ShareSubtitle, "\n"))
		msg = append(msg, message.Text("视频链接：", card.ShortLink, "\n"))
	case 16:
		msg = append(msg, message.Text(card.User.Name, "在", time.Unix(int64(card.Item.UploadTime), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n"))
		msg = append(msg, message.Text(card.Item.Description))
		msg = append(msg, message.Image(card.Item.Cover.Default))
	case 64:
		msg = append(msg, message.Text(card.Author.(map[string]interface{})["name"], "在", time.Unix(int64(card.PublishTime), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n"))
		msg = append(msg, message.Text(card.Title, "\n"))
		msg = append(msg, message.Text(card.Summary))
		for i := 0; i < len(card.ImageUrls); i++ {
			msg = append(msg, message.Image(card.ImageUrls[i]))
		}
		if card.ID != 0 {
			msg = append(msg, message.Text("文章链接：https://www.bilibili.com/read/cv", card.ID, "\n"))
		}
	case 256:
		msg = append(msg, message.Text(card.Upper, "在", time.Unix(int64(card.Ctime), 0).Format("2006-01-02 15:04:05"), typeMsg[cType], "\n"))
		msg = append(msg, message.Text(card.Title))
		msg = append(msg, message.Image(card.Cover))
		msg = append(msg, message.Text(card.Intro, "\n"))
		if card.ID != 0 {
			msg = append(msg, message.Text("音频链接：https://www.bilibili.com/audio/au", card.ID, "\n"))
		}

	case 2048:
		msg = append(msg, message.Text(card.User.Uname, typeMsg[cType], "\n"))
		msg = append(msg, message.Text(card.Vest.Content, "\n"))
		msg = append(msg, message.Text(card.Sketch.Title, "\n"))
		msg = append(msg, message.Text(card.Sketch.DescText, "\n"))
		msg = append(msg, message.Image(card.Sketch.CoverURL))
		msg = append(msg, message.Text("分享链接：", card.Sketch.TargetURL, "\n"))
	case 4308:
		if DynCard.Desc.UserProfile.Info.Uname != "" {
			msg = append(msg, message.Text(DynCard.Desc.UserProfile.Info.Uname, typeMsg[cType], "\n"))
		}
		msg = append(msg, message.Image(card.LivePlayInfo.Cover))
		msg = append(msg, message.Text(card.LivePlayInfo.Title, "\n"))
		msg = append(msg, message.Text("房间号：", card.LivePlayInfo.RoomID, "\n"))
		msg = append(msg, message.Text("分区：", card.LivePlayInfo.ParentAreaName))
		if card.LivePlayInfo.ParentAreaName != card.LivePlayInfo.AreaName {
			msg = append(msg, message.Text("-", card.LivePlayInfo.AreaName))
		}
		if card.LivePlayInfo.LiveStatus == 0 {
			msg = append(msg, message.Text("未开播 \n"))
		} else {
			msg = append(msg, message.Text("直播中 ", card.LivePlayInfo.WatchedShow, "\n"))
		}
		msg = append(msg, message.Text("直播链接：", card.LivePlayInfo.Link))
	default:
		msg = append(msg, message.Text("动态id：", DynCard.Desc.DynamicIDStr, "未知动态类型：", cType, "\n"))
	}
	if DynCard.Desc.DynamicIDStr != "" {
		msg = append(msg, message.Text("动态链接：", TURL, DynCard.Desc.DynamicIDStr))
	}
	return
}

// Detail 通过动态id生成消息
func Detail(dynamicIDStr string) (msg []message.MessageSegment, err error) {
	var data []byte
	data, err = web.RequestDataWith(web.NewDefaultClient(), fmt.Sprintf(DetailURL, dynamicIDStr), "GET", referer, web.RandUA())
	if err != nil {
		return
	}
	return Card2msg(gjson.ParseBytes(data).Get("data.card").Raw, 0)
}

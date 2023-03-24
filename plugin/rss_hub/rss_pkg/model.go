package rss_pkg

import (
	"encoding/hex"
	"hash/fnv"
	"time"
)

// ======== RSS ========[START]

//type SingleFeedItem gofeed.Item

func genHashForFeedItem(link, guid string) string {
	idString := link + "||" + guid
	h := fnv.New32()
	_, _ = h.Write([]byte(idString))
	encoded := hex.EncodeToString(h.Sum(nil))
	return encoded
}

// RssChannelView 频道视图
type RssChannelView struct {
	Channel  *RssFeedChannel
	Contents []*RssContent
}

// ======== RSS ========[END]

// ======== DB ========[START]

const (
	tableNameRssFeedChannel = "rss_feed_channel"
	tableNameRssFeedContent = "rss_feed_content"
	tableNameRssSubscribe   = "rss_subscribe"
)

// RssFeedChannel 订阅的RSS频道
type RssFeedChannel struct {
	// Id 自增id
	Id int64 `gorm:"primary_key;AUTO_INCREMENT"`
	// RssHubFeedPath 频道路由 用于区分rss_hub 不同的频道 例如: `/bangumi/tv/calendar/today`
	RssHubFeedPath string `gorm:"column:rss_hub_feed_path;unique;not null" json:"rss_hub_feed_path"`
	// Title 频道标题
	Title string `gorm:"column:title"        json:"title"`
	// ChannelDesc 频道描述
	ChannelDesc string `gorm:"column:channel_desc" json:"channel_desc"`
	// ImageUrl 频道图片
	ImageUrl string `gorm:"column:image_url"    json:"image_url"`
	// Link 频道链接
	Link string `gorm:"column:link"         json:"link"`
	// UpdatedParsed RSS页面更新时间
	UpdatedParsed time.Time `gorm:"column:updated_parsed" json:"updated_parsed"`
	//// Ctime create time
	//Ctime int64 `gorm:"column:ctime;default:current_timestamp"  json:"ctime"`
	// Mtime update time
	Mtime time.Time `gorm:"column:mtime;default:current_timestamp;" json:"mtime"`
}

// TableName ...
func (RssFeedChannel) TableName() string {
	return tableNameRssFeedChannel
}

// IfNeedUpdate ...
func (r RssFeedChannel) IfNeedUpdate(cmp *RssFeedChannel) bool {
	if r.Link != cmp.Link {
		return false
	}
	return r.UpdatedParsed.Unix() < cmp.UpdatedParsed.Unix()
}

// RssContent 订阅的RSS频道的推送信息
type RssContent struct {
	// Id 自增id
	Id               int64     `gorm:"primary_key;AUTO_INCREMENT"`
	HashId           string    `gorm:"column:hash_id;unique"        json:"hash_id"`
	RssFeedChannelId int64     `gorm:"column:rss_feed_channel_id;not null"   json:"rss_feed_channel_id"`
	Title            string    `gorm:"column:title"       json:"title"`
	Description      string    `gorm:"column:description" json:"description"`
	Link             string    `gorm:"column:link"        json:"link"`
	Date             time.Time `gorm:"column:date"        json:"date"`
	Author           string    `gorm:"column:author"      json:"author"`
	Thumbnail        string    `gorm:"column:thumbnail"   json:"thumbnail"`
	Content          string    `gorm:"column:content"     json:"content"`
	//// Ctime create time
	//Ctime int64 `gorm:"column:ctime;default:current_timestamp"  json:"ctime"`
	// Mtime update time
	Mtime time.Time `gorm:"column:mtime;default:current_timestamp;" json:"mtime"`
}

// TableName ...
func (RssContent) TableName() string {
	return tableNameRssFeedContent
}

// RssSubscribe 订阅关系表：群组-RSS频道
type RssSubscribe struct {
	// Id 自增id
	Id int64 `gorm:"primary_key;AUTO_INCREMENT"`
	// 订阅群组
	GroupId int64 `gorm:"column:group_id;not null"`
	// 订阅频道
	RssFeedChannelId int64 `gorm:"column:rss_feed_channel_id;not null"`
	//// Ctime create time
	//Ctime int64 `gorm:"column:ctime;default:current_timestamp"  json:"ctime"`
	// Mtime update time
	Mtime time.Time `gorm:"column:mtime;default:current_timestamp;" json:"mtime"`
}

// TableName ...
func (RssSubscribe) TableName() string {
	return tableNameRssSubscribe
}

// ======== DB ========[END]

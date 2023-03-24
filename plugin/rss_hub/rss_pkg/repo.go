package rss_pkg

import "context"

// RepoContent RSS 推送信息存储接口
type RepoContent interface {
	// UpsertContent 添加一条文章
	UpsertContent(ctx context.Context, content *RssContent) error
	// DeleteSourceContents 删除订阅源的所有文章，返回被删除的文章数
	DeleteSourceContents(ctx context.Context, channelID int64) (int64, error)
	// HashIDExist hash id 对应的文章是否已存在
	IsContentHashIDExist(ctx context.Context, hashID string) (bool, error)
}

// RepoSource RSS 订阅源存储接口
type RepoSource interface {
	// UpsertSource 添加一个订阅源
	UpsertSource(ctx context.Context, rfc *RssFeedChannel) error
	// GetSource 获取一个订阅源信息
	GetSource(ctx context.Context, id int64) (*RssFeedChannel, error)
	// GetSources 获取所有订阅源信息
	GetSources(ctx context.Context) ([]*RssFeedChannel, error)
	// GetSourceByRssHubFeedLink 通过 rssHub 的 feed 链接获取订阅源信息
	GetSourceByRssHubFeedLink(ctx context.Context, url string) (*RssFeedChannel, error)
	// DeleteSource 删除一个订阅源
	DeleteSource(ctx context.Context, id int64) error
	//UpsertSource(ctx context.Context, sourceID uint, newSource *RssFeedChannelRoute) error
}

type RepoSubscribe interface {
	// CreateSubscribe 添加一个订阅
	CreateSubscribe(ctx context.Context, gid, subscribe int64) error
	// DeleteSubscribe 删除一个订阅
	DeleteSubscribe(ctx context.Context, gid int64, subscribeId int64) error
	// GetSubscribeById 获取一个订阅
	GetSubscribeById(ctx context.Context, gid int64, subscribeId int64) (*RssSubscribe, error)
	// GetSubscribes 获取全部订阅
	GetSubscribes(ctx context.Context) ([]*RssSubscribe, error)
	// GetSubscribedChannelsByGroupId 获取该群所有的订阅
	//GetSubscribedChannelsByGroupId(ctx context.Context, gid int64) ([]*RssSubscribe, error)
}

type RepoMultiQuery interface {
	// GetSubscribesBySource 获取一个源对应的所有订阅群组
	GetSubscribesBySource(ctx context.Context, feedPath string) ([]*RssSubscribe, error)
	// GetIfExistedSubscribe 判断一个群组是否已订阅了一个源
	GetIfExistedSubscribe(ctx context.Context, gid int64, feedPath string) (*RssSubscribe, bool, error)
	// GetSubscribedChannelsByGroupId 获取该群所有的订阅
	GetSubscribedChannelsByGroupId(ctx context.Context, gid int64) ([]*RssFeedChannel, error)
}

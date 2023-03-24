package rss_pkg

import (
	"context"
	"errors"
	sql "github.com/FloatTech/sqlite"
	"github.com/sirupsen/logrus"
	"net/http"
)

type RssDomain interface {
	// Subscribe 订阅Rss频道
	Subscribe(ctx context.Context, gid int64, route string) (rv *RssChannelView, isChannelExisted,
		isSubExisted bool, err error)
	// Unsubscribe 取消订阅Rss频道
	Unsubscribe(ctx context.Context, gid int64, route string) (err error)
	// GetSubscribedChannelsByGroupId 获取群组订阅的Rss频道
	GetSubscribedChannelsByGroupId(ctx context.Context, gid int64) (rv []*RssChannelView, err error)
	//GetRssFeedChannel(ctx context.Context, id uint) (rv *RssChannelView, err error)
	//GetRssFeedChannelList(ctx context.Context, gid int64) (rv []*RssChannelView, err error)
	//DeleteRssFeedChannel(ctx context.Context, id uint) (err error)
	//SubscribeRssFeedChannel(ctx context.Context, gid int64, id uint) (err error)
	//UnsubscribeRssFeedChannel(ctx context.Context, gid int64, id uint) (err error)
	//GetRssFeedContentList(ctx context.Context, gid int64, id uint) (rv []*RssContentView, err error)
	//GetRssFeedContent(ctx context.Context, id uint) (rv *RssContentView, err error)

	// SyncJobTrigger 触发同步任务
	// SyncJobTrigger(ctx context.Context) (err error)'

	// SyncJobTrigger 同步Rss频道
	SyncJobTrigger(ctx context.Context) (groupView map[int64][]*RssChannelView, err error)

	// SyncJobTrigger 同步Rss频道

	//SyncJobTrigger(ctx context.Context) (groupView map[int64]*RssChannelView, err error)

	//// StopSyncJob 停止同步任务
	//StopSyncJob()
	//// StartSyncJob 启动同步任务
	//StartSyncJob()
}

// rssDomain RssRepo定义
type rssDomain struct {
	storage      RepoStorage
	rssHubClient *RssHubClient
}

// NewRssDomain 新建RssDomain，调用方保证单例模式
func NewRssDomain(dbPath string) (RssDomain, error) {
	return newRssDomain(dbPath)
}

func newRssDomain(dbPath string) (*rssDomain, error) {
	repo := &rssDomain{
		storage: &repoStorage{
			db: sql.Sqlite{
				DBPath: dbPath + "rss_hub.db",
			},
		},
		rssHubClient: &RssHubClient{Client: http.DefaultClient},
	}
	err := repo.storage.initDB()
	if err != nil {
		logrus.Errorf("[rss_hub NewRssDomain] open db error: %v", err)
		return nil, err
	}
	// 启动同步任务
	//repo.rssCronTab.cronTab.Start()
	return repo, nil
}

// Subscribe QQ群订阅Rss频道
func (repo *rssDomain) Subscribe(ctx context.Context, gid int64, feedPath string) (
	rv *RssChannelView, isChannelExisted, isSubExisted bool, err error) {
	// 验证
	feed, err := repo.rssHubClient.FetchFeed(rssHubMirrors[0], feedPath)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub Subscribe] add source error: %v", err)
		return
	}
	logrus.WithContext(ctx).Infof("[rss_hub Subscribe] try get source success: %v", len(feed.Title))
	// 新建source结构体
	rv = convertFeedToRssChannelView(0, feedPath, feed)
	rfExisted, err := repo.storage.GetSourceByRssHubFeedLink(ctx, feedPath)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub Subscribe] query source by feedPath error: %v", err)
		return
	}
	// 如果已经存在
	if rfExisted != nil {
		logrus.WithContext(ctx).Infof("[rss_hub Subscribe] source existed: %v", rfExisted)
		isChannelExisted = true
	}
	// 保存
	err = repo.storage.UpsertSource(ctx, rv.Channel)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub Subscribe] save source error: %v", err)
		return
	}
	logrus.Infof("[rss_hub Subscribe] save/update source success %v", rv.Channel.Id)
	// 添加群号到订阅
	subscribe, err := repo.storage.GetSubscribeById(ctx, gid, rv.Channel.Id)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub Subscribe] query subscribe error: %v", err)
		return
	}
	logrus.WithContext(ctx).Infof("[rss_hub Subscribe] query subscribe success: %v", subscribe)
	// 如果已经存在，直接返回
	if subscribe != nil {
		isSubExisted = true
		logrus.WithContext(ctx).Infof("[rss_hub Subscribe] subscribe existed: %v", subscribe)
		return
	}
	// 如果不存在，保存
	err = repo.storage.CreateSubscribe(ctx, gid, rv.Channel.Id)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub Subscribe] save subscribe error: %v", err)
		return
	}
	logrus.WithContext(ctx).Infof("[rss_hub Subscribe] success: %v", len(rv.Contents))
	return
}

// Unsubscribe 群组取消订阅
func (repo *rssDomain) Unsubscribe(ctx context.Context, gid int64, feedPath string) (err error) {
	rf, ifExisted, err := repo.storage.GetIfExistedSubscribe(ctx, gid, feedPath)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub Subscribe] query sub by route error: %v", err)
		return errors.New("数据库错误")
	}
	logrus.WithContext(ctx).Infof("[rss_hub Subscribe] query source by route success: %v", rf)
	// 如果不存在订阅关系，直接返回
	if !ifExisted || rf == nil {
		logrus.WithContext(ctx).Infof("[rss_hub Subscribe] source existed: %v", ifExisted)
		return errors.New("频道不存在")
	}
	err = repo.storage.DeleteSubscribe(ctx, gid, rf.Id)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub Subscribe] delete source error: %v", err)
		return errors.New("删除失败")
	}
	return
}

// GetSubscribedChannelsByGroupId 获取群对应的订阅的频道信息
func (repo *rssDomain) GetSubscribedChannelsByGroupId(ctx context.Context, gid int64) (rv []*RssChannelView, err error) {
	rv = make([]*RssChannelView, 0)
	channels, err := repo.storage.GetSubscribedChannelsByGroupId(ctx, gid)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub GetSubscribedChannelsByGroupId] GetSubscribedChannelsByGroupId error: %v", err)
		return
	}
	logrus.WithContext(ctx).Infof("[rss_hub GetSubscribedChannelsByGroupId] query subscribe success: %v", len(channels))
	for _, cn := range channels {
		rv = append(rv, &RssChannelView{
			Channel: cn,
		})
	}

	return
}

package rss_pkg

import (
	"context"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

// SyncRssFeedNoNotice 同步所有频道
// 1. 获取所有频道
// 2. 遍历所有频道，检查频道是否更新
// 3. 如果更新，获取更新的内容，但是返回的数据不分组
func (repo *rssDomain) SyncRssFeedNoNotice(ctx context.Context) (updated map[int64]*RssChannelView, err error) {
	updated = make(map[int64]*RssChannelView)
	// 获取所有频道
	channels, err := repo.storage.GetSources(ctx)
	if err != nil {
		return
	}
	// 遍历所有源，获取每个channel对应的rss内容
	rssChannelView := make([]*RssChannelView, 0)
	for _, channel := range channels {
		var feed *gofeed.Feed
		// 从site获取rss内容
		feed, err = repo.rssHubClient.FetchFeed(rssHubMirrors[0], channel.RssHubFeedPath)
		if err != nil {
			return nil, err
		}
		rv := convertFeedToRssChannelView(0, channel.RssHubFeedPath, feed)
		rssChannelView = append(rssChannelView, rv)
	}
	// 检查频道是否更新
	for _, cv := range rssChannelView {
		var needUpdate bool
		needUpdate, err = repo.processRssChannelUpdate(ctx, cv.Channel)
		if err != nil {
			logrus.WithContext(ctx).Errorf("[rss_hub SyncRssFeedNoNotice] process rss cv update error: %v", err)
			err = nil
			continue
		}
		logrus.WithContext(ctx).Infof("[rss_hub SyncRssFeedNoNotice] cv %s, need update(real): %v", cv.Channel.RssHubFeedPath, needUpdate)
		needUpdate = true
		// 如果需要更新，更新content db
		if needUpdate {
			var updateChannelView = &RssChannelView{Channel: cv.Channel, Contents: []*RssContent{}}
			for _, content := range cv.Contents {
				content.RssFeedChannelId = cv.Channel.Id
				var existed bool
				existed, err = repo.processRssContentUpdate(ctx, content)
				if err != nil {
					logrus.WithContext(ctx).Errorf("[rss_hub SyncRssFeedNoNotice] upsert content error: %v", err)
					err = nil
					continue
				}
				if !existed {
					updateChannelView.Contents = append(updateChannelView.Contents, content)
					logrus.WithContext(ctx).Infof("[rss_hub SyncRssFeedNoNotice] cv %s, add new content: %v", cv.Channel.RssHubFeedPath, content.Title)
				}
			}
			updated[updateChannelView.Channel.Id] = updateChannelView
			logrus.WithContext(ctx).Infof("[rss_hub SyncRssFeedNoNotice] cv %s, new contents: %v", cv.Channel.RssHubFeedPath, len(updateChannelView.Contents))
		}
	}
	return
}

func (repo *rssDomain) processRssChannelUpdate(ctx context.Context, channel *RssFeedChannel) (needUpdate bool, err error) {
	var channelSrc *RssFeedChannel
	channelSrc, err = repo.storage.GetSourceByRssHubFeedLink(ctx, channel.RssHubFeedPath)
	if err != nil {
		return
	}
	if channelSrc == nil {
		logrus.WithContext(ctx).Errorf("[rss_hub SyncRssFeedNoNotice] channel not found: %v", channel.RssHubFeedPath)
		return
	}
	channel.Id = channelSrc.Id
	// 检查是否需要更新到db
	if channelSrc.IfNeedUpdate(channel) {
		needUpdate = true
		// 保存
		err = repo.storage.UpsertSource(ctx, channel)
		if err != nil {
			logrus.WithContext(ctx).Errorf("[rss_hub SyncRssFeedNoNotice] upsert source error: %v", err)
			return
		}
	}
	return
}

func (repo *rssDomain) processRssContentUpdate(ctx context.Context, content *RssContent) (existed bool, err error) {
	existed, err = repo.storage.IsContentHashIDExist(ctx, content.HashId)
	if err != nil {
		return
	}
	// 不需要更新&不需要发送
	if existed {
		return
	}
	// 保存
	err = repo.storage.UpsertContent(ctx, content)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub SyncRssFeedNoNotice] upsert content error: %v", err)
		return
	}
	return
}

// SyncJobTrigger 同步任务，按照群组订阅情况做好map切片
func (repo *rssDomain) SyncJobTrigger(ctx context.Context) (groupView map[int64][]*RssChannelView, err error) {
	groupView = make(map[int64][]*RssChannelView)
	// 获取所有Rss频道
	// 获取所有频道
	updatedChannelView, err := repo.SyncRssFeedNoNotice(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub SyncJobTrigger] sync rss feed error: %v", err)
		return
	}
	logrus.WithContext(ctx).Infof("[rss_hub SyncJobTrigger] updated channels: %v", len(updatedChannelView))
	subscribes, err := repo.storage.GetSubscribes(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub SyncJobTrigger] get subscribes error: %v", err)
		return
	}
	for _, subscribe := range subscribes {
		groupView[subscribe.GroupId] = append(groupView[subscribe.GroupId], updatedChannelView[subscribe.RssFeedChannelId])
	}
	return
}

//func (repo *rssDomain) processDiffChannel(ctx context.Context, channelSrc *RssFeedChannel, channelNew *RssFeedChannel) (err error) {
//	// 检查频道是否更新
//	if channelSrc.IfNeedUpdate(channelNew) {
//		// 更新频道信息
//		if err = repo.storage.UpsertSource(ctx, channelNew); err != nil {
//			return
//		}
//
//	}
//	// 检查是否更新
//	return
//}
//
//func (repo *rssDomain) processDiffContent() {
//	//if err = repo.storage.UpsertContent(ctx, rssItem); err != nil {
//	//	return
//	//}
//	return
//}

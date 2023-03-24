package rss_pkg

import (
	"context"
	"errors"
	"fmt"
	sql "github.com/FloatTech/sqlite"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

type RepoStorage interface {
	RepoContent
	RepoSource
	RepoSubscribe
	RepoMultiQuery
	initDB() error
}

// repoStorage db struct for rss
type repoStorage struct {
	db  sql.Sqlite
	orm *gorm.DB
}

// GetSubscribesBySource Impl
func (s *repoStorage) GetSubscribesBySource(ctx context.Context, feedPath string) ([]*RssSubscribe, error) {
	logrus.WithContext(ctx).Infof("[rss_hub GetSubscribesBySource] feedPath: %s", feedPath)
	//TODO implement me
	panic("implement me")
}

// GetIfExistedSubscribe Impl
func (s *repoStorage) GetIfExistedSubscribe(ctx context.Context, gid int64, feedPath string) (*RssSubscribe, bool, error) {
	rs := RssSubscribe{}
	err := s.orm.Model(&RssSubscribe{}).Joins(fmt.Sprintf("%s left join %s on %s.rss_feed_channel_id=%s.id", tableNameRssSubscribe, tableNameRssFeedChannel, tableNameRssSubscribe, tableNameRssFeedChannel)).
		Where(&RssSubscribe{GroupId: gid}, &RssFeedChannel{RssHubFeedPath: feedPath}).Select("rss_subscribe.*").First(&rs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		logrus.WithContext(ctx).Errorf("[rss_hub GetIfExistedSubscribe] error: %v", err)
		return nil, false, err
	}
	return &rs, true, nil
}

// initDB ...
func (s *repoStorage) initDB() (err error) {
	err = s.db.Open(time.Hour * 24)
	if err != nil {
		logrus.Errorf("[rss_hub NewRssDomain] open db error: %v", err)
		return
	}
	//if err = s.orm.Create(tableNameRssFeedChannel, &RssFeedChannel{}); err != nil {
	//	logrus.Errorf("[rss_hub NewRssDomain] Create db table RssFeedChannel error: %v", err)
	//	return
	//}
	//if err = s.db.Create(tableNameRssFeedContent, &RssContent{}); err != nil {
	//	logrus.Errorf("[rss_hub NewRssDomain] Create db table RssContent error: %v", err)
	//	return
	//}
	//if err = s.db.Create(tableNameRssSubscribe, &RssSubscribe{}); err != nil {
	//	logrus.Errorf("[rss_hub NewRssDomain] Create db table RssSubscribe error: %v", err)
	//	return
	//}
	s.orm, err = gorm.Open("sqlite3", s.db.DBPath)
	if err != nil {
		return
	}
	s.orm.LogMode(true)
	s.orm.AutoMigrate(&RssFeedChannel{}).AutoMigrate(&RssContent{}).AutoMigrate(&RssSubscribe{})
	s.orm.Table(tableNameRssSubscribe).AddUniqueIndex("uk_cid_gid", "rss_feed_channel_id", "group_id")
	return nil
}

// ==================== RepoSource ==================== [Start]

// UpsertSource Impl
func (s *repoStorage) UpsertSource(ctx context.Context, source *RssFeedChannel) (err error) {
	//rfc := RssFeedChannel{
	//	RssHubFeedPath: source.RssHubFeedPath,
	//	Title:          source.Title,
	//	ChannelDesc:    source.ChannelDesc,
	//	ImageUrl:       source.ImageUrl,
	//	Link:           source.Link,
	//	UpdatedParsed:  source.UpdatedParsed,
	//}
	// Update columns to default value on `id` conflict
	err = s.orm.Take(source, "rss_hub_feed_path = ?", source.RssHubFeedPath).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = s.orm.Create(source).Omit("id").Error
			if err != nil {
				logrus.WithContext(ctx).Errorf("[rss_hub] add source error: %v", err)
				return
			}
		}
		return
	}
	err = s.orm.Model(source).Updates(source).Omit("id").Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub] update source error: %v", err)
		return
	}
	logrus.Println("[rss_hub] add source success: ", source.Id)
	return nil
}

// GetSource Impl
func (s *repoStorage) GetSource(ctx context.Context, id int64) (source *RssFeedChannel, err error) {
	source = &RssFeedChannel{}
	err = s.db.Find(tableNameRssFeedChannel, source, fmt.Sprintf("id = %d", id))
	if err != nil {
		if errors.Is(err, sql.ErrNullResult) {
			return nil, errors.New("source not found")
		}
		logrus.WithContext(ctx).Errorf("[rss_hub] get source error: %v", err)
		return nil, err
	}
	return
}

// GetSources Impl
func (s *repoStorage) GetSources(ctx context.Context) (sources []*RssFeedChannel, err error) {
	sources = []*RssFeedChannel{}
	err = s.orm.Find(&sources, "id > 0").Error
	if err != nil {
		if errors.Is(err, sql.ErrNullResult) {
			return nil, errors.New("source not found")
		}
		logrus.WithContext(ctx).Errorf("[rss_hub] get sources error: %v", err)
		return
	}
	logrus.WithContext(ctx).Infof("[rss_hub] get sources success: %d", len(sources))
	return
}

// GetSourceByRssHubFeedLink Impl
func (s *repoStorage) GetSourceByRssHubFeedLink(ctx context.Context, rssHubFeedLink string) (source *RssFeedChannel, err error) {
	source = &RssFeedChannel{}
	err = s.db.Query(fmt.Sprintf("select * from %s where rss_hub_feed_path = '%s'", tableNameRssFeedChannel, rssHubFeedLink), source)
	if err != nil {
		if errors.Is(err, sql.ErrNullResult) {
			return nil, nil
		}
		logrus.WithContext(ctx).Errorf("[rss_hub] get source error: %v", err)
		return
	}
	return
}

// DeleteSource Impl
func (s *repoStorage) DeleteSource(ctx context.Context, id int64) (err error) {
	err = s.db.Del(tableNameRssFeedChannel, fmt.Sprintf("id = %d", id))
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub] storage.DeleteSource: %v", err)
		if errors.Is(err, sql.ErrNullResult) {
			return errors.New("source not found")
		}
		return
	}
	return nil
}

// ==================== RepoSource ==================== [End]

// ==================== RepoContent ==================== [Start]

// UpsertContent Impl
func (s *repoStorage) UpsertContent(ctx context.Context, content *RssContent) (err error) {
	// check params
	if content == nil {
		err = errors.New("content is nil")
		return
	}
	// check params.RssHubFeedPath and params.HashId
	if content.RssFeedChannelId < 0 || content.HashId == "" || content.Title == "" {
		err = errors.New("content.RssFeedChannelRoute or content.HashId is empty")
		return
	}
	err = s.orm.Create(content).Omit("id").Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub] storage.UpsertContent: %v", err)
		return
	}
	return
}

// DeleteSourceContents Impl
func (s *repoStorage) DeleteSourceContents(ctx context.Context, channelID int64) (rows int64, err error) {
	err = s.orm.Delete(&RssSubscribe{}).Where(&RssSubscribe{RssFeedChannelId: channelID}).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub] storage.DeleteSourceContents: %v", err)
		return
	}
	return
}

// IsContentHashIDExist Impl
func (s *repoStorage) IsContentHashIDExist(ctx context.Context, hashID string) (res bool, err error) {
	res = s.db.CanFind(tableNameRssFeedContent, fmt.Sprintf("hash_id = '%s'", hashID))
	logrus.WithContext(ctx).Debugf("[rss_hub] storage.IsContentHashIDExist: %v", res)
	return
}

// ==================== RepoContent ==================== [End]

// ==================== RepoSubscribe ==================== [Start]

// CreateSubscribe Impl
func (s *repoStorage) CreateSubscribe(ctx context.Context, gid, rssFeedChannelId int64) (err error) {
	// check subscribe
	if rssFeedChannelId < 0 || gid == 0 {
		err = errors.New("gid or rssFeedChannelId is empty")
		return
	}
	err = s.orm.Create(&RssSubscribe{GroupId: gid, RssFeedChannelId: rssFeedChannelId}).Omit("id").Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub] storage.CreateSubscribe: %v", err)
		return
	}
	return
}

// DeleteSubscribe Impl
func (s *repoStorage) DeleteSubscribe(ctx context.Context, gid int64, subscribeId int64) (err error) {
	err = s.orm.Delete(&RssSubscribe{}, "rss_feed_channel_id = ? and group_id = ?", subscribeId, gid).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rss_hub] storage.DeleteSubscribe error: %v", err)
		return
	}
	return
}

// GetSubscribeById Impl
func (s *repoStorage) GetSubscribeById(ctx context.Context, gid int64, subscribeId int64) (res *RssSubscribe, err error) {
	res = &RssSubscribe{}
	err = s.orm.First(res, &RssSubscribe{GroupId: gid, RssFeedChannelId: subscribeId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logrus.WithContext(ctx).Errorf("[rss_hub] storage.GetSubscribeById: %v", err)
		return nil, err
	}
	return
}

// GetSubscribedChannelsByGroupId Impl
func (s *repoStorage) GetSubscribedChannelsByGroupId(ctx context.Context, gid int64) (res []*RssFeedChannel, err error) {
	res = make([]*RssFeedChannel, 0)
	err = s.orm.Model(&RssFeedChannel{}).Joins(fmt.Sprintf("join %s on rss_feed_channel_id=%s.id", tableNameRssSubscribe, tableNameRssFeedChannel)).
		Where(&RssSubscribe{GroupId: gid}).
		Find(&res).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
			return
		}
		logrus.WithContext(ctx).Errorf("[rss_hub] storage.GetSubscribedChannelsByGroupId: %v", err)
		return
	}
	return
}

// GetSubscribes Impl
func (s *repoStorage) GetSubscribes(ctx context.Context) (res []*RssSubscribe, err error) {
	res = make([]*RssSubscribe, 0)
	err = s.orm.Find(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
			return
		}
		logrus.WithContext(ctx).Errorf("[rss_hub] storage.GetSubscribes: %v", err)
		return
	}
	return
}

// ==================== RepoSubscribe ==================== [End]

package domain

import (
	"context"
	"errors"
	"fmt"
	sql "github.com/FloatTech/sqlite"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

// RepoStorage 定义RepoStorage接口
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
	logrus.WithContext(ctx).Infof("[rsshub GetSubscribesBySource] feedPath: %s", feedPath)
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
		logrus.WithContext(ctx).Errorf("[rsshub GetIfExistedSubscribe] error: %v", err)
		return nil, false, err
	}
	return &rs, true, nil
}

// initDB ...
func (s *repoStorage) initDB() (err error) {
	err = s.db.Open(time.Hour * 24)
	if err != nil {
		logrus.Errorf("[rsshub NewRssDomain] open db error: %v", err)
		return
	}
	//if err = s.orm.Create(tableNameRssFeedChannel, &RssFeedChannel{}); err != nil {
	//	logrus.Errorf("[rsshub NewRssDomain] Create db table RssFeedChannel error: %v", err)
	//	return
	//}
	//if err = s.db.Create(tableNameRssFeedContent, &RssContent{}); err != nil {
	//	logrus.Errorf("[rsshub NewRssDomain] Create db table RssContent error: %v", err)
	//	return
	//}
	//if err = s.db.Create(tableNameRssSubscribe, &RssSubscribe{}); err != nil {
	//	logrus.Errorf("[rsshub NewRssDomain] Create db table RssSubscribe error: %v", err)
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
				logrus.WithContext(ctx).Errorf("[rsshub] add source error: %v", err)
				return
			}
		}
		return
	}
	err = s.orm.Model(source).Updates(source).Omit("id").Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rsshub] update source error: %v", err)
		return
	}
	logrus.Println("[rsshub] add source success: ", source.ID)
	return nil
}

// GetSource Impl
func (s *repoStorage) GetSource(ctx context.Context, fID int64) (source *RssFeedChannel, err error) {
	source = &RssFeedChannel{}
	err = s.db.Find(tableNameRssFeedChannel, source, fmt.Sprintf("id = %d", fID))
	if err != nil {
		if errors.Is(err, sql.ErrNullResult) {
			return nil, errors.New("source not found")
		}
		logrus.WithContext(ctx).Errorf("[rsshub] get source error: %v", err)
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
		logrus.WithContext(ctx).Errorf("[rsshub] get sources error: %v", err)
		return
	}
	logrus.WithContext(ctx).Infof("[rsshub] get sources success: %d", len(sources))
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
		logrus.WithContext(ctx).Errorf("[rsshub] get source error: %v", err)
		return
	}
	return
}

// DeleteSource Impl
func (s *repoStorage) DeleteSource(ctx context.Context, fID int64) (err error) {
	err = s.db.Del(tableNameRssFeedChannel, fmt.Sprintf("id = %d", fID))
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rsshub] storage.DeleteSource: %v", err)
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
	if content.RssFeedChannelID < 0 || content.HashId == "" || content.Title == "" {
		err = errors.New("content.RssFeedChannelRoute or content.HashId is empty")
		return
	}
	err = s.orm.Create(content).Omit("id").Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rsshub] storage.UpsertContent: %v", err)
		return
	}
	return
}

// DeleteSourceContents Impl
func (s *repoStorage) DeleteSourceContents(ctx context.Context, channelID int64) (rows int64, err error) {
	err = s.orm.Delete(&RssSubscribe{}).Where(&RssSubscribe{RssFeedChannelID: channelID}).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rsshub] storage.DeleteSourceContents: %v", err)
		return
	}
	return
}

// IsContentHashIDExist Impl
func (s *repoStorage) IsContentHashIDExist(ctx context.Context, hashID string) (res bool, err error) {
	res = s.db.CanFind(tableNameRssFeedContent, fmt.Sprintf("hash_id = '%s'", hashID))
	logrus.WithContext(ctx).Debugf("[rsshub] storage.IsContentHashIDExist: %v", res)
	return
}

// ==================== RepoContent ==================== [End]

// ==================== RepoSubscribe ==================== [Start]

// CreateSubscribe Impl
func (s *repoStorage) CreateSubscribe(ctx context.Context, gid, rssFeedChannelID int64) (err error) {
	// check subscribe
	if rssFeedChannelID < 0 || gid == 0 {
		err = errors.New("gid or rssFeedChannelId is empty")
		return
	}
	err = s.orm.Create(&RssSubscribe{GroupId: gid, RssFeedChannelID: rssFeedChannelID}).Omit("id").Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rsshub] storage.CreateSubscribe: %v", err)
		return
	}
	return
}

// DeleteSubscribe Impl
func (s *repoStorage) DeleteSubscribe(ctx context.Context, gid int64, subscribeID int64) (err error) {
	err = s.orm.Delete(&RssSubscribe{}, "rss_feed_channel_id = ? and group_id = ?", subscribeID, gid).Error
	if err != nil {
		logrus.WithContext(ctx).Errorf("[rsshub] storage.DeleteSubscribe error: %v", err)
		return
	}
	return
}

// GetSubscribeById Impl
func (s *repoStorage) GetSubscribeById(ctx context.Context, gid int64, subscribeID int64) (res *RssSubscribe, err error) {
	res = &RssSubscribe{}
	err = s.orm.First(res, &RssSubscribe{GroupId: gid, RssFeedChannelID: subscribeID}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logrus.WithContext(ctx).Errorf("[rsshub] storage.GetSubscribeById: %v", err)
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
		logrus.WithContext(ctx).Errorf("[rsshub] storage.GetSubscribedChannelsByGroupId: %v", err)
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
		logrus.WithContext(ctx).Errorf("[rsshub] storage.GetSubscribes: %v", err)
		return
	}
	return
}

// ==================== RepoSubscribe ==================== [End]

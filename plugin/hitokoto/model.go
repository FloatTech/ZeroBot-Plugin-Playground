package hitokoto

import (
	"encoding/json"
	"os"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

const (
	chinesebqbURL = "https://cdn.jsdelivr.net/gh/zhaoolee/yiyan_spider@master/yiyan_data.json"
)

// hdb 表情包数据库全局变量
var hdb *hitokotodb

// hitokotodb 表情包数据库
type hitokotodb gorm.DB

// initialize 初始化
func initialize(dbpath string) *hitokotodb {
	var err error
	if _, err = os.Stat(dbpath); err != nil || os.IsNotExist(err) {
		// 生成文件
		f, err := os.Create(dbpath)
		if err != nil {
			return nil
		}
		defer f.Close()
	}
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	hdb := (*hitokotodb)(db)
	err = db.AutoMigrate(&hitokoto{})
	if err != nil {
		panic(err)
	}
	var count int64
	err = db.Model(&hitokoto{}).Count(&count).Error
	if err != nil {
		panic(err)
	}
	if count == 0 {
		err = hdb.truncateAndInsert()
		if err != nil {
			panic(err)
		}
	}
	err = db.Model(&hitokoto{}).Count(&count).Error
	if err != nil {
		panic(err)
	}
	logrus.Infoln("[hitokoto]加载", count, "条一言")
	return (*hitokotodb)(db)
}

type hitokoto struct {
	ID         int    `json:"id" gorm:"column:id"`
	Hitokoto   string `json:"hitokoto" gorm:"column:hitokoto"`
	Type       string `json:"type" gorm:"column:type"`
	From       string `json:"from" gorm:"column:from"`
	FromWho    string `json:"from_who" gorm:"column:from_who"`
	Creator    string `json:"creator" gorm:"column:creator"`
	CreatorUID int    `json:"creator_uid" gorm:"column:creator_uid"`
	Reviewer   int    `json:"reviewer" gorm:"column:reviewer"`
	UUID       string `json:"uuid" gorm:"column:uuid"`
	CreatedAt  string `json:"created_at" gorm:"column:created_at"`
	Category   string `json:"catogory" gorm:"column:category"`
}

// TableName 表名
func (hitokoto) TableName() string {
	return "hitokoto"
}

func (hdb *hitokotodb) getByKey(key string) (b []hitokoto, err error) {
	db := (*gorm.DB)(hdb)
	err = db.Where("hitokoto like ?", "%"+key+"%").Find(&b).Error
	return
}

func (hdb *hitokotodb) truncateAndInsert() (err error) {
	db := (*gorm.DB)(hdb)
	var (
		h    hitokoto
		data []byte
	)
	data, err = web.GetData(chinesebqbURL)
	if err != nil {
		return
	}
	err = db.Exec("drop table hitokoto").Error
	if err != nil {
		return
	}
	err = db.AutoMigrate(&hitokoto{})
	if err != nil {
		return
	}
	hlist := make([]hitokoto, 0, 100)
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		value.ForEach(func(vk, v gjson.Result) bool {
			err := json.Unmarshal(binary.StringToBytes(v.String()), &h)
			if err != nil {
				return true
			}
			h.Category = key.String()
			hlist = append(hlist, h)
			if len(hlist) >= 100 {
				err = db.Create(&hlist).Error
				if err != nil {
					return true
				}
				hlist = hlist[:0]
			}
			return true
		})
		return true
	})
	err = db.Create(&hlist).Error
	return
}

type result struct {
	Category string
	Count    int
}

func (hdb *hitokotodb) getAllCategory() (results []result, err error) {
	db := (*gorm.DB)(hdb)
	err = db.Table("hitokoto").Select("category, count(1) as count").Group("category").Scan(&results).Error
	return
}

func (hdb *hitokotodb) getByCategory(category string) (h []hitokoto, err error) {
	db := (*gorm.DB)(hdb)
	err = db.Where("category = ?", category).Find(&h).Error
	return
}

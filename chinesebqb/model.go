package chinesebqb

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/FloatTech/floatbox/web"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	chinesebqbURL = "https://www.v2fy.com/asset/0i/ChineseBQB/chinesebqb_v2fy.json"
)

type chineseReply struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
	Data   []struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		URL      string `json:"url"`
	} `json:"data"`
}

// bdb 表情包数据库全局变量
var bdb *bqbdb

// bqbdb 表情包数据库
type bqbdb gorm.DB

// initialize 初始化
func initialize(dbpath string) *bqbdb {
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
	bdb := (*bqbdb)(db)
	err = db.AutoMigrate(&bqb{})
	if err != nil {
		panic(err)
	}
	var count int64
	err = db.Model(&bqb{}).Count(&count).Error
	if err != nil {
		panic(err)
	}
	if count == 0 {
		err = bdb.truncateAndInsert()
		if err != nil {
			panic(err)
		}
	}
	err = db.Model(&bqb{}).Count(&count).Error
	if err != nil {
		panic(err)
	}
	logrus.Infoln("[chinesebqb]加载", count, "条表情包")
	return (*bqbdb)(db)
}

// bqb 表情包数据
type bqb struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Name     string `gorm:"column:name"`
	Category string `gorm:"column:category"`
	URL      string `gorm:"column:url"`
}

// TableName 表名
func (bqb) TableName() string {
	return "bqb"
}

func (bdb *bqbdb) getByKey(key string) (b []bqb, err error) {
	db := (*gorm.DB)(bdb)
	err = db.Where("name like ?", "%"+key+"%").Find(&b).Error
	return
}

func (bdb *bqbdb) truncateAndInsert() (err error) {
	db := (*gorm.DB)(bdb)
	var (
		data []byte
		r    chineseReply
	)
	data, err = web.GetData(chinesebqbURL)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return
	}
	err = db.Exec("drop table bqb").Error
	if err != nil {
		return
	}
	err = db.AutoMigrate(&bqb{})
	if err != nil {
		return
	}
	blist := make([]bqb, 0, 100)
	for _, v := range r.Data {
		b := bqb{
			Category: v.Category,
			URL:      v.URL,
		}
		_, back, flag := strings.Cut(v.Name, "-")
		if flag {
			b.Name = back
		} else {
			b.Name = v.Name
		}
		if b.URL != "" {
			blist = append(blist, b)
		}
		if len(blist) >= 100 {
			err = db.Create(&blist).Error
			if err != nil {
				return
			}
			blist = blist[:0]
		}
	}
	err = db.Create(&blist).Error
	return
}

type result struct {
	Category string
	Count    int
}

func (bdb *bqbdb) getAllCategory() (results []result, err error) {
	db := (*gorm.DB)(bdb)
	err = db.Table("bqb").Select("category, count(1) as count").Group("category").Scan(&results).Error
	return
}

func (bdb *bqbdb) getByCategory(category string) (b []bqb, err error) {
	db := (*gorm.DB)(bdb)
	err = db.Where("category = ?", category).Find(&b).Error
	return
}

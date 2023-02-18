package asoularticle

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/FloatTech/floatbox/web"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	articlesURL = "https://asoul.icu/v/articles?pageNum=%v&pageSize=48"
	detailURL   = "https://asoul.icu/v/articles/%v"
)

type articleReply struct {
	Articles []Articles `json:"articles"`
}

// adb 小作文数据库全局变量
var adb *articledb

// articledb 小作文数据库
type articledb gorm.DB

// initialize 初始化
func initialize(dbpath string) *articledb {
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
	adb := (*articledb)(db)
	err = db.AutoMigrate(&Article{})
	if err != nil {
		panic(err)
	}
	var count int64
	err = db.Model(&Article{}).Count(&count).Error
	if err != nil {
		panic(err)
	}
	if count == 0 {
		err = adb.truncateAndInsert()
		if err != nil {
			panic(err)
		}
	}
	err = db.Model(&Article{}).Count(&count).Error
	if err != nil {
		panic(err)
	}
	logrus.Infoln("[asoularticle]加载", count, "条小作文")
	return (*articledb)(db)
}

// Article 小作文数据
type Article struct {
	ID         int    `gorm:"primary_key;AUTO_INCREMENT"`
	Title      string `gorm:"column:title"`
	Author     string `gorm:"column:author"`
	Brief      string `gorm:"column:brief"`
	CreateTime int64  `gorm:"column:create_time"`
	Tags       string `gorm:"column:tags"`
}

// Articles 小作文数据
type Articles struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Brief      string   `json:"brief"`
	CreateTime int64    `json:"createTime"`
	Tags       []string `json:"tags"`
}

// TableName 表名
func (Article) TableName() string {
	return "article"
}

func (adb *articledb) truncateAndInsert() (err error) {
	db := (*gorm.DB)(adb)
	var (
		data    []byte
		r       articleReply
		pageNum = 1
	)
	err = db.Exec("drop table article").Error
	if err != nil {
		return
	}
	err = db.AutoMigrate(&Article{})
	if err != nil {
		return
	}
	data, err = web.GetData(fmt.Sprintf(articlesURL, pageNum))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return
	}
	alist := make([]Article, 0, 100)
	for len(r.Articles) > 0 {
		for _, v := range r.Articles {
			alist = append(alist, Article{
				ID:         v.ID,
				Title:      v.Title,
				Author:     v.Author,
				Brief:      v.Brief,
				CreateTime: v.CreateTime,
				Tags:       strings.Join(v.Tags, "|"),
			})
		}
		pageNum++
		data, err = web.GetData(fmt.Sprintf(articlesURL, pageNum))
		if err != nil {
			return
		}
		err = json.Unmarshal(data, &r)
		if err != nil {
			return
		}
	}
	partition := Partition(alist, 100)
	for _, v := range partition {
		err = db.Create(&v).Error
		if err != nil {
			return
		}
	}
	return
}

func (adb *articledb) randomArticle() (a Article, err error) {
	db := (*gorm.DB)(adb)
	var count int64
	err = db.Model(&Article{}).Count(&count).Offset(rand.Intn(int(count))).Take(&a).Error
	return
}

// Partition 分组函数
func Partition[T any](list []T, size int) (plist [][]T) {
	plen := int(math.Ceil(float64(len(list)) / float64(size)))
	plist = make([][]T, plen)
	for i := 0; i < plen; i++ {
		if i == plen-1 {
			plist[i] = list[i*size:]
			break
		}
		plist[i] = list[i*size : (i+1)*size]
	}
	return
}

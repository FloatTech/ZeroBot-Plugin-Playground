package qzone

import (
	"os"

	_ "github.com/fumiama/sqlite3" // use sql
	"github.com/jinzhu/gorm"
)

// qdb qq空间数据库全局变量
var qdb *qzonedb

// qzonedb qq空间数据库结构体
type qzonedb gorm.DB

// initialize 初始化
func initialize(dbpath string) *qzonedb {
	var err error
	if _, err = os.Stat(dbpath); err != nil || os.IsNotExist(err) {
		// 生成文件
		f, err := os.Create(dbpath)
		if err != nil {
			return nil
		}
		defer f.Close()
	}
	qdb, err := gorm.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}
	qdb.AutoMigrate(&qzoneConfig{}).AutoMigrate(&emotion{})
	return (*qzonedb)(qdb)
}

// qzoneConfig qq空间初始化信息
type qzoneConfig struct {
	ID      uint   `gorm:"primary_key;AUTO_INCREMENT"`
	QQ      int64  `gorm:"column:qq;unique;not null"`
	Skey    string `gorm:"column:skey"`
	Pskey   string `gorm:"column:pskey"`
	Cookies string `gorm:"column:cookies;type:varchar(1024)"`
}

// TableName 表名
func (qzoneConfig) TableName() string {
	return "qzone_config"
}

func (qdb *qzonedb) insertOrUpdate(qq int64, skey, pskey, cookies string) (err error) {
	db := (*gorm.DB)(qdb)
	qc := qzoneConfig{
		QQ:      qq,
		Skey:    skey,
		Pskey:   pskey,
		Cookies: cookies,
	}
	var oqc qzoneConfig
	err = db.Take(&oqc, "qq = ?", qc.QQ).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = db.Create(&qc).Error
		}
		return
	}
	err = db.Model(&oqc).Updates(qc).Error
	return
}

func (qdb *qzonedb) getByUin(qq int64) (qc qzoneConfig, err error) {
	db := (*gorm.DB)(qdb)
	err = db.Take(&qc, "qq = ?", qq).Error
	return
}

// emotion 说说信息
type emotion struct {
	gorm.Model
	QQ     int64  `gorm:"column:qq"`
	Msg    string `gorm:"column:msg"`
	Status int    `gorm:"column:status"` // 1-审核中,2-同意,3-拒绝
	Tag    string `gorm:"column:tag"`
}

// TableName 表名
func (emotion) TableName() string {
	return "emotion"
}

func (qdb *qzonedb) saveEmotion(e emotion) (id int64, err error) {
	db := (*gorm.DB)(qdb)
	err = db.Create(&e).Error
	id = int64(e.ID)
	return
}

func (qdb *qzonedb) getEmotionByID(id int64) (e emotion, err error) {
	db := (*gorm.DB)(qdb)
	err = db.Take(&e, "id = ?", id).Error
	return
}

func (qdb *qzonedb) getEmotionByStatus(status int) (el []emotion, err error) {
	db := (*gorm.DB)(qdb)
	if status == 0 {
		err = db.Find(&el).Error
		return
	}
	err = db.Find(&el, "status = ?", status).Error
	return
}

func (qdb *qzonedb) updateEmotionStatusByID(id int64, status int) (err error) {
	db := (*gorm.DB)(qdb)
	err = db.Model(&emotion{}).Where("id = ?", id).Update("status", status).Error
	return
}

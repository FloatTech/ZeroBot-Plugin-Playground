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
	qdb.AutoMigrate(&qzoneConfig{})
	return (*qzonedb)(qdb)
}

// qzoneConfig qq空间初始化信息
type qzoneConfig struct {
	ID      uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Uin     int64  `gorm:"column:uin;unique;not null"`
	Skey    string `gorm:"column:skey"`
	Pskey   string `gorm:"column:pskey"`
	Cookies string `gorm:"column:cookies;type:varchar(1024)"`
}

// TableName 表名
func (qzoneConfig) TableName() string {
	return "qzone_config"
}

// sleep 更新睡眠时间
func (qdb *qzonedb) insertOrUpdate(uin int64, skey, pskey, cookies string) (err error) {
	db := (*gorm.DB)(qdb)
	qc := qzoneConfig{
		Uin:     uin,
		Skey:    skey,
		Pskey:   pskey,
		Cookies: cookies,
	}
	var oqc qzoneConfig
	err = db.Take(&oqc, "uin = ?", qc.Uin).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = db.Create(&qc).Error
		}
		return
	}
	err = db.Model(&oqc).Updates(qc).Error
	return
}

// getUp 更新起床时间
func (qdb *qzonedb) getByUin(uin int64) (qc qzoneConfig, err error) {
	db := (*gorm.DB)(qdb)
	err = db.Take(&qc, "uin = ?", uin).Error
	return
}

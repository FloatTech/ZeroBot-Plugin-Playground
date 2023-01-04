package qqci

import (
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var adb *appdb

type appdb gorm.DB

type application struct {
	ID        int64  `gorm:"column:id;primary"`
	Appname   string `gorm:"column:appname;unique_index" flag:"a"`
	Gitrepo   string `gorm:"column:gitrepo" flag:"r"`
	Gitbranch string `gorm:"-" flag:"b"`
	Directory string `gorm:"column:directory" flag:"dir"`
	Command   string `gorm:"column:command" flag:"cmd"`
	Action    string `gorm:"-" flag:"act"`
	Upload    string `gorm:"-" flag:"u"`
	Folder    string `gorm:"column:folder" flag:"f"`
	GroupID   int    `gorm:"column:group_id" flag:"gid"`
}

type folder struct {
	FolderID   string `json:"folder_id"`
	FolderName string `json:"folder_name"`
}

// TableName ...
func (application) TableName() string {
	return "application"
}

func initialize(dbpath string) *appdb {
	var err error
	if _, err = os.Stat(dbpath); err != nil || os.IsNotExist(err) {
		// 生成文件
		f, err := os.Create(dbpath)
		if err != nil {
			return nil
		}
		defer f.Close()
	}
	adb, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = adb.AutoMigrate(&application{})
	if err != nil {
		panic(err)
	}
	return (*appdb)(adb)
}

func (adb *appdb) insert(app *application) error {
	db := (*gorm.DB)(adb)
	return db.Model(&application{}).Create(app).Error
}

func (adb *appdb) update(app *application) error {
	db := (*gorm.DB)(adb)
	return db.Model(&application{}).Where("appname = ?", app.Appname).Updates(*app).Error
}

func (adb *appdb) get(app *application) (res application, err error) {
	db := (*gorm.DB)(adb)
	err = db.Model(&application{}).First(&res, "appname = ?", app.Appname).Error
	return
}

func (adb *appdb) getApp(app *application) (a application, err error) {
	db := (*gorm.DB)(adb)
	if app.Appname != "" {
		err = db.Model(&application{}).First(&a, "appname = ?", app.Appname).Error
	} else if app.Gitrepo != "" {
		err = db.Model(&application{}).First(&a, "gitrepo = ?", app.Gitrepo).Error
	}
	return
}

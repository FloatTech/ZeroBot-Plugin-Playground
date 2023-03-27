package fgopickup

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"os"
)

var _dbEngine *orm

type orm struct {
	*gorm.DB
}

func getOrmEngine() *orm {
	return _dbEngine
}

func initialize(dbpath string) *gorm.DB {
	var err error
	if _, err = os.Stat(dbpath); err != nil || os.IsNotExist(err) {
		f, err := os.Create(dbpath)
		if err != nil {
			return nil
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(f)
	}
	gdb, err := gorm.Open(sqlite.Open(dbpath))
	if err != nil {
		panic(err)
	}
	//gdb.AutoMigrate(&pickup{}, &pickupServant{}, &servant{})
	orm := new(orm)
	orm.DB = gdb
	_dbEngine = orm
	return gdb
}

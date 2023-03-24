package fgopickup

import (
	"github.com/sirupsen/logrus"
	"time"
)

type dao struct {
	DbEngine *orm
}

func (d *dao) listPickup() *[]pickup {
	pickup := make([]pickup, 0)
	unixTime := time.Now().Unix()
	err := d.DbEngine.Where("end_time >= ?", unixTime).Find(&pickup).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return &pickup
}

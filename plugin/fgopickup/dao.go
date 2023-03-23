package fgopickup

import (
	"github.com/sirupsen/logrus"
)

type dao struct {
	DbEngine *orm
}

func (d *dao) listPickup() *[]pickup {
	pickup := make([]pickup, 0)
	err := d.DbEngine.Find(&pickup).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return &pickup
}

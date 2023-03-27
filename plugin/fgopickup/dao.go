package fgopickup

import (
	"github.com/sirupsen/logrus"
	"time"
)

type dao struct {
	DBEngine *orm
}

func (d *dao) listPickup() *[]pickup {
	pickup := make([]pickup, 0)
	unixTime := time.Now().Unix()
	err := d.DBEngine.Where("end_time >= ?", unixTime).Find(&pickup).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return &pickup
}

func (d *dao) selectPickup(pickupId int) pickup {
	pickup := pickup{}
	err := d.DBEngine.First(&pickup, pickupId).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return pickup
}

func (d *dao) selectPickupServantIds(pickupId int) []int {
	ids := make([]int, 0)
	err := d.DBEngine.Model(pickupServant{}).Select("servant_id").Where("pickup_id = ?", pickupId).Find(&ids).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return ids
}

func (d *dao) selectServantsByIds(ids []int) *[]servant {
	servants := make([]servant, 0)
	err := d.DBEngine.Find(&servants, ids).Error
	if err != nil {
		logrus.Debugln(err)
	}

	return &servants
}

package fgopickup

import (
	"github.com/sirupsen/logrus"
	"time"
)

type dao struct {
	DBEngine *orm
}

func (d *dao) listPickup() (*[]pickup, error) {
	pickup := make([]pickup, 0)
	unixTime := time.Now().Unix()
	err := d.DBEngine.Where("end_time >= ?", unixTime).Find(&pickup).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return &pickup, err
}

func (d *dao) selectPickup(pickupID int) (pickup, error) {
	pickup := pickup{}
	err := d.DBEngine.First(&pickup, pickupID).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return pickup, err
}

func (d *dao) selectServantIDsByPickupID(pickupID int) ([]int, error) {
	ids := make([]int, 0)
	err := d.DBEngine.Model(pickupServant{}).Select("servant_id").Where("pickup_id = ?", pickupID).Find(&ids).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return ids, err
}

func (d *dao) selectServantsByIDs(ids []int) (*[]servant, error) {
	servants := make([]servant, 0)
	err := d.DBEngine.Find(&servants, ids).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return &servants, err
}

func (d *dao) listServants(page int) (*[]servant, error) {
	pageSize := 50
	servants := make([]servant, 0)
	err := d.DBEngine.Offset(pageSize * (page - 1)).Limit(pageSize).Order("id desc").Find(&servants).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return &servants, err
}

func (d *dao) selectPickupIDsByServantID(id int) ([]int, error) {
	pickupIDs := make([]int, 0)
	err := d.DBEngine.Model(pickupServant{}).Select("pickup_id").Where("servant_id = ?", id).Find(&pickupIDs).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return pickupIDs, err
}

func (d *dao) selectPickupsByIDs(ids []int) (*[]pickup, error) {
	pickups := make([]pickup, 0)
	unixTime := time.Now().Unix()
	err := d.DBEngine.Where("end_time >= ?", unixTime).Find(&pickups, ids).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return &pickups, err
}

func (d *dao) selectServant(id int) (servant, error) {
	servant := servant{}
	err := d.DBEngine.First(&servant, id).Error
	if err != nil {
		logrus.Debugln(err)
	}
	return servant, err
}

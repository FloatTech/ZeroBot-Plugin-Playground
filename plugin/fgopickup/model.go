package fgopickup

type pickup struct {
	Id        int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name      string `gorm:"column:name"`
	StartTime int64  `gorm:"column:start_time"`
	EndTime   int64  `gorm:"column:end_time"`
	Banner    string `gorm:"column:banner"`
}

type pickupServant struct {
	Id        int `gorm:"primary_key;AUTO_INCREMENT"`
	PickupId  int `gorm:"column:pickup_id"`
	ServantId int `gorm:"column:servant_id"`
}

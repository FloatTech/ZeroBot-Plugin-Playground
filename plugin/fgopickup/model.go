package fgopickup

type pickup struct {
	ID        int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name      string `gorm:"column:name"`
	StartTime int64  `gorm:"column:start_time"`
	EndTime   int64  `gorm:"column:end_time"`
	Banner    string `gorm:"column:banner"`
}

type pickupServant struct {
	ID        int `gorm:"primary_key;AUTO_INCREMENT"`
	PickupID  int `gorm:"column:pickup_id"`
	ServantID int `gorm:"column:servant_id"`
}

type servant struct {
	ID     int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name   string `gorm:"column:name"`
	Avatar string `gorm:"column:avatar"`
}

type pickupDetailRes struct {
	Pickup   pickup
	Servants []servant
}

type servantPickupsRes struct {
	ServantName string
	Pickup      []pickup
}

package fgopickup

type pickup struct {
	Id        int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name      string `gorm:"column:name"`
	StartTime string `gorm:"column:start_time"`
	EndTime   string `gorm:"column:end_time"`
	Banner    string `gorm:"column:banner"`
}

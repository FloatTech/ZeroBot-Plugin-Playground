package fgopickup

import (
	"fmt"
)

type dao struct {
	DbEngine *orm
}

func (d *dao) List() *[]pickup {
	pickup := make([]pickup, 0)
	err := d.DbEngine.Find(&pickup).Error
	if err == nil {
		fmt.Println(err)
	}
	fmt.Println(pickup)
	return &pickup
}

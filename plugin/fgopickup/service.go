package fgopickup

import "fmt"

type service struct {
}

func (s *service) getPickups() (*[]pickup, error) {
	dao := dao{DBEngine: getOrmEngine()}
	list, err := dao.listPickup()
	return list, err
}

func (s *service) getPickupDetail(pickupID int) (pickupDetailRes, error) {
	dao := dao{DBEngine: getOrmEngine()}
	pickup, err := dao.selectPickup(pickupID)
	servantIds, err := dao.selectPickupServantIds(pickupID)
	servants, err := dao.selectServantsByIds(servantIds)
	return pickupDetailRes{
		Pickup:   pickup,
		Servants: *servants,
	}, err
}

func (s *service) getPickup(pickupID int) (pickup, error) {
	dao := dao{DBEngine: getOrmEngine()}
	pickup, err := dao.selectPickup(pickupID)
	return pickup, err
}

func (s *service) getPickupTimeGap(id int) int {
	fmt.Println(id)
	return 0
}

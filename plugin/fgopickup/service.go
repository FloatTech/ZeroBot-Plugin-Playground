package fgopickup

type service struct {
}

func (s *service) getPickups() *[]pickup {
	dao := dao{DBEngine: getOrmEngine()}
	list := dao.listPickup()
	return list
}

func (s *service) getPickupDetail(pickupId int) pickupDetailRes {
	dao := dao{DBEngine: getOrmEngine()}
	pickup := dao.selectPickup(pickupId)
	servantIds := dao.selectPickupServantIds(pickupId)
	servants := dao.selectServantsByIds(servantIds)
	return pickupDetailRes{
		Pickup:   pickup,
		Servants: *servants,
	}
}

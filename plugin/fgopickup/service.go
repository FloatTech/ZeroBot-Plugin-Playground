package fgopickup

type service struct {
}

func (s *service) getPickups() []pickup {
	dao := dao{DbEngine: getOrmEngine()}
	list := dao.listPickup()
	return *list
}

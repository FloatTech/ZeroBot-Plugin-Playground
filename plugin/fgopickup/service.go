package fgopickup

type service struct {
}

func (s *service) getPickups() []pickup {
	dao := dao{DBEngine: getOrmEngine()}
	list := dao.listPickup()
	return *list
}

package fgopickup

import "time"

type service struct {
}

func (s *service) getPickups() (*[]pickup, error) {
	dao := dao{DBEngine: getOrmEngine()}
	list, err := dao.listPickup()
	return list, err
}

func (s *service) getPickupDetail(pickupID int) (pickupDetailRes, error) {
	dao := dao{DBEngine: getOrmEngine()}
	res := pickupDetailRes{}
	pickup, err := dao.selectPickup(pickupID)
	if err != nil {
		return res, err
	}
	servantIds, err := dao.selectPickupServantIds(pickupID)
	if err != nil {
		return res, err
	}
	servants, err := dao.selectServantsByIds(servantIds)
	if err != nil {
		return res, err
	}
	res.Pickup = pickup
	res.Servants = *servants
	return res, err
}

func (s *service) getPickupTimeGap(pickupID int) (int, error) {
	dao := dao{DBEngine: getOrmEngine()}
	pickup, err := dao.selectPickup(pickupID)
	startTime := time.Now().Unix()
	endTime := pickup.StartTime
	days := getDiffDaysBySeconds(startTime, endTime)
	return days, err
}

func getDiffDays(start, end time.Time) int {
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.Local)

	return int(end.Sub(start).Hours() / 24)
}

func getDiffDaysBySeconds(start, end int64) int {
	return getDiffDays(time.Unix(start, 0), time.Unix(end, 0))
}

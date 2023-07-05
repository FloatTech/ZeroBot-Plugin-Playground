package repo

import (
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/types"
	"github.com/sirupsen/logrus"
)

type GlobalService struct{}

// ====== global ====

func (GlobalService) NewModel() interface{} {
	return &Global{}
}

func (GlobalService) Get(id string) interface{} {
	return GetGlobal()
}

func (GlobalService) Find(model interface{}) types.Page {
	return types.Page{}
}

func (GlobalService) Edit(model interface{}) bool {
	global, ok := model.(*Global)
	if !ok {
		return false
	}
	err := InsertGlobal(*global)
	if err != nil {
		logrus.Error(err)
		return false
	} else {
		return true
	}
}

func (GlobalService) Del(key string) bool {
	return false
}

// ====== end ====

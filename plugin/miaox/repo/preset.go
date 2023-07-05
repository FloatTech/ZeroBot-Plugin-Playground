package repo

import (
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/types"
	sql "github.com/FloatTech/sqlite"
	"github.com/sirupsen/logrus"
)

type PresetService struct{}

// ====== global ====

func (PresetService) NewModel() interface{} {
	return &PresetScene{}
}

func (PresetService) Get(id string) interface{} {
	return GetPresetScene(id)
}

func (PresetService) Find(model interface{}) types.Page {
	objptr, ok := model.(*PresetScene)
	if !ok {
		return types.Page{}
	}
	condition := BuildCondition(*objptr)
	objptrs, err := sql.FindAll[PresetScene](cmd.sql, "preset_scene", condition)
	if err != nil {
		logrus.Error(err)
		return types.Page{}
	}

	total, err := cmd.Count("preset_scene", condition)
	if err != nil {
		logrus.Error(err)
		return types.Page{}
	}

	newTokens := make([]interface{}, len(objptrs))
	for i, t := range objptrs {
		newTokens[i] = t
	}
	return types.Page{
		Total: total,
		List:  newTokens,
	}
}

func (PresetService) Edit(model interface{}) bool {
	token, ok := model.(*PresetScene)
	if !ok {
		return false
	}
	cmd.Lock()
	defer cmd.Unlock()

	if cmd.sql.Insert("preset_scene", token) != nil {
		return false
	}
	return true
}

func (PresetService) Del(key string) bool {
	if err := cmd.sql.Del("preset_scene", "where key = '"+key+"'"); err != nil {
		logrus.Error(err)
		return false
	}
	return true
}

// ====== end ====

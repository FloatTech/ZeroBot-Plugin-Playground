package repo

import (
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/types"
	sql "github.com/FloatTech/sqlite"
	"github.com/sirupsen/logrus"
)

type TokenService struct{}

// ====== global ====

func (TokenService) NewModel() interface{} {
	return &Token{}
}

func (TokenService) Get(id string) interface{} {
	return GetToken(id)
}

func (TokenService) Find(model interface{}) types.Page {
	token, ok := model.(*Token)
	if !ok {
		return types.Page{}
	}
	condition := BuildCondition(*token)
	tokens, err := sql.FindAll[Token](cmd.sql, "token", condition)
	if err != nil {
		logrus.Error(err)
		return types.Page{}
	}

	total, err := cmd.Count("token", condition)
	if err != nil {
		logrus.Error(err)
		return types.Page{}
	}

	newTokens := make([]interface{}, len(tokens))
	for i, t := range tokens {
		newTokens[i] = t
	}
	return types.Page{
		Total: total,
		List:  newTokens,
	}
}

func (TokenService) Edit(model interface{}) bool {
	token, ok := model.(*Token)
	if !ok {
		return false
	}
	cmd.Lock()
	defer cmd.Unlock()

	if cmd.sql.Insert("token", token) != nil {
		return false
	}
	return true
}

func (TokenService) Del(key string) bool {
	if err := cmd.sql.Del("token", "where key = '"+key+"'"); err != nil {
		logrus.Error(err)
		return false
	}
	return true
}

// ====== end ====

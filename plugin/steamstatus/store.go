package steamstatus

import (
	"fmt"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
)

const TableName = "steamstatus" // 存储数据库名称

// player 用户状态存储结构体
type player struct {
	SteamId       string `json:"steamid"`       // 监听用户标识ID
	PersonaName   string `json:"personaname"`   // 用户昵称
	Target        string `json:"target"`        // 信息推送群组
	GameId        string `json:"gameid"`        // 游戏ID
	GameExtraInfo string `json:"gameextrainfo"` // 游戏信息
	LastUpdate    int64  `json:"lastupdate"`    // 更新时间
}

// streamDB 继承方法的存储结构
type streamDB struct {
	db *sql.Sqlite
	sync.RWMutex
	Table string // 表名称
}

var database *streamDB // 持久化数据库对象，用于监听器使用

// initStore 初始化数据库
func initStore() error { // 初始化数据库信息
	database = &streamDB{
		db:    &sql.Sqlite{},
		Table: TableName,
	}
	database.db.DBPath = engine.DataFolder() + "steamstatus.db"

	// 创建数据库链接
	err := database.db.Open(time.Hour * 24)
	if err != nil {
		return err
	}
	return database.db.Create(database.Table, &player{})
}

// update 如果主键不存在则插入一条新的数据，如果主键存在直接复写
func (sql *streamDB) update(dbInfo player) error {
	sql.Lock()
	defer sql.Unlock()
	return sql.db.Insert(database.Table, &dbInfo)
}

// find 根据主键查信息
func (sql *streamDB) find(steamId string) (dbInfo player, err error) {
	sql.Lock()
	defer sql.Unlock()
	if !sql.db.CanFind(database.Table, "where steam_id = "+steamId) {
		return player{}, nil // 规避没有该用户数据的报错
	}
	err = sql.db.Find(database.Table, &dbInfo, "where steam_id = "+steamId)
	return
}

// findAll 查询所有库信息
func (sql *streamDB) findAll() (dbInfos []player, err error) {
	sql.Lock()
	defer sql.Unlock()
	// TODO 数据量分页，应该暂时不需要
	err = sql.db.FindFor(database.Table, dbInfos, "", func() error {
		return fmt.Errorf("【steamstatus插件】扫描数据库错误，错误为：%+v", err)
	})
	return
}

// del 删除指定数据
func (sql *streamDB) del(steamId string) error {
	sql.Lock()
	defer sql.Unlock()
	return sql.db.Del(database.Table, "where steam_id = "+steamId)
}

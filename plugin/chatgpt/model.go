package chatgpt

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	fcext "github.com/FloatTech/floatbox/ctxext"
	sql "github.com/FloatTech/sqlite"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type model struct {
	sql *sql.Sqlite
	sync.RWMutex
}

type mode struct {
	Mode    string `db:"mode"`
	Content string `db:"content"`
}

type note struct {
	GroupID int64  `db:"groupid"`
	Content string `db:"content"`
}

type key struct {
	QQuid   int64  `db:"qquid"`
	Content string `db:"keys"`
}

type gtoqq struct {
	GroupID int64 `db:"groupid"`
	QQuid   int64 `db:"qquid"`
}

var (
	db    = &model{sql: &sql.Sqlite{}}
	getdb = fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		db.sql.DBPath = engine.DataFolder() + "chatgpt.db"
		err := db.sql.Open(time.Hour * 24)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		err = db.sql.Create("mode", &mode{})
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		err = db.sql.Create("note", &note{})
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		err = db.sql.Create("key", &key{})
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		err = db.sql.Create("gtoqq", &gtoqq{})
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return false
		}
		return true
	})
)

func (db *model) insertmode(modename, content string) (err error) {
	db.Lock()
	defer db.Unlock()
	m := mode{
		Mode:    modename,
		Content: content,
	}
	err = db.sql.Insert("mode", &m)
	if err != nil {
		return
	}
	return
}

func (db *model) findmode(modename string) (content string, err error) {
	db.Lock()
	defer db.Unlock()
	var m mode
	err = db.sql.Find("mode", &m, "where mode = '"+modename+"'")
	if err != nil {
		return
	}
	return m.Content, nil
}

func (db *model) changemode(gid int64, modename string) (err error) {
	content, err := db.findmode(modename)
	if err != nil {
		return
	}
	n := note{
		GroupID: gid,
		Content: content,
	}
	db.Lock()
	defer db.Unlock()
	return db.sql.Insert("note", &n)
}

func (db *model) delgroupmode(gid int64) (err error) {
	db.Lock()
	defer db.Unlock()
	return db.sql.Del("note", "where groupid = "+strconv.FormatInt(gid, 10))
}

func (db *model) findgroupmode(gid int64) (content string, err error) {
	db.Lock()
	defer db.Unlock()
	var m mode
	err = db.sql.Find("note", &m, "where groupid = "+strconv.FormatInt(gid, 10))
	if err != nil {
		return
	}
	return m.Content, nil
}

func (db *model) findformode() (string, error) {
	db.Lock()
	defer db.Unlock()
	var sb strings.Builder
	sb.WriteString("当前所有预设:")
	var m mode
	err := db.sql.FindFor("mode", &m, "", func() error {
		sb.WriteString("\n")
		sb.WriteString(m.Mode)
		return nil
	})
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (db *model) insertkey(qquid int64, content string) (err error) {
	db.Lock()
	defer db.Unlock()
	m := key{
		QQuid:   qquid,
		Content: content,
	}
	return db.sql.Insert("key", &m)
}

func (db *model) delkey(gid int64) (err error) {
	db.Lock()
	defer db.Unlock()
	return db.sql.Del("key", "where qquid = "+strconv.FormatInt(gid, 10))
}

func (db *model) findkey(gid int64) (content string, err error) {
	db.Lock()
	defer db.Unlock()
	var m key
	err = db.sql.Find("key", &m, "where qquid = "+strconv.FormatInt(gid, 10))
	if err != nil {
		return "", errors.New("账号未绑定OpenAI-apikey,请私聊设置key后使用")
	}
	return m.Content, nil
}

func (db *model) insertgkey(qquid, guid int64) (err error) {
	db.Lock()
	defer db.Unlock()
	var n key
	err = db.sql.Find("key", &n, "where qquid = "+strconv.FormatInt(qquid, 10))
	if err != nil || n.Content == "" {
		return errors.New("授权账号未绑定OpenAI-apikey,请私聊设置key以后使用")
	}
	m := gtoqq{
		GroupID: guid,
		QQuid:   qquid,
	}
	return db.sql.Insert("gtoqq", &m)
}

func (db *model) delgkey(gid int64) (err error) {
	db.Lock()
	defer db.Unlock()
	return db.sql.Del("gtoqq", "where groupid = "+strconv.FormatInt(gid, 10))
}

func (db *model) findgtoqq(gid int64) (qquid int64, err error) {
	db.Lock()
	defer db.Unlock()
	var m gtoqq
	err = db.sql.Find("gtoqq", &m, "where groupid = "+strconv.FormatInt(gid, 10))
	if err != nil {
		return 0, errors.New("没有用户授权key")
	}
	return m.QQuid, nil
}

func (db *model) findgkey(gid int64) (content string, err error) {
	db.Lock()
	defer db.Unlock()
	var m gtoqq
	err = db.sql.Find("gtoqq", &m, "where groupid = "+strconv.FormatInt(gid, 10))
	if err != nil {
		return "", errors.New("未设置OpenAI-apikey,请私聊设置key以后授权本群使用")
	}
	var n key
	err = db.sql.Find("key", &n, "where qquid = "+strconv.FormatInt(m.QQuid, 10))
	if err != nil || n.Content == "" {
		return "", errors.New("授权账号未绑定OpenAI-apikey,请私聊设置key以后使用")
	}
	return n.Content, nil
}

func getkey(ctx *zero.Ctx) (key string, err error) {
	// 先从群聊中查找API Key
	if ctx.Event.GroupID != 0 {
		if key, err = db.findgkey(ctx.Event.GroupID); err == nil {
			return key, nil
		}
	}
	// 再从个人中查找API Key
	if key, err = db.findkey(-ctx.Event.UserID); err == nil {
		return key, nil
	}
	// 最后从全局中查找API Key
	if key, err = db.findgkey(-1); err == nil {
		return key, nil
	}
	// 如果都没有设置则会返回错误提示
	return "", err
}

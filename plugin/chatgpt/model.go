package chatgpt

import (
	"strconv"
	"strings"
	"sync"
	"time"

	fcext "github.com/FloatTech/floatbox/ctxext"
	sql "github.com/FloatTech/sqlite"
	ctrl "github.com/FloatTech/zbpctrl"
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
		m := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
		_ = m.Manager.Response(chatgptapikeygid)
		_ = m.Manager.GetExtra(chatgptapikeygid, &apiKey)
		if apiKey == "" {
			ctx.SendChain(message.Text("ERROR: 未设置OpenAI apikey"))
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
	err = db.sql.Insert("note", &n)
	if err != nil {
		return
	}
	return
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

package repo

import (
	"errors"
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/vars"
	"github.com/FloatTech/floatbox/ctxext"
	sql "github.com/FloatTech/sqlite"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"
)

type command struct {
	sql *sql.Sqlite
	sync.RWMutex
}

type Global struct {
	Id        int    `db:"id" json:"id"`
	Proxy     string `db:"proxy" json:"proxy"`           // 代理
	NbServ    string `db:"nb_serv" json:"nb_serv"`       // newbing 服务地址
	Bot       string `db:"bot" json:"bot"`               // AI类型
	MaxTokens int    `db:"max_tokens" json:"max_tokens"` // openai-api 最大Tokens
	Preset    string `db:"preset" json:"preset"`         // 默认预设
	DrawServ  string `db:"draw_serv" json:"draw_serv"`   // AI作画服务地址
	DrawBody  string `db:"draw_body" json:"draw_body"`   // 作画Json模版
}

type Token struct {
	Key    string `db:"key" json:"key" query:"like"`
	Type   string `db:"type" json:"type" query:"="` // 类型
	Email  string `db:"email" json:"email"`         // 邮箱
	Passwd string `db:"passwd" json:"passwd"`       // 密码
	AppId  string `db:"claude_bot" json:"app_id"`   // Claude APPID
	Token  string `db:"token" json:"token"`         // 凭证
	Expire string `db:"expire" json:"expire"`       // 过期日期
}

type PresetScene struct {
	Key     string `db:"key" json:"key" query:"like"`
	Type    string `db:"type" json:"type" query:"="` // 类型
	Content string `db:"content" json:"content"`     // 预设内容
	Message string `db:"message" json:"message"`     // 消息模版
	Chain   string `db:"chain" json:"chain"`         // 拦截处理器
}

var (
	cmd = &command{
		sql: &sql.Sqlite{},
	}

	OnceOnSuccess = ctxext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		ready, err := postRef()
		if err != nil {
			ctx.Send(err.Error())
		}
		return ready
	})
)

func init() {
	// 等待ZeroBot初始化
	go func() {
		for {
			if vars.E != nil {
				_, _ = postRef()
				return
			}
			time.Sleep(time.Second)
		}
	}()
}

func postRef() (bool, error) {
	if vars.E == nil {
		return false, errors.New("ZeroBot未初始化")
	}

	cmd.sql.DBPath = vars.E.DataFolder() + "storage.db"
	err := cmd.sql.Open(time.Hour * 24)
	if err != nil {
		return false, err
	}

	// 初始化数据表
	err = cmd.sql.Create("global", &Global{})
	if err != nil {
		return false, err
	}

	err = cmd.sql.Create("token", &Token{})
	if err != nil {
		return false, err
	}

	err = cmd.sql.Create("preset_scene", &PresetScene{})
	if err != nil {
		return false, err
	}

	return true, nil
}

// 构建查询条件
func BuildCondition(model any) string {
	var condition = ""
	v := reflect.ValueOf(model)
	for index := 0; index < v.NumField(); index++ {
		db, ok1 := v.Type().Field(index).Tag.Lookup("db")
		query, ok2 := v.Type().Field(index).Tag.Lookup("query")
		if !ok1 || !ok2 {
			continue
		}
		s := v.Field(index).String()
		if s == "" {
			continue
		}
		if query == "like" {
			condition += db + " " + query + " '%" + s + "%' and "
		} else {
			condition += db + " " + query + " '" + s + "' and "
		}
	}

	if condition != "" {
		cut, ok := strings.CutSuffix(condition, " and ")
		if ok {
			condition = "where " + cut
		} else {
			condition = "where " + condition
		}
	}
	return condition
}

func (c *command) Count(table string, condition string) (num int, err error) {
	if c.sql.DB == nil {
		return 0, sql.ErrNilDB
	}
	stmt, err := cmd.sql.DB.Prepare("SELECT COUNT(1) FROM " + wraptable(table) + condition + ";")
	if err != nil {
		return 0, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return 0, err
	}
	if rows.Err() != nil {
		return 0, rows.Err()
	}
	if rows.Next() {
		err = rows.Scan(&num)
	}
	err = rows.Close()
	if err != nil {
		return 0, err
	}
	return num, err
}

func wraptable(table string) string {
	first := []rune(table)[0]
	if first < unicode.MaxLatin1 && unicode.IsDigit(first) {
		return "[" + table + "]"
	} else {
		return "'" + table + "'"
	}
}

func GetGlobal() Global {
	var g Global
	if err := cmd.sql.Find("global", &g, ""); err != nil {
		g = Global{
			Id:  1,
			Bot: "openai-api",
		}
	}
	return g
}

func InsertGlobal(g Global) error {
	cmd.Lock()
	defer cmd.Unlock()
	return cmd.sql.Insert("global", &g)
}

func SetProxy(p string) error {
	cmd.Lock()
	defer cmd.Unlock()
	global := GetGlobal()
	global.Proxy = p
	return cmd.sql.Insert("global", &global)
}

func InsertToken(token Token) error {
	cmd.Lock()
	defer cmd.Unlock()
	var t Token
	err := cmd.sql.Find("token", &t, "where type='"+token.Type+"' and key='"+token.Key+"'")
	if err != nil {
		return cmd.sql.Insert("token", &token)
	} else {
		return errors.New("`" + token.Key + "`已存在")
	}
}

func UpdateToken(t Token) {
	cmd.Lock()
	defer cmd.Unlock()
	if err := cmd.sql.Insert("token", &t); err != nil {
		logrus.Warn(err)
	}
}

func GetToken(key string) *Token {
	var t Token
	err := cmd.sql.Find("token", &t, "where key='"+key+"'")
	if err != nil {
		return nil
	}
	return &t
}

func FindTokens(t string) ([]*Token, error) {
	if t != "" {
		return sql.FindAll[Token](cmd.sql, "token", "where type='"+t+"'")
	} else {
		return sql.FindAll[Token](cmd.sql, "token", "")
	}
}

func RemoveToken(key string) {
	cmd.sql.Del("token", "where key='"+key+"'")
}

func GetPresetScene(key string) *PresetScene {
	var p PresetScene
	err := cmd.sql.Find("preset_scene", &p, "where key='"+key+"'")
	if err != nil {
		return nil
	}
	return &p
}

func FindPresetScenes(t string) ([]*PresetScene, error) {
	if t != "" {
		return sql.FindAll[PresetScene](cmd.sql, "preset_scene", "where type='"+t+"'")
	} else {
		return sql.FindAll[PresetScene](cmd.sql, "preset_scene", "")
	}
}

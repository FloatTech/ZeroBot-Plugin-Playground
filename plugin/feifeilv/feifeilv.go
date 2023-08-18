// Package feifeilv 可能媒用
package feifeilv

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	version     = "10000"
	sample      = 1
	act         = "getMsg"
	textType    = "text"
	feifeiURL   = "https://m.feifeilv.top"
	configURL   = feifeiURL + "/api/config/get"
	libraryURL  = feifeiURL + "/api/library/get"
	msgURL      = feifeiURL + "/api/msg/get"
	randomIndex = "/pagesImp/random/index"
)

// CommonRequest 常规请求
type CommonRequest struct {
	Version string `json:"version"`
	Time    string `json:"time"`
}

// ConfigResponse 常用响应
type ConfigResponse struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data ConfigData `json:"data"`
}

// ToolList 工具箱
type ToolList struct {
	CreateTime int64  `json:"createTime"`
	State      int    `json:"state"`
	Hot        int    `json:"hot"`
	ID         string `json:"_id"`
	Name       string `json:"name"`
	Href       string `json:"href"`
	Keywords   string `json:"keywords"`
	Type       string `json:"type"`
	Headimg    string `json:"headimg"`
	Theme      string `json:"theme,omitempty"`
}

// ConfigData 配置相关数据
type ConfigData struct {
	Version  string     `json:"version"`
	Notice   string     `json:"notice"`
	ToolList []ToolList `json:"toolList"`
}

// LibraryRequest 语录请求入参
type LibraryRequest struct {
	SupID   string `json:"supId"`
	Type    string `json:"type"`
	Sample  int    `json:"sample"`
	Version string `json:"version"`
	Time    string `json:"time"`
}

// LibraryResponse 随机返回
type LibraryResponse struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data LibraryTotalData `json:"data"`
}

// LibraryData 随机数据
type LibraryData struct {
	ID         string `json:"_id"`
	SupID      string `json:"supId"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	Source     string `json:"source"`
	CreateTime int64  `json:"createTime"`
	State      int    `json:"state"`
}

// LibraryTotalData 全体数据
type LibraryTotalData struct {
	Data  []LibraryData `json:"data"`
	Total int           `json:"total"`
}

// MsgRequest 消息请求
type MsgRequest struct {
	Keywords  string `json:"keywords"`
	Act       string `json:"act"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
	Version   string `json:"version"`
	Time      string `json:"time"`
}

// MsgResponse 消息出参
type MsgResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data MsgTotalData `json:"data"`
}

// MsgData 消息
type MsgData struct {
	ID         string `json:"_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Type       string `json:"type,omitempty"`
	CreateTime int64  `json:"createTime,omitempty"`
	Hot        int    `json:"hot"`
	State      int    `json:"state"`
}

// MsgTotalData 消息统计
type MsgTotalData struct {
	Data  []MsgData `json:"data"`
	Total int       `json:"total"`
}

func init() {
	engine := control.Register("feifeilv", &ctrl.Options[*zero.Ctx]{
		PrivateDataFolder: "feifeilv",
		DisableOnDefault:  false,
		Brief:             "可能媒用",
		Help:              "- 可能媒用\n- 聊天话术 开场白 0",
	})
	cachePath := engine.DataFolder() + "cache/"
	_ = os.RemoveAll(cachePath)
	_ = os.MkdirAll(cachePath, 0755)
	engine.OnFullMatchGroup([]string{"可能媒用"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		next := zero.NewFutureEvent("message", 999, false, ctx.CheckSession())
		recv, cancel := next.Repeat()
		defer cancel()
		tools, err := getAllTool()
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		tools = filterTool(tools)
		tex := "请输入可能媒用序号\n"
		for i, v := range tools {
			tex += fmt.Sprintf("%d. %s\n", i, v.Name)
		}
		base64Str, err := text.RenderToBase64(tex, text.FontFile, 400, 20)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(message.Image("base64://" + binary.BytesToString(base64Str)))
		for {
			select {
			case <-time.After(time.Second * 10):
				num := rand.Intn(len(tools))
				ctx.SendChain(message.Text("时间太久啦！", zero.BotConfig.NickName[0], "帮你选择查询", tools[num].Name))
				library, err := getLibraryData(tools[num].ID)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				var libraryTex string
				if library.Source != "" {
					libraryTex = library.Content + "\n——" + library.Source
				} else {
					libraryTex = library.Content
				}
				ctx.SendChain(message.Text(libraryTex))
				return
			case c := <-recv:
				msg := c.Event.Message.ExtractPlainText()
				num, err := strconv.Atoi(msg)
				if err != nil {
					ctx.SendChain(message.Text("请输入数字!"))
					continue
				}
				if num < 0 || num >= len(tools) {
					ctx.SendChain(message.Text("序号非法!"))
					continue
				}
				library, err := getLibraryData(tools[num].ID)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				var libraryTex string
				if library.Source != "" {
					libraryTex = library.Content + "\n——" + library.Source
				} else {
					libraryTex = library.Content
				}
				ctx.SendChain(message.Text(libraryTex))
				return
			}
		}
	})
	engine.OnRegex(`^聊天话术\s?(\S{1,25})\s?(\d*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		keyword := ctx.State["regex_matched"].([]string)[1]
		pagenum := ctx.State["regex_matched"].([]string)[2]
		pageIndex, _ := strconv.Atoi(pagenum)
		msg, err := getMsg(keyword, pageIndex, 10)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		m := make(message.Message, 0, 10)
		tip := "关键词%v相关的对话一共有%v页, %v条, 目前是第%v页"
		tagMap := make(map[string]struct{}, 0)
		m = append(m, ctxext.FakeSenderForwardNode(ctx, message.Text(fmt.Sprintf(tip, keyword, (msg.Data.Total-1)/10, msg.Data.Total, pageIndex))))
		for _, v := range msg.Data.Data {
			m = append(m, ctxext.FakeSenderForwardNode(ctx, message.Text(v.Content)))
			for _, t := range strings.Split(v.Type, ",") {
				tagMap[t] = struct{}{}
			}
		}
		tagList := make([]string, 0, len(tagMap))
		for k := range tagMap {
			if k != "" {
				tagList = append(tagList, k)
			}
		}
		m = append(m, ctxext.FakeSenderForwardNode(ctx, message.Text("相关标签: ", strings.Join(tagList, ","))))
		ctx.Send(m)
	})
}

func getAllTool() (tools []ToolList, err error) {
	r := CommonRequest{
		Version: version,
		Time:    getTimestamp(),
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return
	}
	data, err := web.PostData(configURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return
	}
	var rsp ConfigResponse
	err = json.Unmarshal(data, &rsp)
	if err != nil {
		return
	}
	tools = rsp.Data.ToolList
	return
}

func filterTool(tools []ToolList) (newTools []ToolList) {
	newTools = make([]ToolList, 0, len(tools))
	for _, v := range tools {
		if v.Href == randomIndex {
			newTools = append(newTools, v)
		}
	}
	return
}

func getLibraryData(supID string) (libray LibraryData, err error) {
	r := LibraryRequest{
		SupID:   supID,
		Type:    textType,
		Sample:  sample,
		Version: version,
		Time:    getTimestamp(),
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return
	}
	data, err := web.PostData(libraryURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return
	}
	var rsp LibraryResponse
	err = json.Unmarshal(data, &rsp)
	if err != nil {
		return
	}
	if len(rsp.Data.Data) == 0 {
		err = errors.New("请求数据为空")
		return
	}
	libray = rsp.Data.Data[0]
	return
}

func getTimestamp() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	hash := md5.Sum([]byte(timestamp))
	return fmt.Sprintf("%x", hash)
}

func getMsg(keyword string, pageIndex, pageSize int) (msg MsgResponse, err error) {
	r := MsgRequest{
		Keywords:  keyword,
		Act:       act,
		Version:   version,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Time:      getTimestamp(),
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return
	}
	data, err := web.PostData(msgURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return
	}
	return
}

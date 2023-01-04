// Package vote 实时投票
package vote

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/golang/freetype"
	"github.com/wcharczuk/go-chart/v2"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("vote", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "投票",
		Help:              "- 投票选择可口可乐还是百事可乐\n",
		PrivateDataFolder: "vote",
	})

	engine.OnPrefix("投票选择", zero.OnlyGroup, getPara).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			rawOptions := strings.Split(ctx.State["args"].(string), "还是")
			options := make([]string, 0, 16)
			for count, option := range rawOptions {
				options = append(options, strconv.Itoa(count)+". "+option)
			}
			if len(options) == 0 || len(options) == 1 {
				ctx.SendChain(message.Text("投票的选项太少, 退出投票"))
				return
			}
			ctx.SendChain(message.Text("投票开始, 设定的投票时间到或者60秒内无人投票就停止投票\n请按序号投票:\n", strings.Join(options, "\n")))
			voteMap := make(map[int]int, 0)
			repeatMap := make(map[int64]int, 0)
			mode := ctx.State["vote_paras"].([2]int)[0]
			voteDuration := ctx.State["vote_paras"].([2]int)[1]
			next := zero.NewFutureEvent("message", 999, false, zero.RegexRule(`^\d{1,2}$`),
				zero.OnlyGroup, zero.CheckGroup(ctx.Event.GroupID))
			recv, cancel := next.Repeat()
			defer cancel()
			after := time.NewTimer(time.Duration(voteDuration) * time.Second)
		EXIT:
			for {
				select {
				case <-time.After(60 * time.Second):
					ctx.SendChain(message.Text("投票结束"))
					break EXIT
				case <-after.C:
					ctx.SendChain(message.Text("投票结束"))
					break EXIT
				case c := <-recv:
					uid := c.Event.UserID
					choose, err := strconv.Atoi(c.Event.Message.String())
					if err != nil {
						ctx.SendChain(message.Text("ERROR:", err))
					}
					if _, ok := repeatMap[uid]; !ok && choose >= 0 && choose < len(options) {
						voteMap[choose]++
					}
					if mode == 2 {
						repeatMap[uid] = 1
					}
				}
			}
			vc := rankByVote(voteMap)
			if len(vc) > 20 {
				vc = vc[:20]
			}
			// 绘图
			b, err := os.ReadFile(text.FontFile)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			font, err := freetype.ParseFont(b)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			if len(vc) == 0 {
				ctx.SendChain(message.Text("无人投票, 退出投票"))
				return
			}
			bars := make([]chart.Value, len(vc))
			result := ""
			for i, v := range vc {
				bars[i] = chart.Value{
					Value: float64(v.Value),
					Label: fmt.Sprintf("%d. %s", v.Key, rawOptions[v.Key]),
				}
				result += fmt.Sprintf("%d. %s: %d票", v.Key, rawOptions[v.Key], v.Value)
			}
			graph := chart.BarChart{
				Font:  font,
				Title: "投票结果",
				Background: chart.Style{
					Padding: chart.Box{
						Top: 40,
					},
				},
				Height:   500,
				BarWidth: 25,
				Bars:     bars,
			}
			drawedFile := engine.DataFolder() + "vote.png"
			f, err := os.Create(drawedFile)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			err = graph.Render(chart.PNG, f)
			_ = f.Close()
			if err != nil {
				_ = os.Remove(drawedFile)
				ctx.SendChain(message.Text(result))
				return
			}
			ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
		})
}

func getPara(ctx *zero.Ctx) bool {
	next := zero.NewFutureEvent("message", 999, false, ctx.CheckSession())
	recv, cancel := next.Repeat()
	i := 0
	paras := [2]int{}
	ctx.SendChain(message.Text("请输入投票模式序号\n1. 可重复投票\n2. 不可重复投票"))
	for {
		select {
		case <-time.After(time.Second * 120):
			ctx.SendChain(message.Text("未输入,退出投票"))
			cancel()
			return false
		case c := <-recv:
			msg := c.Event.Message.ExtractPlainText()
			num, err := strconv.Atoi(msg)
			if err != nil {
				ctx.SendChain(message.Text("请输入数字!"))
				continue
			}
			switch i {
			case 0:
				if num <= 0 || num > 2 {
					ctx.SendChain(message.Text("投票模式非法!"))
					continue
				}
				paras[0] = num
				ctx.SendChain(message.Text("请输入投票时间(单位秒,至少3秒)"))
			case 1:
				if num < 3 {
					ctx.SendChain(message.Text("投票时间非法!"))
					continue
				}
				cancel()
				paras[1] = num
				ctx.State["vote_paras"] = paras
				return true
			}
			i++
		}
	}
}

func rankByVote(voteMap map[int]int) pairlist {
	pl := make(pairlist, len(voteMap))
	i := 0
	for k, v := range voteMap {
		pl[i] = pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type pair struct {
	Key   int
	Value int
}

type pairlist []pair

func (p pairlist) Len() int           { return len(p) }
func (p pairlist) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairlist) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

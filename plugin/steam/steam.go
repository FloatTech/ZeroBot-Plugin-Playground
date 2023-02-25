// Package steamstatus 获取steam用户状态
package steam

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/binary"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	engine = control.Register("steam", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "steam",
		Help: "一款用来查看你steam好友状态的插件\n-----------------------\n" +
			"- 创建steam用户绑定 xxxxxxx （可输入需要绑定的 steamid )\n" +
			"- 删除steam用户绑定 xxxxxxx （删除你创建的对于 steamid 的绑定\n" +
			"- 查询steam用户绑定（查询本群内所有的绑定对象）\n-----------------------\n" +
			"TIP：steamID在用户资料页的链接上面，形如7656119820673xxxx，需要管理员先开启查询",
		PrivateDataFolder: "steam",
	})
)

func init() {
	// 创建绑定流程
	engine.OnRegex(`^创建steam用户绑定\s*(.*)$`, zero.OnlyGroup, getDB).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		steamID := ctx.State["regex_matched"].([]string)[1]
		// 获取用户状态
		playerStatus, err := getPlayerStatus([]string{steamID})
		if err != nil {
			ctx.SendChain(message.Text("创建失败，获取用户信息错误，错误：" + err.Error()))
			return
		}
		if len(playerStatus) == 0 {
			ctx.SendChain(message.Text("需要绑定的用户不存在，请检查id或url"))
			return
		}
		playerData := playerStatus[0]
		// 判断用户是否已经初始化：若未初始化，通过用户的steamID获取当前状态并初始化；若已经初始化则更新用户信息
		info, err := database.find(steamID)
		if err != nil {
			ctx.SendChain(message.Text("创建失败，数据库错误，错误：" + err.Error()))
			return
		}
		// 处理数据
		groupID := strconv.FormatInt(ctx.Event.GroupID, 10)
		if info.Target == "" {
			info = player{
				SteamID:       steamID,
				PersonaName:   playerData.PersonaName,
				Target:        groupID,
				GameID:        playerData.GameID,
				GameExtraInfo: playerData.GameExtraInfo,
				LastUpdate:    time.Now().Unix(),
			}
		} else if !strings.Contains(info.Target, groupID) {
			info.Target = strings.Join([]string{info.Target, groupID}, ",")
		}
		// 更新数据库
		if err = database.update(info); err != nil {
			ctx.SendChain(message.Text("更新绑定失败，错误：" + err.Error()))
			return
		}
		ctx.SendChain(message.Text("设置成功"))
	})
	// 删除绑定流程
	engine.OnRegex(`^删除steam用户绑定\s*(.*)$`, zero.OnlyGroup, getDB).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		steamID := ctx.State["regex_matched"].([]string)[1]
		groupID := strconv.FormatInt(ctx.Event.GroupID, 10)
		// 判断是否已经绑定该steamID，若已绑定就将群列表从推送群列表钟去除
		info, err := database.findWithGroupID(steamID, groupID)
		if err != nil {
			ctx.SendChain(message.Text("删除失败，数据库错误，错误：" + err.Error()))
			return
		}
		if info.SteamID == "" {
			ctx.SendChain(message.Text("所需要删除的绑定不存在，请检查绑定关系。"))
			return
		}
		// 从绑定列表中剔除需要删除的对象
		targets := strings.Split(info.Target, ",")
		newTargets := make([]string, 0)
		for _, target := range targets {
			if target == groupID {
				continue
			}
			newTargets = append(newTargets, target)
		}
		if len(newTargets) == 0 {
			if err = database.del(steamID); err != nil {
				ctx.SendChain(message.Text("删除数据失败，请检查数据库，错误：" + err.Error()))
				return
			}
		} else {
			info.Target = strings.Join(newTargets, ",")
			if err = database.update(info); err != nil {
				ctx.SendChain(message.Text("删除群推送失败，请检查数据库，错误：" + err.Error()))
				return
			}
		}
		ctx.SendChain(message.Text("设置成功"))
	})
	// 查询当前群绑定信息
	engine.OnFullMatch("查询steam用户绑定", zero.OnlyGroup, getDB).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// 获取群信息
		groupID := strconv.FormatInt(ctx.Event.GroupID, 10)
		// 获取所有绑定信息
		infos, err := database.findAll()
		if err != nil {
			ctx.SendChain(message.Text("查询绑定失败，检查数据库，错误：" + err.Error()))
			return
		}
		if len(infos) == 0 {
			ctx.SendChain(message.Text("还未建立过用户绑定关系！"))
			return
		}
		// 遍历所有信息，如果包含该群就收集对应的steamID
		players := make([]string, 0)
		for _, info := range infos {
			if strings.Contains(info.Target, groupID) {
				players = append(players, " "+info.PersonaName+":"+info.SteamID)
			}
		}
		if len(players) == 0 {
			ctx.SendChain(message.Text("查询成功，该群暂时还没有被绑定的用户！"))
			return
		}
		// 组装并返回结果
		logText := fmt.Sprintf(" 查询steam用户绑定成功，该群绑定的用户有: \n%+v \n", strings.Join(players, "\n"))
		data, err := text.RenderToBase64(logText, text.FontFile, 400, 18)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		ctx.SendChain(message.Image("base64://" + binary.BytesToString(data)))
	})
}

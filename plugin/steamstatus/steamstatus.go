// Package steamstatus 获取steam用户状态
package steamstatus

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	engine = control.Register("steamstatus", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "steam视奸",
		Help: "一款用来视奸你steam好友的插件\n-----------------------\n" +
			"- 创建监听 xxxxxxx （可输入需要监听的 steamid )\n" +
			"- 删除监听 xxxxxxx （删除你创建的对于 steamid 的监听）\n-----------------------\n" +
			"TIP：steamID在用户资料页的链接上面，形如7656119820673xxxx",
		PrivateDataFolder: "steamstatus",
	})
)

func init() {
	// 加载密钥信息
	apiKeyFile := engine.DataFolder() + "apikey.txt"
	if file.IsExist(apiKeyFile) {
		apiKeyByte, err := os.ReadFile(apiKeyFile)
		if err != nil {
			panic(err)
		}
		apiKey = strings.TrimSpace(string(apiKeyByte))
	} else { // 如果没有配置文件直接退出
		control.Delete("steamstatus")
		logrus.Info("[steamstatus] 未配置对应的插件链接密钥，主动移除插件")
		return
	}
	// 初始化数据库
	err := initStore()
	if err != nil {
		// 抛错误但是不影响程序整体运行
		logrus.Errorf("[steamstatus] 初始化数据库失败，请检查data文件目录，错误为：%+v", err)
		panic(err)
	}
	// 创建监听流程
	engine.OnRegex(`^创建监听\s*(.*)$`, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		vanityUrl := ctx.State["regex_matched"].([]string)[1]
		// 获取参数，判断是url还是steamId
		steamId := vanityUrl
		//steamId, err = getPlayerSteamIdWithUrl(vanityUrl)
		//if err != nil {
		//	// 这里不直接抛出错误
		//	logrus.Error("[steamstatus] 创建失败，转换steamId接口异常，错误：" + err.Error())
		//	steamId = vanityUrl
		//}
		// 通过steamId来获取用户当前状态
		playerStatus, err := getPlayerStatus([]string{steamId})
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
		info, err := database.find(steamId)
		if err != nil {
			ctx.SendChain(message.Text("创建失败，数据库错误，错误：" + err.Error()))
			return
		}
		// 处理数据
		groupId := strconv.FormatInt(ctx.Event.GroupID, 10)
		if info.Target == "" {
			info = player{
				SteamId:       steamId,
				PersonaName:   playerData.PersonaName,
				Target:        groupId,
				GameId:        playerData.GameId,
				GameExtraInfo: playerData.GameExtraInfo,
				LastUpdate:    time.Now().Unix(),
			}
		} else {
			if !strings.Contains(info.Target, groupId) {
				info.Target = strings.Join([]string{info.Target, groupId}, ",")
			}
		}
		// 更新数据库
		if err := database.update(info); err != nil {
			ctx.SendChain(message.Text("更新监听绑定失败，错误：" + err.Error()))
			return
		}
		ctx.SendChain(message.Text("设置成功"))
	})
	// 删除监听流程
	engine.OnRegex(`^删除监听\s*(.*)$`, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		steamId := ctx.State["regex_matched"].([]string)[1]
		// 判断是否已经监听该steamId，若已监听就将群列表从推送群列表钟去除
		info, err := database.find(steamId)
		if err != nil {
			ctx.SendChain(message.Text("删除失败，数据库错误，错误：" + err.Error()))
			return
		}
		if info.SteamId == "" {
			ctx.SendChain(message.Text("所需要删除的监听不存在，请检查绑定关系。"))
			return
		}
		groupId := strconv.FormatInt(ctx.Event.GroupID, 10)
		if !strings.Contains(info.Target, groupId) {
			ctx.SendChain(message.Text("所需要删除的监听未在当前群进行绑定，请检查绑定关系。"))
			return
		}
		targets := strings.Split(info.Target, ",")
		newTargets := make([]string, 0)
		for _, target := range targets {
			if target == groupId {
				continue
			}
			newTargets = append(newTargets, target)
		}
		if len(newTargets) == 0 {
			if err := database.del(steamId); err != nil {
				ctx.SendChain(message.Text("删除数据失败，请检查数据库，错误：" + err.Error()))
				return
			}
		} else {
			info.Target = strings.Join(newTargets, ",")
			if err := database.update(info); err != nil {
				ctx.SendChain(message.Text("删除群推送失败，请检查数据库，错误：" + err.Error()))
				return
			}
		}
		ctx.SendChain(message.Text("设置成功"))
	})
	// 查询当前群监听信息
	engine.OnFullMatch("查询监听", zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// 获取群信息
		groupId := strconv.FormatInt(ctx.Event.GroupID, 10)
		// 获取所有监听信息
		infos, err := database.findAll()
		if err != nil {
			ctx.SendChain(message.Text("查询监听失败，检查数据库，错误：" + err.Error()))
			return
		}
		// 遍历所有信息，如果包含该群就收集对应的steamId
		players := make([]string, 0)
		for _, info := range infos {
			if strings.Contains(info.Target, groupId) {
				players = append(players, info.PersonaName)
			}
		}
		if len(players) == 0 {
			ctx.SendChain(message.Text("查询成功，该群暂时还没有被监听的用户！"))
		}
		// 组装并返回结果
		result := strings.Join(players, ",\n")
		ctx.SendChain(message.Text("查询成功，该群监听的用户有：\n" + result))
	})
}

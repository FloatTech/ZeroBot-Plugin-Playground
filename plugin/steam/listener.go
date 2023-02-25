package steam

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	engineListener = control.Register("steamlistener", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "steamlistener",
		Help: "管理员绑定监听面板\n-----------------------\n" +
			"- 初始化steam链接密钥 xxxxxxx （密钥在steam网站申请，申请地址：https://steamcommunity.com/dev/registerkey）\n" +
			"- 当前steam链接密钥（查询已经绑定的密钥）\n" +
			"- 拉取steam绑定用户状态（使用插件定时任务开始）\n-----------------------\n" +
			"TIP：需要先进行密钥初始化，绑定用户之后使用定时任务开启查询即可自动播报",
		PrivateDataFolder: "steamlistener",
	})
)

func init() {
	// 从本地获取steamKey
	steamKeyFile := engineListener.DataFolder() + "apikey.txt"
	if file.IsExist(steamKeyFile) {
		apikey, err := os.ReadFile(steamKeyFile)
		if err != nil {
			panic(err)
		}
		apiKey = binary.BytesToString(apikey)
	}
	engineListener.OnRegex(`初始化steam链接密钥\s*(.*)$`, zero.SuperUserPermission, getDB).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// 直接赋值给持久化字短
		apiKey = ctx.State["regex_matched"].([]string)[1]
		// 持久化到本地文件
		if err := os.WriteFile(steamKeyFile, binary.StringToBytes(apiKey), 0777); err != nil {
			ctx.SendChain(message.Text("持久化密钥失败"))
			return
		}
		ctx.SendChain(message.Text("设置链接密钥成功！"))
	})
	engineListener.OnFullMatch(`当前steam链接密钥`, zero.SuperUserPermission, getDB).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if apiKey == "" {
			ctx.SendChain(message.Text("当前尚未配置链接密钥，请先配置对应密钥！"))
			return
		}
		ctx.SendChain(message.Text("链接密钥为：" + apiKey))
	})
	engineListener.OnFullMatch("拉取steam绑定用户状态", getDB).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		listenUserChange(ctx)
	})
}

// listenUserChange 用于监听用户的信息变化
func listenUserChange(ctx *zero.Ctx) {
	// 获取所有处于监听状态的用户信息
	infos, err := database.findAll()
	if err != nil {
		// 挂了就给管理员发消息
		ctx.SendPrivateMessage(zero.BotConfig.SuperUsers[0], message.Text("[steamstatus] 出问题了，快喵一眼。报错："+err.Error()))
		return
	}
	if len(infos) == 0 {
		return
	}
	// 收集这波用户的streamId，然后查当前的状态，并建立信息映射表
	streamIds := make([]string, 0)
	localPlayerMap := make(map[string]player)
	for _, info := range infos {
		streamIds = append(streamIds, info.SteamID)
		localPlayerMap[info.SteamID] = info
	}
	// 将所有用户状态查一遍
	playerStatus, err := getPlayerStatus(streamIds)
	if err != nil {
		// 出错就发消息
		ctx.SendPrivateMessage(zero.BotConfig.SuperUsers[0], message.Text("[steamstatus] 出问题了，快喵一眼。报错："+err.Error()))
		return
	}
	// 遍历返回的信息做对比，假如信息有变化则发消息
	now := time.Now()
	for _, playerInfo := range playerStatus {
		localInfo := localPlayerMap[playerInfo.SteamID]
		// 排除不需要处理的情况
		if localInfo.GameID == "" && playerInfo.GameID == "" {
			continue
		}
		// 打开游戏
		if localInfo.GameID == "" && playerInfo.GameID != "" {
			sendGroupMessageForPlayerGroups(ctx, localInfo, fmt.Sprintf("%+v正在玩%+v",
				playerInfo.PersonaName, playerInfo.GameExtraInfo))
			localInfo.LastUpdate = now.Unix()
		}
		// 更换游戏
		if localInfo.GameID != "" && playerInfo.GameID != localInfo.GameID && playerInfo.GameID != "" {
			sendGroupMessageForPlayerGroups(ctx, localInfo, fmt.Sprintf("%+v玩了%+v分钟后，丢下了%+v，转头去玩%+v",
				playerInfo.PersonaName, (now.Unix()-localInfo.LastUpdate)/60, localInfo.GameExtraInfo, playerInfo.GameExtraInfo))
			localInfo.LastUpdate = now.Unix()
		}
		// 关闭游戏
		if playerInfo.GameID != localInfo.GameID && playerInfo.GameID == "" {
			sendGroupMessageForPlayerGroups(ctx, localInfo, fmt.Sprintf("%+v玩了%+v分钟后，关掉了%+v",
				playerInfo.PersonaName, (now.Unix()-localInfo.LastUpdate)/60, localInfo.GameExtraInfo))
			localInfo.LastUpdate = 0
		}
		// 更新数据
		localInfo.GameID = playerInfo.GameID
		localInfo.GameExtraInfo = playerInfo.GameExtraInfo
		if err = database.update(localInfo); err != nil {
			logrus.Errorf("[steamstatus] 更新数据条目异常，异常对象:[%+v]，错误信息：[%+v]", localInfo, err)
		}
	}
}

// sendGroupMessageForPlayerGroups 为绑定的用户所绑定群发消息
func sendGroupMessageForPlayerGroups(ctx *zero.Ctx, playerInfo player, msg string) {
	groups := strings.Split(playerInfo.Target, ",")
	for _, groupString := range groups {
		group, err := strconv.ParseInt(groupString, 10, 64)
		if err != nil {
			logrus.Errorf("[steamstatus] 出现异常，异常对象:[%+v]，错误信息：[%+v]", playerInfo, err)
			continue
		}
		ctx.SendGroupMessage(group, message.Text(msg))
	}
}

// ----------------------- 远程调用 ----------------------
const (
	URL       = "https://api.steampowered.com/"                          // steam API 调用地址
	StatusURL = "ISteamUser/GetPlayerSummaries/v2/?key=%+v&steamids=%+v" // 根据用户steamID获取用户状态
)

var apiKey string

// getPlayerStatus 获取用户状态
func getPlayerStatus(streamIds []string) ([]player, error) {
	players := make([]player, 0)
	// 校验密钥是否初始化
	if apiKey == "" {
		return players, errors.New("未进行链接密钥初始化")
	}
	// 拼接请求地址
	url := fmt.Sprintf(URL+StatusURL, apiKey, strings.Join(streamIds, ","))
	logrus.Debugf("[steamstatus] getPlayerStatus url：%+v", url)
	// 拉取并解析数据
	data, err := web.GetData(url)
	if err != nil {
		return players, err
	}
	logrus.Debugf("[steamstatus] getPlayerStatus data：%+v \n", string(data))
	index := gjson.Get(string(data), "response.players.#").Uint()
	for i := uint64(0); i < index; i++ {
		players = append(players, player{
			SteamID:       gjson.Get(string(data), fmt.Sprintf("response.players.%d.steamid", i)).String(),
			PersonaName:   gjson.Get(string(data), fmt.Sprintf("response.players.%d.personaname", i)).String(),
			GameID:        gjson.Get(string(data), fmt.Sprintf("response.players.%d.gameid", i)).String(),
			GameExtraInfo: gjson.Get(string(data), fmt.Sprintf("response.players.%d.gameextrainfo", i)).String(),
		})
	}
	logrus.Debugf("[steamstatus] getPlayerStatus players：%+v \n", players)
	return players, nil
}

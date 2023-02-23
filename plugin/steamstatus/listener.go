package steamstatus

import (
	"encoding/json"
	"fmt"
	"github.com/FloatTech/floatbox/web"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
	"strings"
	"time"
)

// startListening 开启循环监听器
func startListening() {
	// 开启一个循环来不停的扫描用户状态
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		// 监听用户变化
		listenUserChange(ctx)
		// 每完成一次循环，等个30秒，然后再继续
		time.Sleep(30 * time.Second)
		return true
	})
}

// listenUserChange 用于监听用户的信息变化
func listenUserChange(ctx *zero.Ctx) {
	// 获取所有处于监听状态的用户信息
	infos, err := database.findAll()
	if err != nil {
		// 挂了就给管理员发消息
		notice(ctx, err)
		return
	}
	// 收集这波用户的streamId，然后查当前的状态，并建立信息映射表
	streamIds := make([]string, 0)
	localPlayerMap := make(map[string]player)
	for _, info := range infos {
		streamIds = append(streamIds, info.SteamId)
		localPlayerMap[info.SteamId] = info
	}
	// 将所有用户状态查一遍
	playerStatus, err := getPlayerStatus(streamIds)
	if err != nil {
		// 出错就发消息
		notice(ctx, err)
		return
	}
	// 遍历返回的信息做对比，假如信息有变化则发消息
	now := time.Now()
	for _, playerInfo := range playerStatus {
		localInfo := localPlayerMap[playerInfo.SteamId]
		if localInfo.GameId == "" && playerInfo.GameId != "" { // 打开游戏
			sendGroupMessageForPlayerGroups(ctx, localInfo, fmt.Sprintf("%+v正在玩%+v",
				playerInfo.PersonaName, playerInfo.GameExtraInfo))
		} else if playerInfo.GameId != localInfo.GameId && playerInfo.GameId != "" { // 更换游戏
			localInfo.LastUpdate = now.Unix()
			sendGroupMessageForPlayerGroups(ctx, localInfo, fmt.Sprintf("%+v玩了%+v分钟后，丢下了%+v，转头去玩%+v",
				playerInfo.PersonaName, (now.Unix()-localInfo.LastUpdate)/60, localInfo.GameExtraInfo, playerInfo.GameExtraInfo))
		} else if playerInfo.GameId != localInfo.GameId && playerInfo.GameId == "" { // 关闭游戏
			sendGroupMessageForPlayerGroups(ctx, localInfo, fmt.Sprintf("%+v玩了%+v分钟后，关掉了%+v",
				playerInfo.PersonaName, (now.Unix()-localInfo.LastUpdate)/60, localInfo.GameExtraInfo))
		} else { // 其它情况不更新数据
			continue
		}
		// 更新数据
		localInfo.GameId = playerInfo.GameId
		localInfo.GameExtraInfo = playerInfo.GameExtraInfo
		if database.update(localInfo) != nil {
			logrus.Errorf("【steamstatus插件】更新数据条目异常，异常对象:[%+v]，错误信息：[%+v]", localInfo, err)
		}
	}
}

// notice 告警
func notice(ctx *zero.Ctx, err error) {
	for _, id := range zero.BotConfig.SuperUsers {
		ctx.SendPrivateMessage(id, message.Text("【steamstatus插件】喵的插件数据库链接炸了，快喵一眼。报错："+err.Error()))
	}
}

// sendGroupMessageForPlayerGroups 为绑定监听的用户所绑定群发消息
func sendGroupMessageForPlayerGroups(ctx *zero.Ctx, playerInfo player, msg string) {
	groups := strings.Split(playerInfo.Target, ",")
	for _, groupString := range groups {
		group, err := strconv.ParseInt(groupString, 64, 10)
		if err != nil {
			logrus.Errorf("【steamstatus插件】数据条目异常，异常对象:[%+v]，错误信息：[%+v]", playerInfo, err)
			continue
		}
		ctx.SendGroupMessage(group, message.Text(msg))
	}
}

// ----------------------- 远程调用 ----------------------
const (
	APIUrl    = "https://api.steampowered.com/"                         // steam API 调用地址
	VanityUrl = "ISteamUser/ResolveVanityURL/v1/?key=%+v&vanityurl=%+v" // 根据用户链接获取用户steamID接口，暂时不用
	StatusUrl = "ISteamUser/GetPlayerSummaries/v2/?key=%+v&steamids%+v" // 根据用户steamID获取用户状态
)

var apiKey string

// getPlayerStatus 获取用户状态
func getPlayerStatus(streamIds []string) ([]player, error) {
	players := make([]player, 0)
	url := fmt.Sprintf(APIUrl+StatusUrl, apiKey, strings.Join(streamIds, ","))
	data, err := web.GetData(url)
	if err != nil {
		return players, err
	}
	if json.Unmarshal(data, &players) != nil {
		return players, err
	}
	return players, nil
}

// 根据用户的链接获取用户的steamId
func getPlayerSteamIdWithUrl(vanityUrl string) (string, error) {
	url := fmt.Sprintf(APIUrl+VanityUrl, apiKey, vanityUrl)
	data, err := web.GetData(url)
	if err != nil {
		return "", err
	}
	// 如果返回结果为42说明这可能是一个SteamId
	if gjson.ParseBytes(data).Get("response.success").Uint() == 42 {
		return vanityUrl, nil
	}
	return gjson.ParseBytes(data).Get("response.steamid").String(), nil
}

package fgopickup

import (
	"strconv"
	"time"

	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 列出所有的卡池
func listPickups(ctx *zero.Ctx) {
	service := service{}
	pickups, err := service.getPickups()

	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err, " 解释: 查询出错！"))
		return
	}

	msg := make(message.Message, len(*pickups))
	for i, pickup := range *pickups {
		msg[i] = ctxext.FakeSenderForwardNode(ctx, getMsgOfSinglePickup(pickup)...)
	}
	ctx.Send(msg)
}

// 以卡池id获取某一个卡池的详情
func pickupDetail(ctx *zero.Ctx) {
	pickupID, err := strconv.Atoi(ctx.State["args"].(string))
	if err != nil || pickupID <= 0 {
		ctx.Send("参数错误！")
		return
	}

	service := service{}
	detail, err := service.getPickupDetail(pickupID)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err, " 解释: 查询出错！"))
		return
	}
	servants := detail.Servants

	days, err := service.getPickupTimeGap(pickupID)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err, " 解释: 查询出错！"))
		return
	}

	msg := make(message.Message, len(servants)+3)
	msg[0] = ctxext.FakeSenderForwardNode(ctx, getMsgOfSinglePickup(detail.Pickup)...)
	msg[1] = ctxext.FakeSenderForwardNode(ctx, message.Text("距今还有", days, "天. ", "以下为UP从者"))

	for i, servant := range servants {
		avatar := message.Image(servant.Avatar)
		name := message.Text("\n", servant.Name)
		msg[i+2] = ctxext.FakeSenderForwardNode(ctx, avatar, name)
	}
	ctx.Send(msg)
}

func getServantPickups(ctx *zero.Ctx) {
	servantID, err := strconv.Atoi(ctx.State["args"].(string))
	if err != nil || servantID <= 0 {
		ctx.Send("参数错误！")
		return
	}

	service := service{}
	res, err := service.getServantPickups(servantID)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err, " 解释: 查询出错！"))
		return
	}

	msg := make(message.Message, len(res.Pickup)+1)
	msg[0] = ctxext.FakeSenderForwardNode(ctx, message.Text("从者:<", res.ServantName, ">未来卡池如下"))
	for i, pickup := range res.Pickup {
		msg[i+1] = ctxext.FakeSenderForwardNode(ctx, getMsgOfSinglePickup(pickup)...)
	}
	ctx.Send(msg)
}

func getMsgOfSinglePickup(pickup pickup) message.Message {
	banner := message.Image(pickup.Banner)
	id := message.Text("\n", strconv.Itoa(pickup.ID), ". ")
	name := message.Text(pickup.Name)
	date := message.Text("\n",
		parseTime(pickup.StartTime), "~", parseTime(pickup.EndTime))
	return message.Message{banner, id, name, date}
}

func parseTime(timeInSeconds int64) string {
	return time.Unix(timeInSeconds, 0).Format("2006-01-02")
}

func getServantList(ctx *zero.Ctx) {
	page, _ := strconv.Atoi(ctx.State["args"].(string))
	service := service{}
	servants, err := service.listServants(page)
	if err != nil || len(*servants) == 0 {
		ctx.SendChain(message.Text("ERROR: ", err, " 解释: 查询出错！"))
		return
	}
	msg := make(message.Message, len(*servants))
	for i, servant := range *servants {
		avatar := message.Image(servant.Avatar)
		desc := message.Text("\n", servant.ID, ". ", servant.Name)
		msg[i] = ctxext.FakeSenderForwardNode(ctx, avatar, desc)
	}
	ctx.Send(msg)
}

package rss_hub

import (
	"context"
	"fmt"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/rss_hub/rss_pkg"
	"github.com/FloatTech/floatbox/ctxext"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zbpCtxExt "github.com/FloatTech/zbputils/ctxext"
	"github.com/fumiama/cron"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 初始化 repo
var (
	rssCron          = cron.New(cron.WithSeconds())
	rssRepo, initErr = rss_pkg.NewRssDomain(engine.DataFolder() + "rss_hub.db")

	// getRssRepo repo 初始化方法，单例
	getRssRepo = ctxext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		_, _ = rssCron.AddFunc("0 0/1 * * * ?", func() {
			rssSync(ctx)
		})
		rssCron.Start()
		logrus.WithContext(context.Background()).Infoln("RSS订阅姬：初始化")
		if initErr != nil {
			ctx.SendChain(message.Text("RSS订阅姬：初始化失败", initErr.Error()))
			return false
		}
		return true
	})
)

var (
	// 注册插件
	engine = control.Register("RssHub", &ctrl.Options[*zero.Ctx]{
		// 默认不启动
		DisableOnDefault: false,
		Brief:            "RssHub订阅姬",
		// 详细帮助
		Help: "RssHub订阅姬desu~ 支持的详细订阅列表可见 https://rsshub.netlify.app/ \n" +
			"- 添加RssHub订阅[订阅源路由] \n" +
			"- 删除RssHub订阅[订阅源路由] \n" +
			"- 查看RssHub订阅列表 \n" +
			"/启动RssHub同步 \n" +
			"/关闭RssHub同步 \n" +
			"/测试RssHub同步 \n",
		// 插件数据存储路径
		PublicDataFolder: "RssHub",
		OnEnable: func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("RSS订阅姬现在启动了哦"))
		},
		OnDisable: func(ctx *zero.Ctx) {
			rssCron.Stop()
			ctx.SendChain(message.Text("RSS订阅姬现在关闭了哦"))
		},
	}).ApplySingle(zbpCtxExt.DefaultSingle)
)

// init 命令路由
func init() {
	// Manage RssHub
	engine.OnCommand("强制RssHub同步", zero.AdminPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
	})
	// 启动
	engine.OnFullMatch("启动RssHub同步", zero.AdminPermission, getRssRepo).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("RssHub同步任务启动ing..."))
		rssCron.Start()
		ctx.SendChain(message.Text("RssHub同步任务启动成功"))
	})
	// 停止
	engine.OnFullMatch("停止RssHub同步", zero.AdminPermission, getRssRepo).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		rssCron.Stop()
		ctx.SendChain(message.Text("RSS订阅姬：同步任务已停止"))
	})
	// 添加订阅
	engine.OnRegex(`^添加RssHub订阅(.+)$`, zero.OnlyGroup, getRssRepo).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		routeStr := ctx.State["regex_matched"].([]string)[1]
		rv, _, isSubExisted, err := rssRepo.Subscribe(context.Background(), ctx.Event.GroupID, routeStr)
		if err != nil {
			ctx.SendChain(message.Text("RSS订阅姬：添加失败", err.Error()))
			return
		}
		if isSubExisted {
			ctx.SendChain(message.Text("RSS订阅姬：已存在，更新成功"))
		} else {
			ctx.SendChain(message.Text("RSS订阅姬：添加成功"))
		}
		// 添加成功，发送订阅源信息
		msg := make(message.Message, 0)
		rawMsgSlice := formatRssFeedToTextMsg(rv)
		for _, rm := range rawMsgSlice {
			msg = append(msg, fakeSenderForwardNode(ctx.Event.SelfID, message.Text(rm)))
		}
		//m := message.Message{zbpCtxExt.FakeSenderForwardNode(ctx, msg...)}
		if id := ctx.Send(msg).ID(); id == 0 {
			ctx.SendChain(message.Text("ERROR: 可能被风控了"))
		}
		//ctx.SendChain(msg...)
	})
	engine.OnRegex(`^删除RssHub订阅(.+)$`, zero.OnlyGroup, getRssRepo).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		routeStr := ctx.State["regex_matched"].([]string)[1]
		err := rssRepo.Unsubscribe(context.Background(), ctx.Event.GroupID, routeStr)
		if err != nil {
			ctx.SendChain(message.Text("RSS订阅姬：删除失败 ", err.Error()))
			return
		}
		// 添加成功，发送订阅源信息
		var msg []message.MessageSegment
		msg = append(msg, message.Text(fmt.Sprintf("RSS订阅姬：删除%s成功", routeStr)))
		ctx.SendChain(msg...)
	})
	engine.OnFullMatch("RssHub订阅列表", zero.OnlyGroup, getRssRepo).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		rv, err := rssRepo.GetSubscribedChannelsByGroupId(context.Background(), ctx.Event.GroupID)
		if err != nil {
			ctx.SendChain(message.Text("RSS订阅姬：查询失败 ", err.Error()))
			return
		}
		// 添加成功，发送订阅源信息
		var msg []message.MessageSegment
		msg = append(msg, message.Text("RSS订阅姬：当前订阅列表"))
		for _, v := range rv {
			msg = append(msg, message.Text(formatRssFeedToTextMsg(v)))
		}
		ctx.SendChain(msg...)
	})

}

func rssSync(ctx *zero.Ctx) {
	m, ok := control.Lookup("RssHub")
	if !ok {
		logrus.Warn("RssHub插件未启用")
		return
	}
	// map[群组]推送内容Views
	groupToFeedsMap, err := rssRepo.SyncJobTrigger(context.Background())
	if err != nil {
		ctx.SendChain(message.Text("RSS订阅姬：同步任务失败 ", err.Error()))
		return
	}
	// 没有更新的群组不推送
	if len(groupToFeedsMap) == 0 {
		logrus.Info("RssHub未发现更新")
		return
	}
	// 遍历群组
	for i, views := range groupToFeedsMap {
		logrus.Infof("RssHub插件在群 %d 触发推送检查", i)
		//msg = append(msg, message.Text("RSS订阅姬定时推送中~\n"))
		// 遍历推送Views，看哪些频道有更新
		isEnabledInGroup := m.IsEnabledIn(i)
		for _, v := range views {
			msg := make(message.Message, 0)
			// 没有更新的频道不推送
			if len(v.Contents) == 0 {
				continue
			}
			rawMsgSlice := formatRssFeedToTextMsg(v)
			for _, rm := range rawMsgSlice {
				msg = append(msg, fakeSenderForwardNode(ctx.Event.SelfID, message.Text(rm)))
			}
			//msg = append(msg, message.Text(formatRssFeedToTextMsg(v)))
			// 没有更新的群组不推送
			if len(msg) == 0 {
				continue
			}
			// 检查群组是否启用了插件
			if isEnabledInGroup {
				logrus.Infof("RssHub插件在群 %d 开始推送", i)
				ctx.SendGroupMessage(i, message.Text(v.Channel.Title+"\n[RSS订阅姬定时推送]\n"))
				//ms := message.Message{zbpCtxExt.FakeSenderForwardNode(ctx, msg)}
				if id := ctx.Send(msg).ID(); id == 0 {
					ctx.SendChain(message.Text("ERROR: 可能被风控了"))
				}
			}
		}
	}
}

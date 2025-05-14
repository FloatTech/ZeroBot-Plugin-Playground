// Package example 这是zbp的插件编写教学示例
package example

// import 用来放置你所需要导入的东西, 萌新推荐使用vscode, 它会帮你干很多事
import (
	"math/rand"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 自定义限制函数, 括号内填入(时间,触发次数)
var examplelimit = ctxext.NewLimiterManager(time.Second*10, 1)

// 这里就是插件主体了
func init() {
	// 既然是zbp, 那就从接入control开始, 在这里注册你的插件以及设置是否默认开启和填写帮助和数据存放路径
	engine := control.Register("example", &ctrl.Options[*zero.Ctx]{
		// 控制插件是否默认启用 true为默认不启用 false反之
		DisableOnDefault: false,
		// 插件的简介
		Brief: "例",
		// 插件的帮助 管理员发送 /用法 example 可见
		Help: "- example 插件的帮助",
		// 插件的背景图, 支持http和本地路径
		// Banner: "",
		// 插件的数据存放路径, 分为公共和私有, 都会在/data下创建目录, 公有需要首字母大写, 私有需要首字母小写
		PublicDataFolder: "Example",
		// PrivateDataFolder: "example",		// 避免问题所以注释了
		// 自定义插件开启时的回复
		OnEnable: func(ctx *zero.Ctx) {
			ctx.Send("插件已启用")
		},
		// 自定义插件关闭时的回复
		OnDisable: func(ctx *zero.Ctx) {
			ctx.Send("插件已禁用")
		},
	})
	// OnFullMatch 完全匹配触发器: 顾名思义, 收到消息 完全匹配 时就会触发, 所以快下个vscode吧（（（
	engine.OnFullMatch("完全匹配触发器").
		// SetBlock 设置是否阻断: 可选参数 设置是否阻断后续的触发, 也就是说如果设置为true, 后续的触发器不会被触发, false反之
		SetBlock(true).
		// Limit 限速器: 可选参数 限制设置的时间内触发器能够触发多少次, 已经有封装好的限制函数, 可以直接使用
		// Limit(ctxext.LimitByGroup) 按群号限制10s内5次触发
		// Limit(ctxext.LimitByUser) 按用户限制10s内5次触发
		// 也可以自己定义限制函数
		Limit(ctxext.LimitByGroup).
		// Handle 处理事件: 必要参数 直接处理事件, 括号内需要填入func(ctx *zero.Ctx){}
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("完全匹配触发器"))
		})
		// 自定义limit使用例
	engine.OnFullMatch("完全匹配触发器2").SetBlock(true).Limit(examplelimit.LimitByGroup).
		// Handle内也可以填入已经封装好的函数
		Handle(onmessage)
	// OnFullMatchGroup 完全匹配触发器组: 支持多个触发词的完全匹配触发器
	engine.OnFullMatchGroup([]string{"完全匹配触发器组", "OnFullMatchGroup"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("完全匹配触发器组\nOnFullMatchGroup"))
	})
	// OnKeyword 关键词匹配触发器: 也就是说一段消息中含有 关键词匹配触发器 时就会触发
	engine.OnKeyword("关键词匹配触发器").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("关键词匹配"))
	})
	// OnKeywordGroup 关键词匹配触发器组: 支持多个触发词的关键词匹配触发器
	engine.OnKeywordGroup([]string{"关键词匹配触发器组", "OnKeywordGroup"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("关键词匹配触发器组\nOnKeywordGroup"))
	})
	// OnCommand 命令匹配触发器: 需要在触发的文本前加上设置的prefix(前缀), 默认前缀为 /
	engine.OnCommand("命令触发器").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("命令触发器"))
	})
	// OnCommandGroup 命令匹配触发器组: 支持多个触发词的命令匹配触发器
	engine.OnCommandGroup([]string{"命令触发器组", "OnCommandGroup"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("命令触发器组\nOnCommandGroup"))
	})
	// OnPrefix 前缀匹配触发器: 一段消息内的开头为 前缀触发器 时, 不论后面的内容是什么, 就会触发
	engine.OnPrefix("前缀触发器").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("前缀触发器"))
	})
	// OnPrefixGroup 前缀匹配触发器组: 支持多个触发词的前缀匹配触发器
	engine.OnPrefixGroup([]string{"前缀触发器组", "OnPrefixGroup"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("前缀触发器组\nOnPrefixGroup"))
	})
	// OnSuffix 后缀匹配触发器: 一段消息内的结尾为 后缀触发器 时, 不论前面的内容是什么, 就会触发
	engine.OnSuffix("后缀触发器").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("后缀触发器"))
	})
	// OnSuffixGroup 后缀匹配触发器组: 支持多个触发词的后缀匹配触发器
	engine.OnSuffixGroup([]string{"后缀触发器组", "OnSuffixGroup"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("后缀触发器组\nOnSuffixGroup"))
	})
	// OnRegex 正则匹配触发器: 一段消息内的内容与正则表达式匹配时, 就会触发
	// 推荐一个在线正则网站 https://regex101.com/ 进去后在左边的FLAVOR选择Golang
	engine.OnRegex(`^正则表达式匹配触发器$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("正则表达式匹配触发器"))
	})
	engine.OnFullMatch("完全匹配触发器3").SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		// 这里是go语言的特性 使用 := 声明并赋值变量, 生命周期为这个func
		// ctx.Event.UserID 则是发送这条消息的人的qq号, 给下面的at(@)使用
		uid := ctx.Event.UserID
		// SendChain 允许你无限拼接message下去
		ctx.SendChain(message.At(uid), message.Text("Hello world!"), message.Text("\n你好, 世界！"))
	})
	engine.OnFullMatch("随机文本").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// 发送文本支持多条随机
		ctx.SendChain(message.Text([]string{
			"你好啊",
			"做出了回应",
			"好的哟",
		}[rand.Intn(3)]))
		// 另一种实现 来自zbp的插件 atri
		ctx.SendChain(randText(
			"你好啊",
			"做出了回应",
			"好的哟",
		))
	})
	// 加上 zero.OnlyToMe 这个规则 变成了只有 @bot 或者 叫bot名字 才能触发
	engine.OnFullMatch("OnlyToMe", zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// Send 快捷发送
		ctx.Send(
			message.Message{
				message.Text("找我有什么事啊？"),
				message.Text("找我有什么事啊？"),
				message.Text("找我有什么事啊？"),
			},
		)
	})
	// 还有很多规则, 例如
	// zero.OnlyGroup 仅在群组中触发
	// zero.OnlyPrivate 仅在私聊中触发
	// zero.OnlyGuild 仅在频道中触发
	// zero.OnlyPublic 仅在群组和频道中都触发
	// 规则实际上是一个bool值, 所以也可以自定义规则, 只有true时才会继续执行
	engine.OnKeyword("关键词匹配触发器2", zero.OnlyGroup, zero.OnlyPrivate, zero.OnlyGuild, zero.OnlyPublic).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		// 这里避坑！如果带有 Poke|Record , 那么后面的东西就发不出来了
		ctx.SendChain(message.Poke(uid), message.Record("填入在线音频链接或者file:///加上本地路径"), message.Text("关键词匹配触发器2"))
		// 解决办法是分开进行发送
		ctx.SendChain(message.Poke(uid))
		ctx.SendChain(message.Record("填入在线音频链接或者file:///加上本地路径"))
		ctx.SendChain(message.Text("关键词匹配触发器2"))
	})
	// 自定义规则使用例
	engine.OnKeywordGroup([]string{"我在哪发消息给你", "我在哪"}, customrule).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// go语言的特性, 声明一个变量
		var uid, gid int64
		// 发送者所在群号, 当不在群时默认为0, 赋值群号到gid上, 因为声明过了, 无需:=
		gid = ctx.Event.GroupID
		// 发送者qq号
		uid = ctx.Event.UserID
		// 规则是个bool值, 所以可以写成if的形式来判断
		if zero.OnlyGroup(ctx) {
			// 发送消息到指定群号
			ctx.SendGroupMessage(gid, message.Text("你在群里"))
			return
		}
		if zero.OnlyPrivate(ctx) {
			// 发送消息到指定QQ号, 需要添加对方为好友
			ctx.SendPrivateMessage(uid, message.Text("你在私聊里"))
			return
		}
	})
	engine.OnFullMatch("我是什么职位", customrule2).SetBlock(true).Handle(func(_ *zero.Ctx) {
	})
}

// randText 随机文本
func randText(text ...string) message.Segment {
	return message.Text(text[rand.Intn(len(text))])
}

// 可以填入Handle的封装好的函数
func onmessage(ctx *zero.Ctx) {
	ctx.SendChain(message.Text("完全匹配触发器2"))
}

// 自定义的规则, 只有群组和私聊中能触发
func customrule(ctx *zero.Ctx) bool {
	if zero.OnlyGroup(ctx) || zero.OnlyPrivate(ctx) {
		return true
	}
	return false
}

// 自定义的规则, 判断你处于什么权限
func customrule2(ctx *zero.Ctx) bool {
	// 判断是否是主人
	if zero.SuperUserPermission(ctx) {
		ctx.Send("你是主人")
		return true
	}
	// 判断是否是群主及以上
	if zero.OwnerPermission(ctx) {
		ctx.Send("你是群主")
		return true
	}
	// 判断是否是管理员及以上
	if zero.AdminPermission(ctx) {
		ctx.Send("你是管理员")
		return true
	}
	// 判断是否是私聊者或管理员以上
	if zero.UserOrGrpAdmin(ctx) {
		ctx.Send("你是私聊者")
		return true
	}
	// 如果都不是, 则是普通成员
	ctx.Send("你是普通成员")
	return true
}

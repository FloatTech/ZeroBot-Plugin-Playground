// Package example 这是zbp的插件编写教学示例
package example

// import 用来放置你所需要导入的东西，萌新推荐使用vscode，它会帮你干很多事
import (
	"math/rand"

	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 这里就是插件主体了
func init() {
	// 既然是zbp，那就从接入control开始，在这里注册你的插件以及设置是否默认开启和填写帮助和数据存放路径
	engine := control.Register("example", &control.Options{
		// 控制插件是否默认启用 true为默认不启用 false反之
		DisableOnDefault: false,
		Help:             "",
	})
	// 完全匹配触发器，顾名思义，收到消息 test 时就会触发，所以快下个vscode吧（（（
	engine.OnFullMatch("test").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			// 这里是go语言的特性 使用 := 声明并赋值变量，生命周期为这个func
			// ctx.Event.UserID 则是发送这条消息的人的qq号，给下面的at(@)使用
			uid := ctx.Event.UserID
			// SendChain 允许你无限拼接message下去
			ctx.SendChain(message.At(uid), message.Text("Hello,world!"), message.Text("\n你好，世界！"))
		})
	// 完全匹配触发器组，支持多个触发词
	engine.OnFullMatchGroup([]string{"测试", "你好"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			// 同样，发送文本也支持多条随机
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
	// 同样是完全匹配触发器，加上 zero.OnlyToMe 变成了只有 @bot 或者 叫bot名字 才能触发
	engine.OnFullMatch("Test", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			// Send 快捷发送
			ctx.Send("找我有什么事啊？")
		})
	// 关键词匹配触发器，也就是说，一段话内，有 戳我 时就会触发
	engine.OnKeyword("戳我").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			// 这里避坑！如果带有 Poke ，那么后面的东西就发不出来了
			ctx.SendChain(message.Poke(uid), message.Text("戳你！"))
			// 解决办法是分成两条
			ctx.SendChain(message.Poke(uid))
			ctx.SendChain(message.Text("戳你！"))
		})
}

// randText 随机文本
func randText(text ...string) message.MessageSegment {
	return message.Text(text[rand.Intn(len(text))])
}

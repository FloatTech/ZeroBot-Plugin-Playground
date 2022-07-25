// Package example2 这是小锅饭的zbp的插件编写教学示例
package example2

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type user struct {
	Name string `flag:"n"`
	Age  int    `flag:"a"`
}

func init() {
	engine := control.Register("example2", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "小锅饭的示例\n- hello\n- 完全匹配\n- 完全匹配组1 | 完全匹配组2\n- 关键词你好\n- /命令\n- 前缀你好\n- 你好后缀\n" +
			"- 小锅饭正在洗澡\n- /用户 -n 小锅饭 -a 3\n- 牛逼\n- 消息检测\n- 文本\n- 艾特\n" +
			"- 回复\n- 语音\n- 表情\n- 音乐\n- tts\n- 戳一戳\n",
	})
	// hello
	engine.OnFullMatch("hello").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("hello world!"))
	})

	// matcher
	engine.OnFullMatch("完全匹配").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		matched := ctx.State["matched"].(string)
		ctx.SendChain(message.Text("完全匹配的匹配词: ", matched))
	})
	engine.OnFullMatchGroup([]string{"完全匹配组1", "完全匹配组2"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		matched := ctx.State["matched"].(string)
		ctx.SendChain(message.Text("完全匹配组的匹配词: ", matched))
	})
	engine.OnKeyword("关键词").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		keyword := ctx.State["keyword"].(string)
		ctx.SendChain(message.Text("关键词匹配的关键词: ", keyword))
	})
	engine.OnCommand("命令").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		command := ctx.State["command"].(string)
		args := ctx.State["args"].(string)
		ctx.SendChain(message.Text("命令匹配的命令: ", command, "\n命令匹配的参数: ", args))
	})
	engine.OnPrefix("前缀").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		prefix := ctx.State["prefix"].(string)
		args := ctx.State["args"].(string)
		ctx.SendChain(message.Text("前缀匹配的前缀: ", prefix, "\n前缀匹配的参数: ", args))
	})
	engine.OnRegex(`(.*)正在(.*)`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		regexMatched := ctx.State["regex_matched"].([]string)
		ctx.SendChain(message.Text("正则匹配的匹配组: ", regexMatched))
	})
	engine.OnSuffix("后缀").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		suffix := ctx.State["suffix"].(string)
		args := ctx.State["args"].(string)
		ctx.SendChain(message.Text("后缀匹配的后缀: ", suffix, "\n后缀匹配的参数: ", args))
	})
	engine.OnShell("用户", user{}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		args := ctx.State["args"].([]string)
		u := ctx.State["flag"].(*user)
		ctx.SendChain(message.Text("shell匹配的结构体: ", u, "\nshell匹配的参数: ", args))
	})

	// message
	engine.OnFullMatch("文本").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("文本"))
	})
	engine.OnFullMatch("艾特").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("艾特你"))
	})
	engine.OnFullMatch("回复").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("回复你"))
	})
	engine.OnFullMatch("图片").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Image("https://gitcode.net/anto_july/imagematerials/-/blob/main/need/0.png"))
	})
	engine.OnFullMatch("语音").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Record("https://vtbkeyboard.moe/api/audio/672328094/%E7%8C%AB%E7%8C%AB%E8%B7%9F%E9%BC%A0%E9%BC%A0%E4%B8%8D%E5%86%B2%E7%AA%81.mp3"))
	})
	engine.OnFullMatch("表情").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Face(1))
	})
	engine.OnFullMatch("音乐").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Music("163", 28949129))
	})
	engine.OnFullMatch("tts").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.TTS("你好"))
	})
	engine.OnFullMatch("戳一戳").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Poke(ctx.Event.UserID))
	})

	// rule
	engine.OnFullMatch("牛逼", zero.AdminPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("你很牛逼"))
	})
	engine.OnFullMatch("消息检测", checkRule).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("检测完毕"))
	})
}

func checkRule(ctx *zero.Ctx) bool {
	if zero.OnlyPrivate(ctx) {
		ctx.SendChain(message.Text("这是私聊信息"))
	}
	if zero.OnlyGroup(ctx) {
		ctx.SendChain(message.Text("这是群聊信息"))
	}
	if zero.OnlyToMe(ctx) {
		ctx.SendChain(message.Text("这是@bot的信息"))
	}
	if zero.OnlyGuild(ctx) {
		ctx.SendChain(message.Text("这是频道信息"))
	}
	if zero.OnlyPublic(ctx) {
		ctx.SendChain(message.Text("这是群聊或者频道信息"))
	}
	if zero.AdminPermission(ctx) {
		ctx.SendChain(message.Text("你有管理员权限"))
	}
	if zero.OwnerPermission(ctx) {
		ctx.SendChain(message.Text("你有群主权限"))
	}
	if zero.SuperUserPermission(ctx) {
		ctx.SendChain(message.Text("你有主人权限"))
	}
	return true
}

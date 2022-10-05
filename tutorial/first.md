## 浅尝辄止~尝试做一个好看的拼盘吧w

经过了上一次的尝试之后~ 我们已经踏入了第一步 > 让机器人跑自己需要的插件w
 

啊咧~我总不能只想让我的机器人发送 hello world 吧 我还想整点别的w

那好啊~我们就开始编写啊w



### emmm...我想让我的机器人可以发送图片w

好哦~w 记得咱之前提的嘛 在**engine**中

我们可以通过调用**ctx**工具箱的方式 为我们的机器人使用专业对口的小东西

```go
engine.OnFullMatch("hello").Handle(func(ctx *zero.Ctx) { 
		ctx.SendChain(message.Image("https://manual-lucy.himoyo.cn/test.png"))
	}) 
```

此刻 我们的机器人在收到信息 hello 后 就会发送来自于 "https://manual-lucy.himoyo.cn/test.png"

的图片 ~~私货(((~~

以此类推 我们也可以条用 Record 的方式 来进行相关语音的发送

相关可以调用的如下 : 

```go
       // 这里避坑！如果带有 Poke|Record , 那么后面的东西就发不出来了
		ctx.SendChain(message.Poke(uid), message.Record("填入在线音频链接或者file:///加上本地路径"), message.Text("关键词匹配触发器2"))
		// 解决办法是分开进行发送
		ctx.SendChain(message.Poke(uid)) //戳一戳
		ctx.SendChain(message.Record("填入在线音频链接或者file:///加上本地路径"))
		ctx.SendChain(message.Text("关键词匹配触发器2"))
		// 发送消息到指定群号
		ctx.SendGroupMessage(gid, message.Text("你在群里"))
       // 发送消息到指定QQ号, 需要添加对方为好友
		ctx.SendPrivateMessage(uid, message.Text("你在私聊里"))

// gid 指的是 活动群组 一般使用 ctx.Event.GroupID | uid 即 指的是 活动用户 使用 ctx.Event.UserID
```

## 好啦....感觉这样没意思 那?我还有什么好玩的呢?

其实呢~在zbp插件中 就有一处挺好玩的

```go
engine.OnFullMatch("", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var nickname = zero.BotConfig.NickName[0]
			ctx.SendChain(message.Text(
				[]string{
					nickname + "在此，有何贵干~",
					"(っ●ω●)っ在~",
					"这里是" + nickname + "(っ●ω●)っ",
					nickname + "不在呢~",
				}[rand.Intn(4)],
			))
		})
```

通过听取bot的名字 实现了喊到即出现的情景~算是一个不错的案例w



同时 你也可以~

```go
// 戳一戳
	engine.On("notice/notify/poke", zero.OnlyToMe).SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			var nickname = zero.BotConfig.NickName[0]
			switch {
			case poke.Load(ctx.Event.GroupID).AcquireN(3):
				// 5分钟共8块命令牌 一次消耗3块命令牌
				time.Sleep(time.Second * 1)
				ctx.SendChain(message.Text("请不要戳", nickname, " >_<"))
			case poke.Load(ctx.Event.GroupID).Acquire():
				// 5分钟共8块命令牌 一次消耗1块命令牌
				time.Sleep(time.Second * 1)
				ctx.SendChain(message.Text("喂(#`O′) 戳", nickname, "干嘛！"))
			default:
				// 频繁触发，不回复
			}
		})
```

利用命令牌的方式 实现了~相关的戳一戳生气w






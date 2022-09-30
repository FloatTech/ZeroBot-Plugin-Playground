## Hey ya~ Hello World~

貌似编程的第一步都是先跑个HelloWorld呢~

既然我们要先学第一步 那就先写个 Hello World 玩吧~


这是第一个实例

```go
package main // 包名

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
    "github.com/wdvxdr1123/ZeroBot/message"
    zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	// 既然是zbp, 那就从接入control开始, 在这里注册你的插件以及设置是否默认开启和填写帮助
	engine := control.Register("example", &ctrl.Options[*zero.Ctx]{
		// 控制插件是否默认启用 true为默认不启用 false反之
		DisableOnDefault: false,
		// 插件的帮助 管理员发送 /用法 example 可见
		Help: "- example 插件的帮助",
	})
	engine.OnFullMatch("hello").Handle(func(ctx *zero.Ctx) { 
        // OnFullMath 即 全局匹配 当收到此消息时做出的反应
		ctx.Send(message.Text("hello world!"))
         // ctx.Send 即 快捷发送消息 | 合并转发 
		 // 支持 ctx.Send(("hello world!")
	}) 
}

```

其中 

```go
	engine := control.Register("example", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "- example 插件的帮助",
	})
```

在Golang中 我们是通过引用包中的"轮子"来实现机器人交互

这边使用的包是 ZeroBot | 很多功能我们都是通过引用来实现功能的w

在第一步 我们注册了这个插件 即利用 "Control" 被定义为 **engine** 即驱动的方式 注册了一个叫 "example" 的插件

我们大多数的功能都是需要走轮子 即已经是被调用的Control的**engine**来实现

所以呢 匹配的方式的前缀 都可以看到 "engine" 为开头

既然工具注册好了 那么我们就应该使用相关的工具来实现相关的功能了

这个工具呢 叫做 "ctx" 它可以帮助我们执行很多相关的功能

首先呢 我们需要使用监听的工具 而engine就可以帮助我们实现w

其中 **engine.OnFullMatch("hello")** 就是通过的**OnFullMatch**工具去识别收到的是否有这些东西


那么 我们既然使用了**OnFullMatch**工具 那么我们的监听听到了**hello**后 就会通过一个小小的**handle**,去执行我们需要的执行的东西

即

```
.Handle(func(ctx *zero.Ctx) { 
	}) 
```

通过Handle小工具 我们可以执行接下来的操作~

那~既然我要发送 那我就应该告诉程序 我应该要发送什么?

当然~ **ctx**工具和**engine**是一样的~都有自己需要的格式来进行发送

之所以分开使用 因为各有所长~所以在实际使用中我们需要使用不同的工具来应对所需

诸如~ **send**工具就可以用来调用其中的发送包~  帮助你发送所需要的内容~ 主要是发送文字

而他的大师兄~**sendChain**即是多样的发送工具 可以调用更多的功能w


当我们通过导入插件的方式运行时~ 发送 **hello** 即收到回复 **hello world**

至此 我们的**Hello World**已经成功跑起来了哦w 这是你编写的第一个插件w




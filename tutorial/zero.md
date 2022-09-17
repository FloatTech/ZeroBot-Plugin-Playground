## Hey ya~ Hello World~

~~貌似编程的第一步都是先跑个HelloWorld呢~~

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
         // ctx.Send 即发送单一类型信息 仅支持发送一条类似的
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

这是我们要做的第一步~即 将这个插件的信息报告给总部~

这个总部在golang中 是被定义为 **engine** 

所以呢 我们许多的功能 都需要通过总部发声 即是 **engine.**

让总部知道有这个插件的存在,总部才会认可这个插件并且去使用~

其中 **Help** 一类 即是报告他的内容~让总部知道这些插件有什么用 总部会记录下并且在用户需要的时候调出~



既然是总部发言 那么 我们就需要通过总部 收到消息后执行相关的指令

这就意味这我们需要通过总部调用相关的匹配工具 让总部能够收到并且能正确的指引我方使用正确的指令

其中 **engine.OnFullMatch("hello")** 就是通过总部的**OnFullMatch**工具去识别收到的是否有这些东西

那么 我们既然使用了**OnFullMatch**工具 那么我们的监听器听到了**hello**后 就会通过一个小小的**handle**,去执行我们需要的执行的东西

即

```
.Handle(func(ctx *zero.Ctx) { 
	}) 
```

通过Handle小工具 我们可以执行接下来的操作~



那~既然我要发送 那我就应该告诉程序 我应该要发送什么?

那~**ctx**就可以作为我们的信纸,通过格式化的信息,清晰的去传递所需要的内容

当然~**ctx**工具和**engine**是一样的~都有自己需要的格式来进行发送

~之所以分开使用 因为各有所长~所以在实际使用中我们需要使用不同的工具来应对所需

诸如~**send**工具就可以用来调用其中的发送包~帮助你发送所需要的内容~主要是发送文字

而他的大师兄~**sendChain**即是多样的发送工具 可以调用更多的功能w



当我们通过导入插件的方式运行时~ 发送 **hello** 即收到回复 **hello world**

至此 我们的**Hello World**已经成功跑起来了哦w 这是你编写的第一个插件w




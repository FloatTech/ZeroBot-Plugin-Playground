package miaox

import (
	"context"
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/chain"
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/cmd"
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/repo"
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/repo/store"
	pTypes "github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/types"
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/utils"
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/vars"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/bincooo/MiaoX"
	"github.com/bincooo/MiaoX/types"
	xvars "github.com/bincooo/MiaoX/vars"
	wapi "github.com/bincooo/openai-wapi"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var help = `
- @Bot + 文本内容
- 昵称前缀 + 文本内容
- 预设列表
- [开启|切换]预设 + [预设名]
- 删除凭证 + [key]
- 添加凭证 + [key]:[value]
- 切换AI + [AI类型：openai-api、openai-web、claude、bing-(c|b|p|s)]
- 作画 + [tags]
tips：
  配置项采用web方式配置：请访问 http://127.0.0.1:8082 查看
`
var (
	engine = control.Register("miaox", &ctrl.Options[*zero.Ctx]{
		Help:              help,
		Brief:             "AI适配器",
		DisableOnDefault:  false,
		PrivateDataFolder: "miaox",
	})

	//mgr types.BotManager
	lmt types.Limiter
)

func init() {
	vars.E = engine

	var err error
	if vars.Loading, err = os.ReadFile(engine.DataFolder() + "load.gif"); err != nil {
		panic(err)
	}

	lmt = MiaoX.NewCommonLimiter()
	if e := lmt.RegChain("args", &chain.ArgsInterceptor{}); e != nil {
		panic(e)
	}
	if e := lmt.RegChain("online", &chain.OnlineInterceptor{}); e != nil {
		panic(e)
	}

	engine.OnRegex(`^添加凭证\s+(\S+)`, zero.OnlyPrivate, repo.OnceOnSuccess).SetBlock(true).
		Handle(insertTokenCommand)
	engine.OnRegex(`^删除凭证\s+(\S+)`, zero.OnlyPrivate, repo.OnceOnSuccess).SetBlock(true).
		Handle(deleteTokenCommand)
	engine.OnFullMatch("凭证列表", zero.OnlyPrivate, repo.OnceOnSuccess).SetBlock(true).
		Handle(tokensCommand)
	engine.OnRegex(`[开启|切换]预设\s(\S+)`, zero.OnlyToMe, repo.OnceOnSuccess).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(enablePresetSceneCommand)
	engine.OnRegex(`切换AI\s(\S+)`, zero.AdminPermission, repo.OnceOnSuccess).SetBlock(true).
		Handle(switchAICommand)
	engine.OnFullMatch("预设列表", zero.OnlyToMe, repo.OnceOnSuccess).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(presetScenesCommand)
	engine.OnPrefix("作画", repo.OnceOnSuccess).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(drawCommand)
	engine.OnMessage(zero.OnlyToMe, repo.OnceOnSuccess, excludeOnMessage).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(conversationCommand)

	cmd.Register("/api/global", repo.GlobalService{}, cmd.NewMenu("global", "全局配置"))
	cmd.Register("/api/preset", repo.PresetService{}, cmd.NewMenu("preset", "预设配置"))
	cmd.Register("/api/token", repo.TokenService{}, cmd.NewMenu("token", "凭证配置"))

	Run(":8082")
}

// 自定义优先级
func excludeOnMessage(ctx *zero.Ctx) bool {
	msg := ctx.MessageString()
	exclude := []string{"添加凭证 ", "删除凭证 ", "凭证列表", "开启预设 ", "切换预设 ", "预设列表", "切换AI ", "作画", "/", "!"}
	for _, value := range exclude {
		if strings.HasPrefix(msg, value) {
			return false
		}
	}
	return true
}

// 聊天
func conversationCommand(ctx *zero.Ctx) {
	cctx, err := createConversationContext(ctx, "")
	if err != nil {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("发生异常: "+err.Error()))
		return
	}

	cctx.Prompt = parseMessage(ctx)
	args := cctx.Data.(pTypes.ConversationContextArgs)
	args.Current = strconv.FormatInt(ctx.Event.Sender.ID, 10)
	args.Nickname = ctx.Event.Sender.NickName
	cctx.Data = args

	response := make(chan types.PartialResponse)
	if e := lmt.Join(cctx, response); e != nil {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(e.Error()))
		close(response)
		return
	}

	lmtHandle := func() {
		defer close(response)
		delay := utils.NewDelay(ctx)
		defer delay.Close()
		for {
			r := <-response
			if len(r.Message) > 0 {
				segment := utils.StringToMessageSegment(r.Message)
				ctx.SendChain(append(segment, message.Reply(ctx.Event.MessageID))...)
				delay.Defer()
			}
			if r.Error != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(r.Error.Error()))
				break
			}

			if r.Status == xvars.Closed {
				break
			}
		}
	}

	go lmtHandle()
}

// AI作画
func drawCommand(ctx *zero.Ctx) {
	prompt := ctx.State["args"].(string)
	if prompt == "" {
		return
	}

	prompt = strings.ReplaceAll(prompt, "，", ",")
	global := repo.GetGlobal()
	if global.DrawServ == "" {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请先联系管理员设置AI作画API"))
		return
	}

	logrus.Info("接收到作画请求，开始作画：", prompt)
	ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("这就开始画 ~"))
	imgBytes, err := utils.DrawAI(global.DrawServ, prompt, global.DrawBody)
	if err != nil {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("作画失败："+err.Error()))
		return
	}

	ctx.SendChain(message.Reply(ctx.Event.MessageID), message.ImageBytes(imgBytes))
}

// 添加凭证
func insertTokenCommand(ctx *zero.Ctx) {
	value := ctx.State["regex_matched"].([]string)[1]
	pattern := `^([^|]+)\:(.+)`
	r, _ := regexp.Compile(pattern)
	matches := r.FindStringSubmatch(value)
	logrus.Infoln(matches)
	if matches[1] == "" || matches[2] == "" {
		ctx.Send("添加失败，请按格式填写")
		return
	}
	global := repo.GetGlobal()
	billing, err := wapi.Query(context.Background(), matches[2], global.Proxy)
	if err != nil {
		logrus.Warn(err)
	}
	if billing.System-billing.Soft <= 0 {
		ctx.Send("添加失败，凭证余额为0")
		return
	}
	err = repo.InsertToken(repo.Token{
		Key:   matches[1],
		Token: matches[2],
		Type:  xvars.OpenAIAPI,
	})
	if err != nil {
		ctx.Send("添加失败: " + err.Error())
	} else {
		ctx.Send("添加成功，余额为" + strconv.FormatFloat(billing.System-billing.Soft, 'f', 2, 64))
	}
}

// 删除凭证
func deleteTokenCommand(ctx *zero.Ctx) {
	value := ctx.State["regex_matched"].([]string)[1]
	token := repo.GetToken(value)
	if token == nil {
		ctx.Send("`" + value + "`不存在")
		return
	}
	repo.RemoveToken(value)
	ctx.Send("`" + value + "`已删除")
}

// 凭证列表
func tokensCommand(ctx *zero.Ctx) {
	doc := "凭证列表：\n"
	tokens, err := repo.FindTokens("")
	if err != nil {
		ctx.Send(doc + "None.")
		return
	}
	if len(tokens) <= 0 {
		ctx.Send(doc + "None.")
		return
	}
	for _, token := range tokens {
		doc += token.Type + " | " + token.Key + "\n"
	}
	ctx.Send(doc)
}

// 开启/切换预设
func enablePresetSceneCommand(ctx *zero.Ctx) {
	value := ctx.State["regex_matched"].([]string)[1]
	presetScene := repo.GetPresetScene(value)
	if presetScene == nil {
		ctx.Send("`" + value + "`预设不存在")
		return
	}

	cctx, err := createConversationContext(ctx, "")
	if err != nil {
		ctx.Send("获取上下文出错: " + err.Error())
		return
	}
	if cctx.Bot != presetScene.Type {
		ctx.Send("当前AI类型无法使用`" + value + "`预设")
		return
	}

	cctx.Preset = presetScene.Content
	cctx.Format = presetScene.Message
	cctx.Chain += presetScene.Chain
	lmt.Remove(cctx.Bot)
	store.DeleteOnline(cctx.Id)
	updateConversationContext(cctx)
	ctx.Send("已切换`" + value + "`预设")
}

// 预设场景列表
func presetScenesCommand(ctx *zero.Ctx) {
	doc := "预设列表：\n"
	preset, err := repo.FindPresetScenes("")
	if err != nil {
		ctx.Send(doc + "None.")
		return
	}
	if len(preset) <= 0 {
		ctx.Send(doc + "None.")
		return
	}
	for _, token := range preset {
		doc += token.Type + " | " + token.Key + "\n"
	}
	ctx.Send(doc)
}

func switchAICommand(ctx *zero.Ctx) {
	bot := ctx.State["regex_matched"].([]string)[1]
	var cctx types.ConversationContext
	switch bot {
	case xvars.OpenAIAPI,
		xvars.OpenAIWeb,
		xvars.Claude,
		xvars.Bing + "-c",
		xvars.Bing + "-b",
		xvars.Bing + "-p",
		xvars.Bing + "-s":
		deleteConversationContext(ctx)
		c, err := createConversationContext(ctx, bot)
		if err != nil {
			ctx.Send(err.Error())
			return
		}
		cctx = c
	default:
		ctx.Send("未知的AI类型：`" + bot + "`")
		return
	}

	lmt.Remove(cctx.Bot)
	store.DeleteOnline(cctx.Id)
	ctx.Send("已切换`" + bot + "`AI模型")
}

func Run(addr string) {
	cmd.Run(addr)
	logrus.Info("已开启 `" + addr + "` Web服务")
}

// Package playwright 网页截图
package playwright

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/playwright-community/playwright-go"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/extension/single"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	githubSearchURL = "https://api.github.com/search/repositories?q=%v"
)

func init() {
	engine := control.Register("playwright", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "网页截图",
		Help: "- /网页截图 -p https://zhuanlan.zhihu.com/p/497349204 -w 600 -h 800\n" +
			"- >github xxx\n",
		PrivateDataFolder: "playwright",
	}).ApplySingle(single.New(
		single.WithKeyFn(func(ctx *zero.Ctx) int64 { return ctx.Event.GroupID }),
		single.WithPostFn[int64](func(ctx *zero.Ctx) {
			ctx.Send(
				message.ReplyWithMessage(ctx.Event.MessageID,
					message.Text("已经有正在进行的网页截图..."),
				),
			)
		}),
	))
	cachePath := engine.DataFolder() + "cache/"
	_ = os.RemoveAll(cachePath)
	_ = os.MkdirAll(cachePath, 0755)
	screenshot := func(ctx *zero.Ctx) (picPath string) {
		pageURL := ctx.State["pageURL"].(string)
		width := ctx.State["width"].(float64)
		height := ctx.State["height"].(float64)
		uid := ctx.Event.UserID
		now := time.Now()
		today := now.Format("20060102")
		pwFile := cachePath + strconv.FormatInt(uid, 10) + today + "playwright.png"
		ctx.SendChain(message.Text("少女祈祷中..."))
		pw, err := playwright.Run()
		if err != nil {
			ctx.Send(fmt.Sprintf("could not start playwright: %v", err))
			return
		}
		browser, err := pw.Chromium.Launch()
		if err != nil {
			ctx.Send(fmt.Sprintf("could not launch browser: %v", err))
			return
		}
		device := pw.Devices["Pixel 5"]
		context, err := browser.NewContext(playwright.BrowserNewContextOptions{
			Geolocation: &playwright.BrowserNewContextOptionsGeolocation{
				Longitude: playwright.Float(12.492507),
				Latitude:  playwright.Float(41.889938),
			},
			Permissions:       []string{"geolocation"},
			Viewport:          device.Viewport,
			UserAgent:         playwright.String(device.UserAgent),
			DeviceScaleFactor: playwright.Float(device.DeviceScaleFactor),
			IsMobile:          playwright.Bool(device.IsMobile),
			HasTouch:          playwright.Bool(device.HasTouch),
		})
		if err != nil {
			ctx.Send(fmt.Sprintf("could not create context: %v", err))
			return
		}
		page, err := context.NewPage()
		if err != nil {
			ctx.Send(fmt.Sprintf("could not create page: %v", err))
			return
		}
		if _, err = page.Goto(pageURL, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		}); err != nil {
			ctx.Send(fmt.Sprintf("could not goto: %v", err))
			return
		}
		x := float64(0)
		y := float64(0)
		fullpage := true
		clip := (*playwright.PageScreenshotOptionsClip)(nil)
		if width != 0 && height != 0 {
			clip = &playwright.PageScreenshotOptionsClip{
				X:      &x,
				Y:      &y,
				Width:  &width,
				Height: &height,
			}
		}
		if _, err = page.Screenshot(playwright.PageScreenshotOptions{
			Path:     playwright.String(pwFile),
			FullPage: &fullpage,
			Clip:     clip,
		}); err != nil {
			ctx.Send(fmt.Sprintf("could not create screenshot: %v", err))
			return
		}
		picPath = "file:///" + file.BOTPATH + "/" + pwFile
		if err = browser.Close(); err != nil {
			ctx.Send(fmt.Sprintf("could not close browser: %v", err))
			return
		}
		if err = pw.Stop(); err != nil {
			ctx.Send(fmt.Sprintf("could not stop Playwright: %v", err))
			return
		}
		return
	}
	engine.OnCommand("网页截图", zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		fset := flag.FlagSet{}
		var (
			pageURL string
			width   float64
			height  float64
		)
		fset.StringVar(&pageURL, "p", "https://zhuanlan.zhihu.com/p/497349204", "网页链接")
		fset.Float64Var(&width, "w", 0, "宽度")
		fset.Float64Var(&height, "h", 0, "高度")
		arguments := shell.Parse(ctx.State["args"].(string))
		err := fset.Parse(arguments)
		if err != nil {
			ctx.SendChain(message.Text("Error:", err))
			return
		}
		ctx.State["pageURL"] = pageURL
		ctx.State["width"] = width
		ctx.State["height"] = height
		ctx.SendChain(message.Image(screenshot(ctx)))
	})
	engine.OnRegex(`^>github\s?(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		data, err := web.RequestDataWith(web.NewDefaultClient(), fmt.Sprintf(githubSearchURL, url.QueryEscape(ctx.State["regex_matched"].([]string)[1])), "GET", "", web.RandUA(), nil)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
		}
		// 解析请求
		info := gjson.ParseBytes(data)
		if info.Get("total_count").Int() == 0 {
			ctx.SendChain(message.Text("ERROR: 没有找到这样的仓库"))
			return
		}
		repo := info.Get("items.0")
		ctx.State["pageURL"] = strings.ReplaceAll(repo.Get("html_url").Str, "https://github.com/", "https://starchart.cc/")
		ctx.State["width"] = 0.0
		ctx.State["height"] = 0.0
		starPicURL := screenshot(ctx)
		// 发送结果
		ctx.SendChain(
			message.Text(
				repo.Get("full_name").Str, "\n",
				"Description: ",
				repo.Get("description").Str, "\n",
				"Star/Fork/Issue: ",
				repo.Get("watchers").Int(), "/", repo.Get("forks").Int(), "/", repo.Get("open_issues").Int(), "\n",
				"Language: ",
				notnull(repo.Get("language").Str), "\n",
				"License: ",
				notnull(strings.ToUpper(repo.Get("license.key").Str)), "\n",
				"Last pushed: ",
				repo.Get("pushed_at").Str, "\n",
				"Jump: ",
				repo.Get("html_url").Str, "\n",
			),
			message.Image(
				"https://opengraph.githubassets.com/0/"+repo.Get("full_name").Str,
			).Add("cache", 0),
			message.Image(starPicURL))
	})
}

func notnull(text string) string {
	if text == "" {
		return "None"
	}
	return text
}

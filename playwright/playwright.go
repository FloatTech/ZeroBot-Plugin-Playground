// Package playwright 网页截图
package playwright

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/playwright-community/playwright-go"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/extension/single"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("playwright", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "网页截图\n- /网页截图 -p https://zhuanlan.zhihu.com/p/497349204 -w 600 -h 800\n",
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
	engine.OnCommand("网页截图").SetBlock(true).Handle(func(ctx *zero.Ctx) {
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
		ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + pwFile))
		if err = browser.Close(); err != nil {
			ctx.Send(fmt.Sprintf("could not close browser: %v", err))
			return
		}
		if err = pw.Stop(); err != nil {
			ctx.Send(fmt.Sprintf("could not stop Playwright: %v", err))
			return
		}
	})
}

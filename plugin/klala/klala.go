// Package klala 星穹铁道图鉴查询
package klala

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

func init() { // 主函数
	en := control.Register("klala", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "星穹铁道图鉴查询",
		Help: "- 更新图鉴\n" +
			"- *xx图鉴",
		PrivateDataFolder: "klala",
	})
	en.OnRegex(`^\*(.*)图鉴$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if file.IsNotExist(en.DataFolder() + "star-rail-atlas") {
			ctx.SendChain(message.Text("请先发送\"更新图鉴\"!"))
			return
		}
		word := ctx.State["regex_matched"].([]string)[1] // 关键字
		if word == "" {
			return
		}
		t, err := os.ReadFile(en.DataFolder() + "star-rail-atlas/path.json") // 获取文件
		if err != nil {
			ctx.SendChain(message.Text("获取路径文件失败", err))
			return
		}
		var paths wikimap
		_ = json.Unmarshal(t, &paths)
		var path string
		var ok bool
		path, ok = paths.Light[word]
		if !ok {
			path, ok = paths.Role[word]
			if !ok {
				ctx.SendChain(message.Text("-未找到图鉴"))
				return
			}
		}
		ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + en.DataFolder() + "star-rail-atlas" + path))
	})
	en.OnRegex(`^更新图鉴$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var cmd *exec.Cmd
		if file.IsNotExist(en.DataFolder() + "star-rail-atlas") {
			cmd = exec.Command("git", "clone", "https://github.com/Nwflower/star-rail-atlas.git")
		} else {
			cmd = exec.Command("git", "pull")
		}
		cmd.Dir = file.BOTPATH + "/" + en.DataFolder()
		output, err := cmd.CombinedOutput()
		if err != nil {
			ctx.SendChain(message.Text("运行失败: ", err, "\n", helper.BytesToString(output)))
			return
		}
		ctx.SendChain(message.Text("运行成功: ", helper.BytesToString(output)))
	})
}

type wikimap struct {
	Light map[string]string `json:"lightcone"`
	Role  map[string]string `json:"role"`
}

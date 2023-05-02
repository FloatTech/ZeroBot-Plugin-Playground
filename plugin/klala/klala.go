// Package klala 星穹铁道图鉴查询
package klala

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

const (
	gitteURL = "https://gitee.com/Nwflower/star-rail-atlas.git"
	//githubURL = "https://github.com/Nwflower/star-rail-atlas.git"
)

func init() { // 主函数
	en := control.Register("klala", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "星穹铁道图鉴查询",
		Help: "- *(强制)更新图鉴\n" +
			"- *图鉴列表\n" +
			"- *xx图鉴\n" +
			"- *xx材料|素材",
		PrivateDataFolder: "klala",
	})
	en.OnRegex(`^\*(.*)(材料|素材|图鉴)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
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
		//匹配类型
		var path string
		var ok bool
		switch ctx.State["regex_matched"].([]string)[2] {
		case "材料", "素材":
			path, ok = paths.findhow(word)
		case "图鉴":
			path, ok = paths.findwhat(word)
		}
		if !ok {
			ctx.SendChain(message.Text("未找到该", ctx.State["regex_matched"].([]string)[2]))
			return
		}
		ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + en.DataFolder() + "star-rail-atlas" + path))
	})
	en.OnRegex(`^*(强制)?更新图鉴$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var cmd *exec.Cmd
		var p = file.BOTPATH + "/" + en.DataFolder()
		if ctx.State["regex_matched"].([]string)[1] != "" {
			if err := os.RemoveAll(p + "star-rail-atlas"); err != nil {
				ctx.SendChain(message.Text("-删除失败", err))
				return
			}
		}
		if file.IsNotExist(en.DataFolder() + "star-rail-atlas") {
			cmd = exec.Command("git", "clone", gitteURL)
			cmd.Dir = p
		} else {
			cmd = exec.Command("git", "pull")
			cmd.Dir = p + "star-rail-atlas/"
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			ctx.SendChain(message.Text("运行失败: ", err, "\n", helper.BytesToString(output)))
			return
		}
		ctx.SendChain(message.Text("运行成功: ", helper.BytesToString(output)))
	})
	en.OnRegex(`^*图鉴列表$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if file.IsNotExist(en.DataFolder() + "star-rail-atlas") {
			ctx.SendChain(message.Text("请先发送\"更新图鉴\"!"))
			return
		}
		index := []string{"role.yaml", "lightcone.yaml", "material for role.yaml"}
		var t [3][]byte
		var err error
		for i := 0; i < 3; i++ {
			t[i], err = os.ReadFile(en.DataFolder() + "star-rail-atlas/index/" + index[i])
			if err != nil {
				ctx.SendChain(message.Text("获取路径文件失败", err))
				return
			}
		}
		ctx.Send(message.Message{
			ctxext.FakeSenderForwardNode(ctx, message.Text(string(t[0]))),
			ctxext.FakeSenderForwardNode(ctx, message.Text(string(t[1]))),
			ctxext.FakeSenderForwardNode(ctx, message.Text(string(t[2]))),
		})
	})
}

// 寻找图鉴
func (paths wikimap) findwhat(word string) (path string, ok bool) {
	if path, ok = paths.Role[word]; ok {
		return
	}
	if path, ok = paths.Light[word]; ok {
		return
	}
	return
}
func (paths wikimap) findhow(word string) (path string, ok bool) {
	if path, ok = paths.Material[word]; ok {
		return
	}
	return
}

type wikimap struct {
	Light    map[string]string `json:"lightcone"`
	Role     map[string]string `json:"role"`
	Material map[string]string `json:"material for role"`
}

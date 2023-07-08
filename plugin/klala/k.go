// Package klala 星穹铁道
package klala

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var (
	lastExecutionTime int64
	cds               int64 = 5
	initdata                = ctxext.DoOnceOnSuccess(downdata)
	cryptic           string
	preFix            = `^[*＊]`
)

func init() { // 主函数
	en := control.Register("kkk", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "星穹铁道面板查询",
		Help: "- *xx面板\n" +
			"- *更新面板\n" +
			"- *绑定xxx\n" +
			"- *设置CD为xs",
	})
	en.OnRegex(preFix+`(.*)面板$`, initdata).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		wife := getWifeOrWq()
		currentTime := time.Now().Unix()
		uid := strconv.Itoa(getuid(strconv.FormatInt(ctx.Event.UserID, 10)))
		if uid == "0" {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("-未绑定uid\n-第一次使用请发送\"*绑定xxx\""))
			return
		}
		key := ctx.State["regex_matched"].([]string)[1]
		if key == "" {
			return
		}
		if key == "更新" {
			if currentTime-lastExecutionTime < cds {
				ctx.SendChain(message.Text("-全局时间冷却中,剩余时间", cds-currentTime+lastExecutionTime, "s"))
				return
			}
			lastExecutionTime = currentTime
			msg, err := saveRoel(uid)
			if err != nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text(err))
				return
			}
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("\n", msg))
			return
		}
		wifeid := wife.findnames("wife", key)
		key = wife.idmap("wife", wifeid)
		if key == "" {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("-请输入角色全名"))
			return
		}
		data, err := os.ReadFile(jsPath + uid + ".klala")
		if err != nil {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("-未找到本地缓存数据,请\"*更新面板\""))
			return
		}
		var t thisdata
		err = json.Unmarshal(data, &t)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		// 获取角色序列
		var n = -1
		// 匹配角色
		for i, v := range t.RoleData {
			if key == v.Name {
				n = i
				break
			}
		}
		if n == -1 { // 在返回数据中未找到想要的角色
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("-该角色未展示"))
			return
		}
		imagePath, err := t.drawcard(n)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + imagePath)) // 输出
	})

	en.OnRegex(preFix+`绑定(?:星铁)?[Uu]?[Ii]?[Dd]?\s*(\d+)`, initdata).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		currentTime := time.Now().Unix()
		suid := ctx.State["regex_matched"].([]string)[1] // 获取uid
		int64uid, err := strconv.ParseInt(suid, 10, 64)
		if suid == "" || int64uid < 100000000 || int64uid > 1000000000 || err != nil {
			// ctx.SendChain(message.Text("-请输入正确的uid"))
			return
		}
		sqquid := strconv.Itoa(int(ctx.Event.UserID))
		file, _ := os.OpenFile(uidPath+sqquid+".klala", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = file.WriteString(suid)
		file.Close()
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("-绑定uid", suid, "成功,尝试获取角色信息"))
		if currentTime-lastExecutionTime < cds {
			ctx.SendChain(message.Text("-全局时间冷却中,剩余时间", cds-currentTime+lastExecutionTime, "s"))
			return
		}
		lastExecutionTime = currentTime
		msg, err := saveRoel(suid)
		if err != nil {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text(err))
			return
		}
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("\n", msg))
	})
	en.OnRegex(preFix+`设置CD为(\d+)s`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		cs := ctx.State["regex_matched"].([]string)[1] // 获取uid
		c, _ := strconv.ParseInt(cs, 10, 64)
		if c < 5 {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("-CD太短惹"))
			return
		}
		cds = c
		ctx.SendChain(message.Text("-设置CD为", cs, "S"))
	})
	en.OnRegex(preFix+`(强制)?更新klala$`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var cmd *exec.Cmd
		var p = file.BOTPATH + "/data/klala/"
		if ctx.State["regex_matched"].([]string)[1] != "" {
			if err := os.RemoveAll(p + "kkk"); err != nil {
				ctx.SendChain(message.Text("-删除失败", err))
				return
			}
		}
		if file.IsNotExist("data/klala/kkk") {
			cmd = exec.Command("git", "clone", "https://gitee.com/lianhong2758/kkk.git")
			cmd.Dir = p
		} else {
			cmd = exec.Command("git", "pull")
			cmd.Dir = p + "kkk"
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			ctx.SendChain(message.Text("运行失败: ", err, "\n", helper.BytesToString(output)))
			return
		}
		o := helper.BytesToString(output)
		if len(o) > 500 {
			o = o[:500] + "\n..."
		}
		ctx.SendChain(message.Text("运行成功: ", o))
	})
}

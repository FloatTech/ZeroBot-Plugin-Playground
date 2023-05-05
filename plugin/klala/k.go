// Package klala 星穹铁道面板查询
package klala

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/file"
	img "github.com/FloatTech/imgfactory"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var (
	lastExecutionTime int64
	cds               int64 = 150
	initdata                = ctxext.DoOnceOnSuccess(downdata)
)

func init() { // 主函数
	en := control.Register("kkk", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "星穹铁道查询",
		Help: "- *xx面板\n" +
			"- *更新面板\n" +
			"- *绑定xxx",
	})
	en.OnRegex(`^\*(.*)面板$`, initdata).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		wife := GetWifeOrWq("wife")
		currentTime := time.Now().Unix()
		key := ctx.State["regex_matched"].([]string)[1]
		uid := strconv.Itoa(Getuid(strconv.FormatInt(ctx.Event.UserID, 10)))
		if uid == "0" {
			ctx.SendChain(message.Text("-未绑定uid\n-第一次使用请发送\"*绑定xxx\""))
			return
		}
		if key == "" {
			return
		} else if key == "更新" {
			if currentTime-lastExecutionTime < cds {
				ctx.SendChain(message.Text("-全局时间冷却中,剩余时间", cds-currentTime+lastExecutionTime, "s"))
				return
			}
			lastExecutionTime = currentTime
			_, _ = Updata(uid)
			listdata, err := CaseList(uid)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			var list roleList
			err = json.Unmarshal(listdata, &list)
			if err != nil {
				ctx.SendChain(message.Text("-请将想要查看的角色放在展框中,并等待5分钟之后尝试"))
				return
			}
			var msg strings.Builder
			msg.WriteString("-更新成功,您展示的角色为: ")
			for _, v := range list.Data.Characters {
				d, err := GetRole(uid, strconv.Itoa(v.ID))
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				var r roles
				err = json.Unmarshal(d, &r)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				thisdata, err := r.convertData(list.Data.NickName, v)
				if err != nil {
					ctx.SendChain(message.Text("数据映射错误捏：", err))
					return
				}
				es, err := json.Marshal(&thisdata)
				if err != nil {
					ctx.SendChain(message.Text("数据反解析错误捏：", err))
					return
				}
				//_ = downcard(v.DrawcardIcon, strconv.Itoa(v.ID))
				file, _ := os.OpenFile("data/klala/kkk/js/"+uid+strconv.Itoa(v.ID)+".klala", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				_, _ = file.Write(es)
				file.Close()
				msg.WriteString("\n")
				msg.WriteString(wife.Idmap(strconv.Itoa(v.ID)))
			}
			ctx.SendChain(message.Text(msg.String()))
			return
		}
		wifeid := wife.Findnames(key)
		key = wife.Idmap(wifeid)
		if key == "" {
			ctx.SendChain(message.Text("-请输入角色全名"))
			return
		}
		data, err := os.ReadFile("data/klala/kkk/js/" + uid + wifeid + ".klala")
		if err != nil {
			ctx.SendChain(message.Text("-未找到角色数据"))
			return
		}
		var t thisdata
		err = json.Unmarshal(data, &t)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		drawimage, err := t.drawcard()
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ff, err := img.ToBytes(drawimage) // 图片放入缓存
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(message.ImageBytes(ff)) // 输出
	})

	en.OnRegex(`^\*绑定(\d+)`, initdata).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		suid := ctx.State["regex_matched"].([]string)[1] // 获取uid
		int64uid, err := strconv.ParseInt(suid, 10, 64)
		if suid == "" || int64uid < 100000000 || int64uid > 1000000000 || err != nil {
			//ctx.SendChain(message.Text("-请输入正确的uid"))
			return
		}
		sqquid := strconv.Itoa(int(ctx.Event.UserID))
		file, _ := os.OpenFile("data/klala/kkk/uid/"+sqquid+".klala", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = file.WriteString(suid)
		file.Close()
		ctx.SendChain(message.Text("-绑定成功"))
	})
	en.OnRegex(`^\*设置CD为(\d+)s`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		cs := ctx.State["regex_matched"].([]string)[1] // 获取uid
		c, _ := strconv.ParseInt(cs, 10, 64)
		if c < 60 {
			ctx.SendChain(message.Text("-CD太短惹"))
			return
		}
		cds = c
		ctx.SendChain(message.Text("-设置CD为", cs, "S"))
	})
	en.OnRegex(`^*更新Klala$`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
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
		ctx.SendChain(message.Text("运行成功: ", helper.BytesToString(output)))
	})
}

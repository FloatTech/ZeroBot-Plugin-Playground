package allmessage

import (
    "strconv"
    "strings"

    ctrl "github.com/FloatTech/zbpctrl"
    "github.com/FloatTech/zbputils/control"
    zero "github.com/wdvxdr1123/ZeroBot"
    "github.com/wdvxdr1123/ZeroBot/message"
)

var blacklist = map[int64]bool{}

func init() {
    engine := control.AutoRegister(&ctrl.Options[*zero.Ctx]{
        DisableOnDefault:  false,
        Brief:             "群广播-发送消息到所有群聊喵",
        Help:              "此插件仅拱Bot管理员使用；\n支持？黑名单？添加/移除/查看\n用法:\n广播+内容\n广播黑名单添加 123456789\n广播黑名单移除 123456789\n广播黑名单列表",
        PrivateDataFolder: "allmessage",
    })

    engine.OnPrefix("广播").SetBlock(true).Handle(func(ctx *zero.Ctx) {
        if !zero.SuperUserPermission(ctx) {
            ctx.Send("喵～你不是主人，我不能帮你广播哦！")
            return
        }

        raw := strings.TrimSpace(strings.TrimPrefix(ctx.Event.RawMessage, "广播"))
        if raw == "" {
            ctx.Send("呜呜～你没有发内容，广播取消啦！")
            return
        }

        // 黑名单管理指令
        if strings.HasPrefix(raw, "黑名单添加 ") {
            idStr := strings.TrimPrefix(raw, "黑名单添加 ")
            gid, err := strconv.ParseInt(idStr, 10, 64)
            if err != nil {
                ctx.Send("喵～群号格式不对喵！")
                return
            }
            blacklist[gid] = true
            ctx.Send("已将群 " + idStr + " 加入黑名单喵！")
            return
        }

        if strings.HasPrefix(raw, "黑名单移除 ") {
            idStr := strings.TrimPrefix(raw, "黑名单移除 ")
            gid, err := strconv.ParseInt(idStr, 10, 64)
            if err != nil {
                ctx.Send("喵～群号格式不对喵！")
                return
            }
            delete(blacklist, gid)
            ctx.Send("已将群 " + idStr + " 移出黑名单喵！")
            return
        }

        if raw == "黑名单列表" {
            if len(blacklist) == 0 {
                ctx.Send("当前黑名单为空喵～")
                return
            }
            var list []string
            for gid := range blacklist {
                list = append(list, strconv.FormatInt(gid, 10))
            }
            ctx.Send("当前黑名单群号喵：\n" + strings.Join(list, "\n"))
            return
        }

        // 正常广播内容
        var msg message.Message
        if strings.HasPrefix(raw, "图片 ") {
            url := strings.TrimPrefix(raw, "图片 ")
            msg = message.Message{message.Image(url)}
        } else if strings.HasPrefix(raw, "语音 ") {
            url := strings.TrimPrefix(raw, "语音 ")
            msg = message.Message{message.Record(url)}
        } else {
            msg = message.Message{message.Text(raw)}
        }

        groupList := ctx.GetGroupList()
        if !groupList.Exists() {
            ctx.Send("呜呜～我现在拿不到群列表，可能没有权限喵！")
            return
        }

        var success, fail, skipped int
        for _, g := range groupList.Array() {
            gidField := g.Get("group_id")
            if !gidField.Exists() {
                continue
            }
            gid := gidField.Int()
            if blacklist[gid] {
                skipped++
                continue
            }
            if ctx.SendGroupMessage(gid, msg) == 0 {
                fail++
            } else {
                success++
            }
        }

        ctx.Send(message.Text("广播完成啦喵～成功 ", success, " 个群，失败 ", fail, " 个群，跳过黑名单 ", skipped, " 个群！"))
    })
}

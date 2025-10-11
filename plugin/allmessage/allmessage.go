package allmessage

import (
    "strings"

    ctrl "github.com/FloatTech/zbpctrl"
    "github.com/FloatTech/zbputils/control"
    zero "github.com/wdvxdr1123/ZeroBot"
    "github.com/wdvxdr1123/ZeroBot/message"
)

var blacklist = map[int64]bool{
    123456789: true, // 示例黑名单群号
}

func init() {
    engine := control.AutoRegister(&ctrl.Options[*zero.Ctx]{
        DisableOnDefault:  false,
        Brief:             "广播",
        Help:              "bot管理员发送“广播+内容”即可立即广播到所有群（跳过黑名单）",
        PrivateDataFolder: "allmessage",
    })

    engine.OnPrefix("广播 ").SetBlock(true).Handle(func(ctx *zero.Ctx) {
        if !zero.SuperUserPermission(ctx) {
            ctx.Send("喵～你不是主人，我不能帮你广播哦！")
            return
        }

        raw := strings.TrimSpace(strings.TrimPrefix(ctx.Event.RawMessage, "广播"))
        if raw == "" {
            ctx.Send("呜呜～你没有发内容，广播取消啦！")
            return
        }

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

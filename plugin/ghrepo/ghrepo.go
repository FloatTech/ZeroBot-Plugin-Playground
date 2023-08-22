package github

import (
	b64 "encoding/base64"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	repoAPI  = "https://api.github.com/repos/"
	imageAPI = "https://opengraph.githubassets.com/1a/"
)

func init() {
	zero.OnPrefix("https://github.com/").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			repoInfo := ctx.Event.Message.ExtractPlainText()[19:]
			// 去除域名后内容过短，忽略消息
			if len(repoInfo) <= 2 {
				return
			}
			// 删除末尾的 /
			if repoInfo[len(repoInfo)-1] == '/' {
				repoInfo = repoInfo[:len(repoInfo)-1]
			}
			// 检查是否同时含有用户名和仓库名
			repoInfoSlice := strings.Split(repoInfo, "/")
			if len(repoInfoSlice) != 2 {
				return
			}
			// 检查仓库是否存在
			data, err := web.GetData(repoAPI + repoInfo)
			if err != nil {
				log.Errorln("[github]", "Fail to check repo status", err)
				return
			}
			repoStatusMessage := gjson.Get(binary.BytesToString(data), "message").String()
			// 仓库不存在，忽略
			if repoStatusMessage == "Not Found" {
				return
			}
			// 下载仓库图片
			imageData, err := web.GetData(imageAPI + repoInfo)
			if err != nil {
				log.Errorln("[github]", "Fail to download repo image", err)
				return
			}
			// 发送仓库图片
			repoText := "https://github.com/" + repoInfo
			imageB64 := "base64://" + b64.StdEncoding.EncodeToString(imageData)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(repoText), message.Image(imageB64))
		})
}

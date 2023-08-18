// Package grasspicture 草图
package grasspicture

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	grassURL  = "https://oss.grass.starxw.com"
	infoURL   = grassURL + "/service/info"
	imageURL  = grassURL + "/service/image"
	statusURL = grassURL + "/service/status"
	uploadURL = grassURL + "/service/upload"
)

type infoRsp struct {
	Code      int    `json:"code"`
	ID        string `json:"id"`
	LikeCount int    `json:"likeCount"`
}

type statusRsp struct {
	TotalImage          int    `json:"totalImage"`
	TotalImageSize      int    `json:"totalImageSize"`
	TotalImageSizeHuman string `json:"totalImageSizeHuman"`
	WaitImage           string `json:"waitImage"`
	APICountToday       int    `json:"apiCountToday"`
	APIFlowToday        int    `json:"apiFlowToday"`
	APIFlowTodayHuman   string `json:"apiFlowTodayHuman"`
	Service             bool   `json:"service"`
}

type uploadRsp struct {
	Code int `json:"code"`
}

func init() {
	engine := control.Register("grasspicture", &ctrl.Options[*zero.Ctx]{
		PrivateDataFolder: "grasspicture",
		DisableOnDefault:  false,
		Brief:             "草图",
		Help:              "- 生草 | 来张草图\n- 草图信息\n- 草图投稿 (来个大佬解决问题)\n开发文档: https://apifox.com/apidoc/shared-8a633395-412a-460a-87e0-b82012573873",
	})
	cachePath := engine.DataFolder() + "cache/"
	_ = os.RemoveAll(cachePath)
	_ = os.MkdirAll(cachePath, 0755)
	engine.OnFullMatchGroup([]string{"生草", "来张草图"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		data, err := web.GetData(infoURL)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		var info infoRsp
		err = json.Unmarshal(data, &info)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		data, err = web.GetData(fmt.Sprintf("%v?id=%v&type=image", imageURL, info.ID))
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(message.ImageBytes(data))
	})
	engine.OnFullMatchGroup([]string{"草图信息"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		data, err := web.GetData(statusURL)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		var status statusRsp
		err = json.Unmarshal(data, &status)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		text := fmt.Sprintf("---==草图信息==---\n是否正常提供服务: %v\n图片总数: %v\n待审核图片数: %v\n今日调用次数: %v\n图片总大小: %v\n今日图片流量: %v\n", status.Service, status.TotalImage, status.WaitImage, status.APICountToday, status.TotalImageSizeHuman, status.APIFlowTodayHuman)
		ctx.SendChain(message.Text(text))
	})
	engine.OnMessage(zero.AdminPermission, func(ctx *zero.Ctx) bool {
		msg := ctx.Event.Message
		if len(msg) < 1 {
			return false
		}
		if msg[0].Type != "reply" {
			return false
		}
		if !strings.Contains(msg.ExtractPlainText(), "草图投稿") {
			return false
		}
		return true
	}).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			id := ctx.Event.Message[0].Data["id"]
			imageURL := ctx.GetMessage(message.NewMessageIDFromString(id)).Elements[0].Data["url"]
			fileName := cachePath + id + ".jpg"
			err := file.DownloadTo(imageURL, fileName)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			var payload bytes.Buffer
			writer := multipart.NewWriter(&payload)
			file, err := os.Open(fileName)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			defer file.Close()
			part, err := writer.CreateFormFile("file", fileName)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			_, err = io.Copy(part, file)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			err = writer.Close()
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			data, err := web.PostData(uploadURL, writer.FormDataContentType(), &payload)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			var upload uploadRsp
			err = json.Unmarshal(data, &upload)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			var text string
			switch upload.Code {
			case 200:
				text = "草图投稿成功"
			case 400:
				text = "您已被限流"
			case 403:
				text = "您已被加入黑名单"
			case 1000:
				text = "没有发送文件，或发送的文件非有效图片格式"
			case 1001:
				text = "文件体积大于 2MB"
			}
			ctx.SendChain(message.ReplyWithMessage(ctx.Event.Message[0].Data["id"], message.Text(text))...)
		})
}

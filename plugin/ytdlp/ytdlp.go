// Package ytdlp YouTube 视频下载
package ytdlp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	en := control.Register("ytdlp", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "ytdlp 下载器",
		Help:              "/yt-dlp [youtube-video-link]",
		PrivateDataFolder: "ytdlp",
	}).ApplySingle(ctxext.DefaultSingle)
	// 初始化临时文件夹
	tempFileDir := path.Join(en.DataFolder(), "temp")
	err := os.MkdirAll(tempFileDir, 0750)
	if err != nil {
		panic(err)
	}
	// 注册指令
	en.OnPrefix("/yt-dlp", zero.OnlyGroup, zero.SuperUserPermission).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			wd, err := os.Getwd()
			if err != nil {
				ctx.Send("系统错误，无法获取当前程序目录。")
			}
			_, err = exec.LookPath("yt-dlp")
			if err != nil {
				ctx.Send("目标服务器未安装 yt-dlp。")
				return
			}
			url := ctx.State["args"].(string)
			cmdVideoFileSize := exec.Command("yt-dlp", "--print", "%(filesize,filesize_approx)s", url)
			videoFileSizeInByte, err := cmdVideoFileSize.Output()
			if err != nil {
				ctx.Send("获取视频文件大小失败，可能是服务器无法访问 YouTube。")
				return
			}
			videoFileSize := string(videoFileSizeInByte)
			cmdVideoTitle := exec.Command("yt-dlp", "--print", "%(title)s", url)
			videoTitleInByte, err := cmdVideoTitle.Output()
			if err != nil {
				ctx.Send("获取视频标题失败，可能是服务器无法访问 YouTube。")
				return
			}
			videoTitle := string(videoTitleInByte)
			ctx.Send(fmt.Sprintf("视频标题：%s视频大小：%s即将开始下载视频，请稍候。", videoTitle, videoFileSize))
			fileName := fmt.Sprintf("%d_%d.mp4", ctx.Event.Sender.ID, time.Now().Unix())
			videoFilePath := path.Join(tempFileDir, fileName)
			cmdDownload := exec.Command("yt-dlp", url, "-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best", "-o", videoFilePath)
			downloadLogInByte, err := cmdDownload.Output()
			downloadLog := string(downloadLogInByte)
			if err != nil {
				ctx.Send(fmt.Sprintf("视频文件下载失败，可能存在网络波动。\n%s", downloadLog))
				return
			}
			ctx.Send("文件下载成功，正在上传，请稍候。")
			fullVideoFilePath := path.Join(wd, videoFilePath)
			resp := ctx.UploadThisGroupFile(fullVideoFilePath, fileName, "")
			if resp.Status != "ok" {
				respData, err := json.Marshal(resp)
				errMsg := string(respData)
				if err != nil {
					errMsg = err.Error()
				}
				ctx.Send(fmt.Sprintln("文件上传失败。", errMsg))
			}
			os.Remove(videoFilePath)
		})
}

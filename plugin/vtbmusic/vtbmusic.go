// Package vtbmusic vtb点歌
package vtbmusic

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	getGroupListURL = "https://aqua.chat/v1/GetGroupsList"
	getMusicListURL = "https://aqua.chat/v1/GetMusicList"
)

func init() { // 插件主体
	engine := control.Register("vtbmusic", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "vtb点歌",
		Help:              "- vtb点歌",
		PrivateDataFolder: "vtbmusic",
	})

	// 开启
	engine.OnFullMatch(`vtb点歌`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {

		})
}

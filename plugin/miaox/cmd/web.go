package cmd

import (
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/types"
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/vars"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	r     = gin.Default()
	menus []Menu
)

type Menu struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Model string `json:"model"`
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func NewMenu(model string, name string) Menu {
	return Menu{
		Name:  name,
		Model: model,
		Path:  "/" + model,
	}
}

func Register(base string, sv types.ModelService, m Menu) {
	menus = append(menus, m)

	r.GET(base+"/get", func(ctx *gin.Context) {
		id := ctx.DefaultQuery("id", "-1")
		if id == "-1" {
			ctx.JSON(200, &types.JsonRest{
				Code: 500,
				Msg:  "`id`的值不正确",
			})
			return
		}
		model := sv.Get(id)
		ctx.JSON(200, &types.JsonRest{
			Code: 200,
			Data: model,
		})
	})

	r.DELETE(base+"/del", func(ctx *gin.Context) {
		id := ctx.DefaultQuery("id", "-1")
		if id == "-1" {
			ctx.JSON(200, &types.JsonRest{
				Code: 500,
				Msg:  "`id`的值不正确",
			})
			return
		}

		if result := sv.Del(id); !result {
			ctx.JSON(200, &types.JsonRest{
				Code: 500,
				Msg:  "删除失败",
			})
		} else {
			ctx.JSON(200, &types.JsonRest{
				Code: 200,
			})
		}
	})

	r.POST(base+"/edit", func(ctx *gin.Context) {
		model := sv.NewModel()
		if err := ctx.BindJSON(&model); err != nil {
			ctx.JSON(500, err)
			return
		}
		if result := sv.Edit(model); !result {
			ctx.JSON(200, types.JsonRest{
				Code: 500,
				Msg:  "执行失败",
			})
		} else {
			ctx.JSON(200, types.JsonRest{
				Code: 200,
			})
		}
	})

	r.POST(base+"/page", func(ctx *gin.Context) {
		model := sv.NewModel()
		if err := ctx.BindJSON(&model); err != nil {
			ctx.JSON(500, err)
			return
		}
		page := sv.Find(model)
		ctx.JSON(200, types.JsonRest{
			Code: 200,
			Data: page,
		})
	})
}

func Run(addr string) {
	r.Static("/model", vars.E.DataFolder()+"/static/model")
	r.Static("/assets", vars.E.DataFolder()+"/static/assets")
	r.StaticFile("/", vars.E.DataFolder()+"/static/index.html")
	r.StaticFile("/favicon.ico", vars.E.DataFolder()+"/static/favicon.ico")
	r.GET("/api/menu", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code": 200,
			"data": menus,
		})
	})
	r.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/")
	})
	go func() {
		err := r.Run(addr)
		if err != nil {
			logrus.Error(err)
			os.Exit(0)
		}
	}()
}

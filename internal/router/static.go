package router

import (
	"TTPanel/internal/global"
	"github.com/gin-gonic/gin"
)

type StaticRouter struct{}

func (s *StaticRouter) Init(Router *gin.Engine) {
	//前端面板静态文件路由
	{
		Router.LoadHTMLGlob(global.Config.System.PanelPath + "/Templates/*.html") // npm打包成dist的路径
		Router.StaticFile("/favicon.ico", global.Config.System.PanelPath+"/Templates/favicon.ico")
		Router.Static("/static", global.Config.System.PanelPath+"/Templates/static")   // dist里面的静态资源
		Router.Static("/ssh_cast", global.Config.System.PanelPath+"/data/cast")        // cast的静态资源
		Router.StaticFile("/", global.Config.System.PanelPath+"/Templates/index.html") // 前端网页入口页面
	}
}

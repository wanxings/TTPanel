package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type HostRouter struct{}

func (s *HostRouter) Init(Router *gin.RouterGroup) {
	hostRouter := Router.Group("host")
	hostApi := api.GroupApp.HostApiApp
	{
		hostRouter.POST("AddHostCategory", hostApi.AddHostCategory)       // 添加主机分类
		hostRouter.POST("HostCategoryList", hostApi.HostCategoryList)     // 主机分类列表
		hostRouter.POST("EditHostCategory", hostApi.EditHostCategory)     // 编辑主机分类
		hostRouter.POST("DeleteHostCategory", hostApi.DeleteHostCategory) // 删除主机分类
	}
	{
		hostRouter.POST("AddHost", hostApi.AddHost)       // 添加主机
		hostRouter.POST("HostList", hostApi.HostList)     // 主机列表
		hostRouter.POST("DeleteHost", hostApi.DeleteHost) // 删除主机
	}
	{
		hostRouter.GET("Terminal", hostApi.Terminal) // 终端
	}
	{
		hostRouter.POST("ShortcutCommandList", hostApi.ShortcutCommandList)     // 快捷命令列表
		hostRouter.POST("AddShortcutCommand", hostApi.AddShortcutCommand)       // 添加快捷命令
		hostRouter.POST("DeleteShortcutCommand", hostApi.DeleteShortcutCommand) // 删除快捷命令
		hostRouter.POST("EditShortcutCommand", hostApi.EditShortcutCommand)     // 编辑快捷命令
	}

}

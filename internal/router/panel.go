package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type PanelRouter struct{}

func (s *PanelRouter) Init(Router *gin.RouterGroup) {
	panelRouter := Router.Group("panel")
	panelApi := api.GroupApp.PanelApiApp
	{
		panelRouter.POST("ExtensionList", panelApi.ExtensionList) // 扩展列表
		panelRouter.POST("Base", panelApi.Base)                   // 面板基础信息
		panelRouter.POST("OperatePanel", panelApi.OperatePanel)   // 操作面板
		panelRouter.POST("OperateServer", panelApi.OperateServer) // 操作服务器
		panelRouter.POST("CheckUpdate", panelApi.CheckUpdate)     // 检查更新
		panelRouter.POST("Update", panelApi.Update)               // 更新
	}
}

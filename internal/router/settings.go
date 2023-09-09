package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type SettingsRouter struct{}

func (s *SettingsRouter) Init(Router *gin.RouterGroup) {
	settingsRouter := Router.Group("settings")
	settingsApi := api.GroupApp.SettingsApiApp
	{
		settingsRouter.POST("List", settingsApi.List)                                             //获取系统设置
		settingsRouter.POST("SetBasicAuth", settingsApi.SetBasicAuth)                             //设置基础认证
		settingsRouter.POST("SetPanelPort", settingsApi.SetPanelPort)                             //设置面板端口
		settingsRouter.POST("SetEntrance", settingsApi.SetEntrance)                               //设置面板入口
		settingsRouter.POST("SetEntranceErrorCode", settingsApi.SetEntranceErrorCode)             //设置面板入口错误码
		settingsRouter.POST("SetUser", settingsApi.SetUser)                                       //设置用户
		settingsRouter.POST("SetPanelName", settingsApi.SetPanelName)                             // 设置面板名称
		settingsRouter.POST("SetPanelIP", settingsApi.SetPanelIP)                                 // 设置面板IP
		settingsRouter.POST("SetDefaultWebsiteDirectory", settingsApi.SetDefaultWebsiteDirectory) // 设置默认网站目录
		settingsRouter.POST("SetDefaultBackupDirectory", settingsApi.SetDefaultBackupDirectory)   // 设置默认备份目录
		settingsRouter.POST("SetPanelApi", settingsApi.SetPanelApi)                               // 设置面板API
		settingsRouter.POST("SetLanguage", settingsApi.SetLanguage)                               // 设置面板语言
	}
}

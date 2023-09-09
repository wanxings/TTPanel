package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type MonitorRouter struct{}

func (s *MonitorRouter) Init(Router *gin.RouterGroup) {
	monitorRouter := Router.Group("monitor")
	monitorApi := api.GroupApp.MonitorApiApp
	{
		monitorRouter.POST("Base", monitorApi.Base)                 // 获取系统基本信息
		monitorRouter.POST("Logs", monitorApi.Logs)                 // 获取日志
		monitorRouter.POST("ClearAllLogs", monitorApi.ClearAllLogs) // 清空日志
		monitorRouter.POST("Config", monitorApi.Config)             // 获取配置
		monitorRouter.POST("SaveConfig", monitorApi.SaveConfig)     // 保存配置
	}
	{
		monitorRouter.POST("EventConfig", monitorApi.EventConfig)                 // 事件配置项
		monitorRouter.POST("SaveEventConfig", monitorApi.SaveEventConfig)         // 保存事件配置项
		monitorRouter.POST("EventList", monitorApi.EventList)                     // 事件列表
		monitorRouter.POST("BatchSetEventStatus", monitorApi.BatchSetEventStatus) // 批量设置事件状态
	}
}

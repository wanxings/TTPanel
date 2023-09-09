package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type LogAuditRouter struct{}

func (s *LogAuditRouter) Init(Router *gin.RouterGroup) {
	logAuditRouter := Router.Group("log_audit")
	logAuditApi := api.GroupApp.LogAuditApiApp
	{
		logAuditRouter.POST("PanelOperationLogList", logAuditApi.PanelOperationLogList)   // 面板操作日志列表
		logAuditRouter.POST("ClearPanelOperationLog", logAuditApi.ClearPanelOperationLog) // 清空面板操作日志
	}
	{
		logAuditRouter.POST("LogFileOccupancy", logAuditApi.LogFileOccupancy) // 日志占用
	}
	{
		logAuditRouter.POST("SSHLoginLogList", logAuditApi.SSHLoginLogList) // ssh登录日志
	}

}

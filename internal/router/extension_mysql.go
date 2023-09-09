package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ExtensionMysqlRouter struct{}

func (s *ExtensionMysqlRouter) Init(Router *gin.RouterGroup) {
	mysqlR := Router.Group("extension_mysql")
	mysqlApi := api.GroupApp.ExtensionMysqlApiApp
	{
		mysqlR.POST("Info", mysqlApi.Info)                                   // 获取Mysql的信息
		mysqlR.POST("Install", mysqlApi.Install)                             // 安装Mysql
		mysqlR.POST("Uninstall", mysqlApi.Uninstall)                         // 卸载Mysql
		mysqlR.POST("SetStatus", mysqlApi.SetStatus)                         // 设置Mysql状态
		mysqlR.POST("PerformanceConfig", mysqlApi.PerformanceConfig)         // 获取性能配置
		mysqlR.POST("SavePerformanceConfig", mysqlApi.SavePerformanceConfig) // 获取性能配置
		mysqlR.POST("LoadStatus", mysqlApi.LoadStatus)                       // 负载状态
		mysqlR.POST("ErrorLog", mysqlApi.ErrorLog)                           // 错误日志
		mysqlR.POST("SlowLogs", mysqlApi.SlowLogs)                           // 慢日志
	}
}

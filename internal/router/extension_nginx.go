package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ExtensionNginxRouter struct{}

func (s *ExtensionNginxRouter) Init(Router *gin.RouterGroup) {
	nginxR := Router.Group("extension_nginx")
	nginxApi := api.GroupApp.ExtensionNginxApiApp
	{
		nginxR.POST("Info", nginxApi.Info)                                   // 获取Nginx的信息
		nginxR.POST("Install", nginxApi.Install)                             // 安装Nginx
		nginxR.POST("Uninstall", nginxApi.Uninstall)                         // 卸载nginx
		nginxR.POST("SetStatus", nginxApi.SetStatus)                         // 设置Nginx的运行状态
		nginxR.POST("PerformanceConfig", nginxApi.PerformanceConfig)         // 获取性能配置
		nginxR.POST("SavePerformanceConfig", nginxApi.SavePerformanceConfig) // 保存性能配置
		nginxR.POST("LoadStatus", nginxApi.LoadStatus)                       // 获取负载状态
	}
}

package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ExtensionNodejsRouter struct{}

func (s *ExtensionNodejsRouter) Init(Router *gin.RouterGroup) {
	nodejsR := Router.Group("extension_nodejs")
	nodejsApi := api.GroupApp.ExtensionNodejsApiApp
	{
		nodejsR.POST("Info", nodejsApi.Info)                                 // 扩展信息
		nodejsR.POST("Config", nodejsApi.Config)                             // nodejs配置
		nodejsR.POST("SetRegistrySources", nodejsApi.SetRegistrySources)     // 设置镜像源
		nodejsR.POST("SetVersionUrl", nodejsApi.SetVersionUrl)               // 设置版本地址
		nodejsR.POST("VersionList", nodejsApi.VersionList)                   // 版本列表
		nodejsR.POST("UpdateVersionList", nodejsApi.UpdateVersionList)       // 更新版本列表
		nodejsR.POST("Install", nodejsApi.Install)                           // 安装nodejs
		nodejsR.POST("Uninstall", nodejsApi.Uninstall)                       // 卸载nodejs
		nodejsR.POST("SetDefaultEnv", nodejsApi.SetDefaultEnv)               // 设置默认环境变量
		nodejsR.POST("NodeModulesList", nodejsApi.NodeModulesList)           // node_modules列表
		nodejsR.POST("OperationNodeModules", nodejsApi.OperationNodeModules) // 操作node_modules
	}
}

package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ExtensionPHPRouter struct{}

func (s *ExtensionPHPRouter) Init(Router *gin.RouterGroup) {
	phpR := Router.Group("extension_php")
	phpApi := api.GroupApp.ExtensionPHPApiApp
	{
		phpR.POST("Info", phpApi.Info)                                   // 获取基本信息
		phpR.POST("Install", phpApi.Install)                             // 安装php
		phpR.POST("Uninstall", phpApi.Uninstall)                         // 卸载php
		phpR.POST("Status", phpApi.Status)                               // 获取php状态
		phpR.POST("SetStatus", phpApi.SetStatus)                         // 设置php状态
		phpR.POST("ExtensionList", phpApi.ExtensionList)                 // 获取php扩展列表
		phpR.POST("InstallLib", phpApi.InstallLib)                       // 安装php扩展
		phpR.POST("UninstallLib", phpApi.UninstallLib)                   // 卸载php扩展
		phpR.POST("GeneralConfig", phpApi.GeneralConfig)                 // 获取php通用配置
		phpR.POST("SaveGeneralConfig", phpApi.SaveGeneralConfig)         // 保存php通用配置
		phpR.POST("DisableFunctionList", phpApi.DisableFunctionList)     // 获取禁用函数列表
		phpR.POST("AddDisableFunction", phpApi.AddDisableFunction)       // 添加禁用函数
		phpR.POST("DeleteDisableFunction", phpApi.DeleteDisableFunction) // 删除禁用函数
		phpR.POST("PerformanceConfig", phpApi.PerformanceConfig)         // 获取性能配置
		phpR.POST("SavePerformanceConfig", phpApi.SavePerformanceConfig) // 保存性能配置
		phpR.POST("LoadStatus", phpApi.LoadStatus)                       // 获取负载状态
		phpR.POST("FpmLog", phpApi.FpmLog)                               // fpm日志
		phpR.POST("FpmSlowLog", phpApi.FpmSlowLog)                       // fpm慢日志
		phpR.POST("PHPInfo", phpApi.PHPInfo)                             // PHPInfo
		phpR.POST("PHPInfoHtml", phpApi.PHPInfoHtml)                     // PHPInfo网页版
		phpR.POST("CmdVersion", phpApi.CmdVersion)                       // 获取php命令行版本
		phpR.POST("SetCmdVersion", phpApi.SetCmdVersion)                 // 设置php命令行版本
	}
}

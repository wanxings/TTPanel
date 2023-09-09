package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type LinuxToolsRouter struct{}

func (s *LinuxToolsRouter) Init(Router *gin.RouterGroup) {
	linuxToolsR := Router.Group("linux_tools")
	linuxToolsApi := api.GroupApp.LinuxToolsApiApp
	{ //dns工具
		linuxToolsR.POST("GetDnsConfig", linuxToolsApi.GetDnsConfig)         // 获取dns配置
		linuxToolsR.POST("SaveDnsConfig", linuxToolsApi.SaveDnsConfig)       // 保存dns配置
		linuxToolsR.POST("TestDnsConfig", linuxToolsApi.TestDnsConfig)       // 测试dns配置
		linuxToolsR.POST("RecoverDnsConfig", linuxToolsApi.RecoverDnsConfig) // 恢复dns配置
	}
	//{ //swap工具
	//	linuxToolsR.POST("Install", linuxToolsApi.Install) // 安装docker
	//}
	{ //网卡工具
		linuxToolsR.POST("GetNetworkConfig", linuxToolsApi.GetNetworkConfig) // 获取网卡配置
	}
	{ //时区工具
		linuxToolsR.POST("GetTimeZoneConfig", linuxToolsApi.GetTimeZoneConfig) // 获取时区配置
		linuxToolsR.POST("SetTimeZone", linuxToolsApi.SetTimeZone)             // 设置时区
		linuxToolsR.POST("SyncDate", linuxToolsApi.SyncDate)                   // 同步服务器时间
	}
	{ //host工具
		linuxToolsR.POST("GetHostsConfig", linuxToolsApi.GetHostsConfig) // 获取host配置
		linuxToolsR.POST("AddHosts", linuxToolsApi.AddHosts)             // 添加host
		linuxToolsR.POST("RemoveHosts", linuxToolsApi.RemoveHosts)       // 删除host
	}
}

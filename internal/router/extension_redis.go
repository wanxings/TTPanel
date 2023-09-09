package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ExtensionRedisRouter struct{}

func (s *ExtensionRedisRouter) Init(Router *gin.RouterGroup) {
	redisR := Router.Group("extension_redis")
	redisApi := api.GroupApp.ExtensionRedisApiApp
	{
		redisR.POST("Info", redisApi.Info)                                   // 获取Redis的信息
		redisR.POST("SetStatus", redisApi.SetStatus)                         // 设置Redis的状态
		redisR.POST("Install", redisApi.Install)                             // 安装Redis
		redisR.POST("Uninstall", redisApi.Uninstall)                         // 卸载Redis
		redisR.POST("PerformanceConfig", redisApi.PerformanceConfig)         // 性能配置
		redisR.POST("SavePerformanceConfig", redisApi.SavePerformanceConfig) // 保存性能配置
		redisR.POST("LoadStatus", redisApi.LoadStatus)                       // 获取Redis的负载状态
		redisR.POST("PersistentConfig", redisApi.PersistentConfig)           // 持久化配置
		redisR.POST("SavePersistentConfig", redisApi.SavePersistentConfig)   // 保存持久化配置
	}
}

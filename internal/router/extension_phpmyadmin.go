package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ExtensionPhpmyadminRouter struct{}

func (s *ExtensionPhpmyadminRouter) Init(Router *gin.RouterGroup) {
	phpmyadminRouter := Router.Group("extension_phpmyadmin")
	phpmyadminApi := api.GroupApp.ExtensionPhpmyadminApiApp
	{
		phpmyadminRouter.POST("Info", phpmyadminApi.Info)           // 获取Phpmyadmin的信息
		phpmyadminRouter.POST("Install", phpmyadminApi.Install)     // 安装Phpmyadmin
		phpmyadminRouter.POST("Uninstall", phpmyadminApi.Uninstall) // 卸载Phpmyadmin
		phpmyadminRouter.POST("GetConfig", phpmyadminApi.GetConfig) // 获取配置
		phpmyadminRouter.POST("SetConfig", phpmyadminApi.SetConfig) // 获取配置
	}
}

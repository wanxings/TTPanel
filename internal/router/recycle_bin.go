package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type RecycleBinRouter struct{}

func (s *RecycleBinRouter) Init(Router *gin.RouterGroup) {
	recycleBinRouter := Router.Group("recycleBin")
	recycleBinApi := api.GroupApp.RecycleBinApiApp
	{
		recycleBinRouter.POST("Config", recycleBinApi.Config)                         // 获取配置
		recycleBinRouter.POST("SetConfig", recycleBinApi.SetConfig)                   // 设置配置
		recycleBinRouter.POST("List", recycleBinApi.List)                             // 恢复文件
		recycleBinRouter.POST("RecoveryFile", recycleBinApi.RecoveryFile)             // 恢复文件
		recycleBinRouter.POST("DeleteRecoveryFile", recycleBinApi.DeleteRecoveryFile) // 删除文件
		recycleBinRouter.POST("ClearRecycleBin", recycleBinApi.ClearRecycleBin)       // 清空回收站
	}
}

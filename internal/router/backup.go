package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type BackupRouter struct{}

func (s *BackupRouter) Init(Router *gin.RouterGroup) {
	backupRouter := Router.Group("backup")
	backupApi := api.GroupApp.BackupApiApp
	{
		backupRouter.POST("Database", backupApi.Database) // 备份数据库
		backupRouter.POST("Project", backupApi.Project)   // 备份项目
		backupRouter.POST("Dir", backupApi.Dir)           // 备份目录
		backupRouter.POST("Panel", backupApi.Panel)       // 备份面板
		backupRouter.POST("List", backupApi.List)         // 备份列表
	}

}

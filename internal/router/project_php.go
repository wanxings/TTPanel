package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ProjectPHPRouter struct{}

func (s *ProjectPHPRouter) Init(Router *gin.RouterGroup) {
	phpProjectRouter := Router.Group("project_php")
	phpProjectApi := api.GroupApp.ProjectPHPApiApp
	{
		phpProjectRouter.POST("Create", phpProjectApi.Create)                               // 创建PHP项目
		phpProjectRouter.POST("List", phpProjectApi.List)                                   // PHP项目列表
		phpProjectRouter.POST("ProjectInfo", phpProjectApi.ProjectInfo)                     // PHP项目详情
		phpProjectRouter.POST("SetStatus", phpProjectApi.SetStatus)                         // 设置PHP项目状态
		phpProjectRouter.POST("SetExpireTime", phpProjectApi.SetExpireTime)                 // 设置PHP项目到期时间
		phpProjectRouter.POST("Delete", phpProjectApi.Delete)                               // 删除PHP项目
		phpProjectRouter.POST("UsingPHPVersion", phpProjectApi.UsingPHPVersion)             // 使用的PHP版本
		phpProjectRouter.POST("SwitchUsingPHPVersion", phpProjectApi.SwitchUsingPHPVersion) // 切换使用的PHP版本
	}
	{
		phpProjectRouter.POST("SetRunPath", phpProjectApi.SetRunPath) // 设置运行目录
		phpProjectRouter.POST("SetPath", phpProjectApi.SetPath)       // 设置项目目录
		phpProjectRouter.POST("SetUserIni", phpProjectApi.SetUserIni) // 设置防跨站攻击
	}
}

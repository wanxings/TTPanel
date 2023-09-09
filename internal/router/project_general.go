package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ProjectGeneralRouter struct{}

func (s *ProjectGeneralRouter) Init(Router *gin.RouterGroup) {
	generalProjectRouter := Router.Group("project_general")
	generalProjectApi := api.GroupApp.ProjectGeneralApiApp
	{
		generalProjectRouter.POST("Create", generalProjectApi.Create)                       // 创建通用项目
		generalProjectRouter.POST("List", generalProjectApi.List)                           // 获取通用项目列表
		generalProjectRouter.POST("Delete", generalProjectApi.Delete)                       // 删除通用项目
		generalProjectRouter.POST("Info", generalProjectApi.Info)                           // 获取通用项目详情
		generalProjectRouter.POST("SaveProjectConfig", generalProjectApi.SaveProjectConfig) // 保存项目配置
		generalProjectRouter.POST("SetStatus", generalProjectApi.SetStatus)                 // 设置项目状态
	}
}

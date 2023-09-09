package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type CronTaskRouter struct{}

func (t *CronTaskRouter) Init(Router *gin.RouterGroup) {
	taskRouter := Router.Group("task")
	taskApi := api.GroupApp.CronTaskApiApp
	{
		taskRouter.POST("BatchCreateCronTask", taskApi.BatchCreateCronTask) // 创建-批量
		taskRouter.POST("List", taskApi.List)                               // 列表
		taskRouter.POST("Details", taskApi.Details)                         // 任务详情
		taskRouter.POST("Edit", taskApi.Edit)                               // 编辑任务
		taskRouter.POST("BatchSetStatus", taskApi.BatchSetStatus)           // 设置状态-批量
		taskRouter.POST("BatchDelete", taskApi.BatchDelete)                 // 删除-批量
		taskRouter.POST("BatchRun", taskApi.BatchRun)                       // 执行任务-批量
		taskRouter.POST("GetLog", taskApi.GetLog)                           // 获取日志
		taskRouter.POST("ClearLog", taskApi.ClearLog)                       // 清空日志
	}
}

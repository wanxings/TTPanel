package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type QueueTaskRouter struct{}

func (s *QueueTaskRouter) Init(Router *gin.RouterGroup) {
	queueTaskRouter := Router.Group("queueTask")
	queueTaskApi := api.GroupApp.QueueTaskApiApp
	{
		queueTaskRouter.POST("RunningCount", queueTaskApi.RunningCount) // 运行中的任务数量
		queueTaskRouter.POST("TaskList", queueTaskApi.TaskList)         // 任务列表
		queueTaskRouter.POST("DelTask", queueTaskApi.DelTask)           // 删除任务
		queueTaskRouter.POST("ClearTask", queueTaskApi.ClearTask)       // 清空任务
	}
}

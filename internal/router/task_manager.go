package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type TaskManagerRouter struct{}

func (m *TaskManagerRouter) Init(Router *gin.RouterGroup) {
	taskManagerRouter := Router.Group("manager")
	taskManagerApi := api.GroupApp.TaskManagerApiApp
	{
		taskManagerRouter.POST("ProcessList", taskManagerApi.ProcessList) //进程列表
		taskManagerRouter.POST("KillProcess", taskManagerApi.KillProcess) //结束进程-批量
	}
	{
		taskManagerRouter.POST("StartupList", taskManagerApi.StartupList) //启动项列表

	}
	{
		taskManagerRouter.POST("ServiceList", taskManagerApi.ServiceList)     //系统服务列表
		taskManagerRouter.POST("DeleteService", taskManagerApi.DeleteService) //删除服务
		taskManagerRouter.POST("SetRunLevel", taskManagerApi.SetRunLevel)     //设置运行级别状态
	}
	{
		taskManagerRouter.POST("ConnectionList", taskManagerApi.ConnectionList) //网络列表
	}
	{
		taskManagerRouter.POST("LinuxUserList", taskManagerApi.LinuxUserList)     //Linux用户列表
		taskManagerRouter.POST("DeleteLinuxUser", taskManagerApi.DeleteLinuxUser) //删除Linux用户
	}
}

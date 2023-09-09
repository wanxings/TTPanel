package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model/request"
	"github.com/gin-gonic/gin"
)

type TaskManagerApi struct{}

// ProcessList
// @Tags      manager
// @Summary   进程列表
// @Router    manager/ProcessList [post]
func (m *TaskManagerApi) ProcessList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.TaskManagerServiceApp.ProcessList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, 0, 0, 0)
}

// KillProcess
// @Tags      manager
// @Summary   结束进程
// @Router    manager/KillProcess [post]
func (m *TaskManagerApi) KillProcess(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.KillProcessR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("KillProcess.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.TaskManagerServiceApp.KillProcess(param.Pid, param.Force)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByTaskManager, helper.MessageWithMap("task_manager.KillProcess", map[string]any{"Pid": param.Pid, "Force": param.Force}))
	response.ToResponse(data)
}

// StartupList
// @Tags      manager
// @Summary   启动项列表
// @Router    manager/StartupList [post]
func (m *TaskManagerApi) StartupList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.TaskManagerServiceApp.StartupList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, 0, 0, 0)
}

// ServiceList
// @Tags      manager
// @Summary   系统服务列表
// @Router    manager/ServiceList [post]
func (m *TaskManagerApi) ServiceList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.TaskManagerServiceApp.ServiceList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, 0, 0, 0)
}

// DeleteService
// @Tags      manager
// @Summary   删除服务
// @Router    manager/DeleteService [post]
func (m *TaskManagerApi) DeleteService(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.DeleteServiceR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DeleteService.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.TaskManagerServiceApp.DeleteService(param.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByTaskManager, helper.MessageWithMap("task_manager.DeleteService", map[string]any{"Name": param.Name}))

	response.ToResponse(data)
}

// SetRunLevel
// @Tags      manager
// @Summary   设置运行级别状态
// @Router    manager/SetRunLevel [post]
func (m *TaskManagerApi) SetRunLevel(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetRunLevelR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("SetRunLevel.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.TaskManagerServiceApp.SetRunLevel(param.Name, param.Level)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByTaskManager, helper.MessageWithMap("task_manager.SetRunLevel", map[string]any{"Name": param.Name, "RunLevel": param.Level}))
	response.ToResponse(data)
}

// ConnectionList
// @Tags      manager
// @Summary   网络连接列表
// @Router    manager/ConnectionList [post]
func (m *TaskManagerApi) ConnectionList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.TaskManagerServiceApp.ConnectionList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// LinuxUserList
// @Tags      manager
// @Summary   Linux用户列表
// @Router    manager/LinuxUserList [post]
func (m *TaskManagerApi) LinuxUserList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.TaskManagerServiceApp.LinuxUserList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// DeleteLinuxUser
// @Tags      manager
// @Summary   删除Linux用户
// @Router    manager/DeleteLinuxUser [post]
func (m *TaskManagerApi) DeleteLinuxUser(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.DeleteLinuxUserR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DeleteLinuxUser.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.TaskManagerServiceApp.DeleteLinuxUser(param.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByTaskManager, helper.MessageWithMap("task_manager.DeleteLinuxUser", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

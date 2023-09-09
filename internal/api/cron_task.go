package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/pkg/util"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type CronTaskApi struct{}

// BatchCreateCronTask
// @Tags      Task
// @Summary   创建-批量
// @Router    /task/BatchCreateCronTask [post]
func (t *CronTaskApi) BatchCreateCronTask(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchCreateCronTaskR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchCreate.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	errList := ServiceGroupApp.CronTaskServiceApp.BatchCreate(param.List)
	if errList != nil {
		var errS []string
		for _, err := range errList {
			errS = append(errS, err.Error())
		}
		response.ToErrorResponse(errcode.ServerError.WithDetails(errS...))
		return
	}
	taskList := make([]string, 0)
	for _, v := range param.List {
		taskList = append(taskList, v.Name)
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByCronTask, helper.MessageWithMap("cron_task.Create", map[string]any{"Name": strings.Join(taskList, ",")}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))

}

// List
// @Tags      Task
// @Summary   计划任务列表
// @Router    /task/List [post]
func (t *CronTaskApi) List(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TaskListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("List.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, total, err := ServiceGroupApp.CronTaskServiceApp.List(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), 0, 1)
}

// Details
// @Tags      Task
// @Summary   计划任务详情
// @Router    /task/Details [post]
func (t *CronTaskApi) Details(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TaskIdR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Details.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.CronTaskServiceApp.Get(param.Id)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)

}

// Edit
// @Tags      Edit
// @Summary   编辑计划任务
// @Router    /task/Edit [post]
func (t *CronTaskApi) Edit(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.EditR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Edit.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.CronTaskServiceApp.Edit(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByCronTask, helper.MessageWithMap("cron_task.Edit", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// BatchSetStatus
// @Tags      BatchSetStatus
// @Summary   设置状态-批量
// @Router    /task/BatchSetStatus [post]
func (t *CronTaskApi) BatchSetStatus(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchSetStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Edit.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	taskNameList, err := ServiceGroupApp.CronTaskServiceApp.BatchSetStatus(param.IDs, param.Status)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByCronTask, helper.MessageWithMap("cron_task.SetStatus", map[string]any{"Name": strings.Join(taskNameList, ","), "Status": param.Status}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// BatchDelete
// @Tags      BatchDelete
// @Summary   删除-批量
// @Router    /task/BatchSetStatus [post]
func (t *CronTaskApi) BatchDelete(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchIDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Edit.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	taskNameList, err := ServiceGroupApp.CronTaskServiceApp.BatchDelete(param.IDs)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByCronTask, helper.MessageWithMap("cron_task.Delete", map[string]any{"Name": strings.Join(taskNameList, ",")}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// BatchRun
// @Tags      BatchRun
// @Summary   执行-批量
// @Router    /task/BatchRun [post]
func (t *CronTaskApi) BatchRun(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchIDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchRun.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	taskNameList, err := ServiceGroupApp.CronTaskServiceApp.BatchRun(param.IDs)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByCronTask, helper.MessageWithMap("cron_task.Execute", map[string]any{"Name": strings.Join(taskNameList, ",")}))
	response.ToResponseMsg(helper.Message("tips.ExecuSuccess"))
}

// GetLog
// @Tags      GetLog
// @Summary   获取日志
// @Router    /task/GetLog [post]
func (t *CronTaskApi) GetLog(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.GetLogR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("GetLog.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.CronTaskServiceApp.GetExecutionLog(param.Id, param.Line)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// ClearLog
// @Tags      ClearLog
// @Summary   清空日志
// @Router    /task/ClearLog [post]
func (t *CronTaskApi) ClearLog(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TaskIdR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ClearLog.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//查询任务详情
	taskInfo, err := (&model.CronTask{ID: param.Id}).Get(global.PanelDB)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
	}
	if taskInfo.ID == 0 {
		response.ToErrorResponse(errcode.ServerError.WithDetails(errors.New("not found").Error()))
		return
	}
	//清空执行日志
	shellPath := global.Config.System.ServerPath + "/cron/" + taskInfo.Hash + ".log"
	_, _ = util.ExecShell("echo '' > " + shellPath)

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByCronTask, helper.MessageWithMap("cron_task.ClearLog", map[string]any{"Name": taskInfo.Name}))
	response.ToResponseMsg(helper.Message("tips.ClearSuccess"))
}

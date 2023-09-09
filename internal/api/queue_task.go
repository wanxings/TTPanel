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

type QueueTaskApi struct{}

// RunningCount
// @Tags      QueueTask
// @Summary   运行中的任务数量
// @Router    /queueTask/RunningCount [post]
func (s *QueueTaskApi) RunningCount(c *gin.Context) {
	response := app.NewResponse(c)
	count, err := ServiceGroupApp.QueueTaskServiceApp.RunningCount()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(count)
}

// TaskList
// @Tags      QueueTask
// @Summary   任务列表
// @Router    /queueTask/TaskList [post]
func (s *QueueTaskApi) TaskList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.QueueTaskListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("QueueTask.TaskList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)

	data, total, err := ServiceGroupApp.QueueTaskServiceApp.TaskList(param.Status, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// DelTask
// @Tags      QueueTask
// @Summary   删除任务
// @Router    /queueTask/DelTask [post]
func (s *QueueTaskApi) DelTask(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.QueueTaskDelR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("QueueTask.DelTask.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取任务
	queueTask, err := ServiceGroupApp.QueueTaskServiceApp.GetQueueTaskByID(param.ID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//删除任务
	err = ServiceGroupApp.QueueTaskServiceApp.DelTask(queueTask)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Delete", map[string]any{"Name": queueTask.Name}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// ClearTask
// @Tags      QueueTask
// @Summary   清空任务
// @Router    /queueTask/ClearTask [post]
func (s *QueueTaskApi) ClearTask(c *gin.Context) {
	response := app.NewResponse(c)
	err := ServiceGroupApp.QueueTaskServiceApp.ClearTask()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.Message("queuetask.ClearCompletedTasks"))
	response.ToResponseMsg(helper.Message("tips.ClearSuccess"))
}

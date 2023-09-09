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

type MonitorApi struct{}

// Base
// @Tags     Monitor
// @Summary   获取基础信息
// @Router    /api/monitor/Base [get]
func (s *MonitorApi) Base(c *gin.Context) {
	response := app.NewResponse(c)
	//获取基础信息
	data, err := ServiceGroupApp.MonitorServiceApp.Base()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// Logs
// @Tags     Monitor
// @Summary   获取日志
// @Router    /api/monitor/Logs [get]
func (s *MonitorApi) Logs(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.LogsR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Monitor.Logs.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取日志
	data, err := ServiceGroupApp.MonitorServiceApp.Logs(param.Type, param.StartTime, param.EndTime)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// ClearAllLogs
// @Tags     Monitor
// @Summary   清空日志
// @Router    /api/monitor/ClearAllLogs [get]
func (s *MonitorApi) ClearAllLogs(c *gin.Context) {
	response := app.NewResponse(c)

	//清空日志
	err := ServiceGroupApp.MonitorServiceApp.ClearAllLogs()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByMonitor, helper.Message("monitor.Clear"))
	response.ToResponse(gin.H{})
}

// Config
// @Tags     Monitor
// @Summary   获取配置
// @Router    /api/monitor/Config [get]
func (s *MonitorApi) Config(c *gin.Context) {
	response := app.NewResponse(c)
	//获取配置
	response.ToResponse(global.Config.Monitor)
}

// SaveConfig
// @Tags     Monitor
// @Summary   保存配置
// @Router    /api/monitor/SaveConfig [post]
func (s *MonitorApi) SaveConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.MonitorConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Monitor.SaveConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//保存配置
	err := ServiceGroupApp.MonitorServiceApp.SaveConfig(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByMonitor, helper.Message("monitor.SaveConfig"))
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// EventConfig
// @Tags     Monitor
// @Summary
// @Router    /api/monitor/EventConfig [post]
func (s *MonitorApi) EventConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取配置
	EventConfig := ServiceGroupApp.MonitorServiceApp.GetEventConfig()
	response.ToResponse(EventConfig)
}

// SaveEventConfig
// @Tags     Monitor
// @Summary
// @Router    /api/monitor/SaveEventConfig [post]
func (s *MonitorApi) SaveEventConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.MonitorEventConfig{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Monitor.SaveEventConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.MonitorServiceApp.SaveEventConfig(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByMonitor, helper.Message("monitor.SaveEventConfig"))
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// EventList
// @Tags     Monitor
// @Summary
// @Router    /api/monitor/EventList [post]
func (s *MonitorApi) EventList(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.EventListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Monitor.EventList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)

	//获取列表
	list, total, err := ServiceGroupApp.MonitorServiceApp.EventList(param.Status, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseList(list, int(total), param.Limit, param.Page)
}

// BatchSetEventStatus
// @Tags     Monitor
// @Summary 批量设置事件状态
// @Router    /api/monitor/BatchSetEventStatus [post]
func (s *MonitorApi) BatchSetEventStatus(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BatchSetEventStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Monitor.BatchSetEventStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.MonitorServiceApp.BatchSetEventStatus(param.Ids, param.Status)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseMsg(helper.Message("tips.SetStatusSuccess"))
}

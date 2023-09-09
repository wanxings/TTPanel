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

type LogAuditApi struct{}

// PanelOperationLogList
// @Tags      LogAudit
// @Summary   面板操作日志列表
// @Router    /log_audit/PanelOperationLogList [post]
func (s *LogAuditApi) PanelOperationLogList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.OperationLogListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("PanelOperationLogList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)

	data, total, err := ServiceGroupApp.LogAuditServiceApp.PanelOperationLogList(&param, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, total, param.Limit, param.Page)
}

// ClearPanelOperationLog
// @Tags      LogAudit
// @Summary   清空面板操作日志
// @Router    /log_audit/ClearPanelOperationLog [post]
func (s *LogAuditApi) ClearPanelOperationLog(c *gin.Context) {
	response := app.NewResponse(c)
	err := ServiceGroupApp.LogAuditServiceApp.ClearPanelOperationLog()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByLogAudit, helper.Message("log_audit.ClearPanelOperationLog"))
	response.ToResponseMsg(helper.Message("tips.ClearSuccess"))
}

// LogFileOccupancy
// @Tags      LogAudit
// @Summary   日志占用
// @Router    /log_audit/LogFileOccupancy [post]
func (s *LogAuditApi) LogFileOccupancy(c *gin.Context) {
	response := app.NewResponse(c)
	data := ServiceGroupApp.LogAuditServiceApp.LogFileOccupancy()
	response.ToResponse(data)
}

// SSHLoginLogList
// @Tags      LogAudit
// @Summary   ssh登录日志
// @Router    /log_audit/SSHLoginLogList [post]
func (s *LogAuditApi) SSHLoginLogList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SSHLoginLogListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("SSHLoginLogList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)

	data, total, err := ServiceGroupApp.LogAuditServiceApp.SSHLoginLogList(param.Query, param.Status, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, total, param.Limit, param.Page)
}

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

type SSHManageApi struct{}

// GetSSHInfo
// @Tags     SSHManage
// @Summary   获取SSH信息
// @Router    /system/ssh_manage/GetSSHInfo [post]
func (s *SSHManageApi) GetSSHInfo(c *gin.Context) {
	response := app.NewResponse(c)
	ResponseData, err := ServiceGroupApp.SSHManageServiceApp.GetSSHInfo()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(&ResponseData)
}

// SetSSHStatus
// @Tags     SSHManage
// @Summary   设置SSH状态
// @Router    /system/ssh_manage/SetSSHStatus [post]
func (s *SSHManageApi) SetSSHStatus(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetSSHStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ssh_manage.SetSSHStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	if err := ServiceGroupApp.SSHManageServiceApp.SetSSHStatus(param.Action); err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystem, helper.MessageWithMap("ssh_manage.SetStatus", map[string]any{"Status": param.Action}))
	response.ToResponseMsg(helper.Message("tips.SetStatusSuccess"))
}

// OperateSSHKeyLogin
// @Tags     SSHManage
// @Summary   操作SSH密钥登录
// @Router    /system/ssh_manage/OperateSSHKeyLogin [post]
func (s *SSHManageApi) OperateSSHKeyLogin(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.OperateSSHKeyLoginR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ssh_manage.OperateSSHKeyLogin.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	key, err := ServiceGroupApp.SSHManageServiceApp.OperateSSHKeyLogin(param.Action, param.KeyType)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystem, helper.MessageWithMap("ssh_manage.OperateSSHKeyLogin", map[string]any{"Status": param.Action, "KeyType": param.KeyType}))
	response.ToResponse(key)
}

// OperatePasswordLogin
// @Tags     SSHManage
// @Summary   操作密码登录
// @Router    /system/ssh_manage/OperatePasswordLogin [post]
func (s *SSHManageApi) OperatePasswordLogin(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.OperatePasswordLoginR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ssh_manage.OperatePasswordLogin.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SSHManageServiceApp.OperatePasswordLogin(param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystem, helper.MessageWithMap("ssh_manage.OperatePasswordLogin", map[string]any{"Status": param.Action}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// GetSSHLoginStatistics
// @Tags     SSHManage
// @Summary   获取SSH登录统计
// @Router    /system/ssh_manage/GetSSHLoginStatistics [post]
func (s *SSHManageApi) GetSSHLoginStatistics(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.GetSSHLoginStatisticsR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ssh_manage.GetSSHLoginStatistics.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	ResponseData, err := ServiceGroupApp.SSHManageServiceApp.GetSSHLoginStatistics(param.Refresh)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(&ResponseData)
}

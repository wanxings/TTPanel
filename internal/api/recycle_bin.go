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

type RecycleBinApi struct{}

// Config
// @Tags     RecycleBin
// @Summary  获取回收站配置
// @Router    /system/recycleBin/Config [get]
func (s *RecycleBinApi) Config(c *gin.Context) {
	response := app.NewResponse(c)
	//获取回收站配置
	data, err := ServiceGroupApp.RecycleBinServiceApp.Config()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SetConfig
// @Tags     RecycleBin
// @Summary  设置回收站状态
// @Router    /system/recycleBin/SetConfig [post]
func (s *RecycleBinApi) SetConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SetRecycleBinStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("RecycleBin.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//设置回收站状态
	err := ServiceGroupApp.RecycleBinServiceApp.SetStatus(param.ExplorerStatus)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByRecycleBin, helper.MessageWithMap("recyclebin.SetStatus", map[string]any{"Status": param.ExplorerStatus}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// List
// @Tags     RecycleBin
// @Summary  获取回收站列表
// @Router    /system/recycleBin/List [post]
func (s *RecycleBinApi) List(c *gin.Context) {
	response := app.NewResponse(c)

	//获取回收站列表
	data, err := ServiceGroupApp.RecycleBinServiceApp.List()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// RecoveryFile
// @Tags     RecycleBin
// @Summary  从回收站恢复文件
// @Router    /system/recycleBin/RecoveryFile [post]
func (s *RecycleBinApi) RecoveryFile(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.RecoveryFileR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("RecycleBin.RecoveryFile.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//从回收站恢复文件
	recoveryName, err := ServiceGroupApp.RecycleBinServiceApp.RecoveryFile(param.Hash, param.Cover)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByRecycleBin, helper.MessageWithMap("recyclebin.Recovery", map[string]any{"Path": recoveryName}))
	response.ToResponseMsg(helper.Message("tips.RecoverySuccess"))
}

// DeleteRecoveryFile
// @Tags     RecycleBin
// @Summary  删除回收站文件
// @Router    /system/recycleBin/DeleteRecoveryFile [post]
func (s *RecycleBinApi) DeleteRecoveryFile(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.DeleteRecoveryFileR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("RecycleBin.DeleteRecoveryFile.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//删除回收站文件
	deleteName, err := ServiceGroupApp.RecycleBinServiceApp.DeleteRecoveryFile(param.Hash)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByRecycleBin, helper.MessageWithMap("recyclebin.Delete", map[string]any{"Name": deleteName}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// ClearRecycleBin
// @Tags     RecycleBin
// @Summary  清空回收站
// @Router    /system/recycleBin/ClearRecycleBin [post]
func (s *RecycleBinApi) ClearRecycleBin(c *gin.Context) {
	response := app.NewResponse(c)
	//清空回收站
	err := ServiceGroupApp.RecycleBinServiceApp.ClearRecycleBin()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByRecycleBin, helper.Message("recyclebin.Clear"))
	response.ToResponseMsg(helper.Message("tips.ClearSuccess"))
}

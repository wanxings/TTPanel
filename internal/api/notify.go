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

type NotifyApi struct{}

// AddNotifyChannel
// @Tags     Notify
// @Summary   创建通知
// @Router    /api/notify/AddNotifyChannel [post]
func (s *NotifyApi) AddNotifyChannel(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.AddNotifyChannelR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Notify.AddNotifyChannel.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.NotifyServiceApp.AddNotifyChannel(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByNotify, helper.MessageWithMap("notify.Add", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// NotifyChannelList
// @Tags     Notify
// @Summary   通知列表
// @Router    /api/notify/NotifyChannelList [post]
func (s *NotifyApi) NotifyChannelList(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.NotifyChannelListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Notify.NotifyChannelList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	list, total, err := ServiceGroupApp.NotifyServiceApp.NotifyChannelList(param.Query, param.Category, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(list, int(total), param.Limit, param.Page)
}

// TestNotifyChannel
// @Tags     Notify
// @Summary   测试通知
// @Router    /api/notify/TestNotifyChannel [post]
func (s *NotifyApi) TestNotifyChannel(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.TestNotifyChannelR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Notify.TestNotifyChannel.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.NotifyServiceApp.TestNotifyChannel(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.SendSuccess"))
}

// EditNotifyChannel
// @Tags     Notify
// @Summary   编辑通知
// @Router    /api/notify/EditNotifyChannel [post]
func (s *NotifyApi) EditNotifyChannel(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.EditNotifyChannelR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Notify.EditNotifyChannel.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.NotifyServiceApp.EditNotifyChannel(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByNotify, helper.MessageWithMap("notify.Edit", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// DeleteNotifyChannel
// @Tags     Notify
// @Summary   删除通知
// @Router    /api/notify/DeleteNotifyChannel [post]
func (s *NotifyApi) DeleteNotifyChannel(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.DeleteNotifyChannelR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Notify.DeleteNotifyChannel.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.NotifyServiceApp.DeleteNotifyChannel(param.ID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByNotify, helper.MessageWithMap("notify.Delete", map[string]any{"ID": param.ID}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

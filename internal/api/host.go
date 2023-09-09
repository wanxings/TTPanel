package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"github.com/gin-gonic/gin"
)

type HostApi struct{}

// AddHostCategory
// @Tags      Host
// @Summary   添加主机分类
// @Router    /host/AddHostCategory [post]
func (s *HostApi) AddHostCategory(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.AddHostCategoryR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.AddHostCategory.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.HostServiceApp.AddHostCategory(param.Name, param.Remark)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByHostManager, helper.MessageWithMap("host.AddHostCategory", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.AddSuccess"))
}

// HostCategoryList
// @Tags      Host
// @Summary   主机分类列表
// @Router    /host/HostCategoryList [post]
func (s *HostApi) HostCategoryList(c *gin.Context) {
	response := app.NewResponse(c)
	list, total, err := ServiceGroupApp.HostServiceApp.HostCategoryList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(list, int(total), 0, 0)
}

// EditHostCategory
// @Tags      Host
// @Summary   编辑主机分类
// @Router    /host/EditHostCategory [post]
func (s *HostApi) EditHostCategory(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.EditHostCategoryR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.AddHostCategory.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.HostServiceApp.EditHostCategory(param.ID, param.Name, param.Remark)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByHostManager, helper.MessageWithMap("host.EditHostCategory", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// DeleteHostCategory
// @Tags      Host
// @Summary   删除主机分类
// @Router    /host/DeleteHostCategory [post]
func (s *HostApi) DeleteHostCategory(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.DeleteHostCategoryR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.AddHostCategory.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.HostServiceApp.DeleteHostCategory(param.ID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByHostManager, helper.MessageWithMap("host.DeleteHostCategory", map[string]any{"ID": param.ID}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// AddHost
// @Tags      Host
// @Summary   添加主机
// @Router    /host/AddHost [post]
func (s *HostApi) AddHost(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.AddHostR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.AddHostCategory.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.HostServiceApp.AddHost(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByHostManager, helper.MessageWithMap("host.AddHost", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.AddSuccess"))
}

// HostList
// @Tags      Host
// @Summary   主机列表
// @Router    /host/HostList [post]
func (s *HostApi) HostList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.HostListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.HostList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)

	//获取列表
	data, total, err := ServiceGroupApp.HostServiceApp.HostList(param.Query, param.CId, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// DeleteHost
// @Tags      Host
// @Summary   删除主机
// @Router    /host/DeleteHost [post]
func (s *HostApi) DeleteHost(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.DeleteHostR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.DeleteHost.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	host, err := ServiceGroupApp.HostServiceApp.DeleteHost(param.ID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByHostManager, helper.MessageWithMap("host.DeleteHost", map[string]any{"Name": host.Name}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// Terminal
// @Tags      Host
// @Summary   主机终端
// @Router    /host/Terminal [post]
func (s *HostApi) Terminal(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TerminalR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.Terminal.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//查询主机信息
	//查询主机是否存在
	host, err := (&model.Host{ID: param.HostId}).Get(global.PanelDB)
	if err != nil || host.ID == 0 {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByHostManager, helper.MessageWithMap("host.ConnectTerminal", map[string]any{"Name": host.Name}))
	ServiceGroupApp.HostServiceApp.Terminal(c, &param, host)

}

// ShortcutCommandList
// @Tags      Host
// @Summary   快捷命令列表
// @Router    /host/ShortcutCommandList [post]
func (s *HostApi) ShortcutCommandList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.ShortcutCommandListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.ShortcutCommandList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	data, total, err := ServiceGroupApp.HostServiceApp.ShortcutCommandList(param.Query, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// AddShortcutCommand
// @Tags      Host
// @Summary   添加快捷命令
// @Router    /host/AddShortcutCommand [post]
func (s *HostApi) AddShortcutCommand(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.AddShortcutCommandR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.AddShortcutCommand.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.HostServiceApp.AddShortcutCommand(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByHostManager, helper.MessageWithMap("host.AddShortcutCommand", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.AddSuccess"))
}

// DeleteShortcutCommand
// @Tags      Host
// @Summary   删除快捷命令
// @Router    /host/DeleteShortcutCommand [post]
func (s *HostApi) DeleteShortcutCommand(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.DeleteShortcutCommandR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.DeleteShortcutCommand.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	shortcutCommand, err := ServiceGroupApp.HostServiceApp.DeleteShortcutCommand(param.ID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByHostManager, helper.MessageWithMap("host.DeleteShortcutCommand", map[string]any{"Name": shortcutCommand.Name}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// EditShortcutCommand
// @Tags      Host
// @Summary   编辑快捷命令
// @Router    /host/EditShortcutCommand [post]
func (s *HostApi) EditShortcutCommand(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.EditShortcutCommandR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.EditShortcutCommand.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.HostServiceApp.EditShortcutCommand(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByHostManager, helper.MessageWithMap("host.EditShortcutCommand", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

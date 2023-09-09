package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model/request"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ExtensionPHPApi struct{}

// Info
// @Tags     PHP
// @Summary   获取基本信息
// @Router    /extension_php/Info [post]
func (s *ExtensionPHPApi) Info(c *gin.Context) {
	response := app.NewResponse(c)
	//获取基本信息
	data, err := ServiceGroupApp.ExtensionPHPServiceApp.Info()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// Install
// @Tags     PHP
// @Summary   安装php
// @Router    /extension_php/Install [post]
func (s *ExtensionPHPApi) Install(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php.Install.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//安裝php
	err := ServiceGroupApp.ExtensionPHPServiceApp.Install(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Install", map[string]any{"Name": "php-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// Uninstall
// @Tags     PHP
// @Summary   卸载php
// @Router    /extension_php/Uninstall [post]
func (s *ExtensionPHPApi) Uninstall(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php.Uninstall.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//卸载php
	err := ServiceGroupApp.ExtensionPHPServiceApp.Uninstall(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Uninstall", map[string]any{"Name": "php-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// Status
// @Tags     PHP
// @Summary   获取php状态
// @Router    /extension_php/Status [post]
func (s *ExtensionPHPApi) Status(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php.Status.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//获取php状态
	data := ServiceGroupApp.ExtensionPHPServiceApp.Status(param.Version)
	response.ToResponse(data)
}

// SetStatus
// @Tags     PHP
// @Summary   设置php状态
// @Router    /extension_php/SetStatus [post]
func (s *ExtensionPHPApi) SetStatus(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.PHPSetStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//设置php状态
	err := ServiceGroupApp.ExtensionPHPServiceApp.SetStatus(param.Version, param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExtension, helper.MessageWithMap("extension.SetStatus", map[string]any{"Name": "php", "Status": param.Action}))

	response.ToResponseMsg(helper.Message("tips.SetStatusSuccess"))
}

// ExtensionList
// @Tags     PHP
// @Summary   获取php扩展列表
// @Router    /extension_php/ExtensionList [post]
func (s *ExtensionPHPApi) ExtensionList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php.ExtensionList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取php扩展列表
	data, err := ServiceGroupApp.ExtensionPHPServiceApp.ExtensionList(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// InstallLib
// @Tags     PHP
// @Summary   安装php扩展
// @Router    /extension_php/InstallLib [post]
func (s *ExtensionPHPApi) InstallLib(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.InstallLibR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php.InstallLib.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//安装php扩展
	err := ServiceGroupApp.ExtensionPHPServiceApp.InstallLib(param.Version, param.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Install", map[string]any{"Name": "php-" + param.Version + "-" + param.Name}))

	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// UninstallLib
// @Tags     PHP
// @Summary   卸载php扩展
// @Router    /extension_php/UninstallLib [post]
func (s *ExtensionPHPApi) UninstallLib(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.InstallLibR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php.UninstallLib.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//卸载php扩展
	err := ServiceGroupApp.ExtensionPHPServiceApp.UninstallLib(param.Version, param.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Uninstall", map[string]any{"Name": "php-" + param.Version + "-" + param.Name}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// GeneralConfig
// @Tags     PHP
// @Summary   获取php通用配置
// @Router    /extension_php/GeneralConfig [post]
func (s *ExtensionPHPApi) GeneralConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.GeneralConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取php通用配置
	data, err := ServiceGroupApp.ExtensionPHPServiceApp.GeneralConfig(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SaveGeneralConfig
// @Tags     PHP
// @Summary   保存php通用配置
// @Router    /extension_php/SaveGeneralConfig [post]
func (s *ExtensionPHPApi) SaveGeneralConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SaveGeneralConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.SaveGeneralConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//保存php通用配置
	err := ServiceGroupApp.ExtensionPHPServiceApp.SaveGeneralConfig(param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// DisableFunctionList
// @Tags     PHP
// @Summary   获取php禁用函数列表
// @Router    /extension_php/DisableFunctionList [post]
func (s *ExtensionPHPApi) DisableFunctionList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.DisableFunctionList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取php禁用函数列表
	data, err := ServiceGroupApp.ExtensionPHPServiceApp.DisableFunctionList(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// AddDisableFunction
// @Tags     PHP
// @Summary   添加php禁用函数
// @Router    /extension_php/AddDisableFunction [post]
func (s *ExtensionPHPApi) AddDisableFunction(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.AddDisableFunctionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.AddDisableFunction.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//添加php禁用函数
	err := ServiceGroupApp.ExtensionPHPServiceApp.AddDisableFunction(param.Version, param.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseMsg(helper.Message("tips.AddSuccess"))
}

// DeleteDisableFunction
// @Tags     PHP
// @Summary   删除php禁用函数
// @Router    /extension_php/DeleteDisableFunction [post]
func (s *ExtensionPHPApi) DeleteDisableFunction(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.AddDisableFunctionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.DeleteDisableFunction.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//删除php禁用函数
	err := ServiceGroupApp.ExtensionPHPServiceApp.DeleteDisableFunction(param.Version, param.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	global.Log.Warnf("php DeleteDisableFunction ,Version:%s  Name:%s", param.Version, param.Name)
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// PerformanceConfig
// @Tags     PHP
// @Summary   获取php性能配置
// @Router    /extension_php/PerformanceConfig [post]
func (s *ExtensionPHPApi) PerformanceConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.PerformanceConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取php性能配置
	data, err := ServiceGroupApp.ExtensionPHPServiceApp.PerformanceConfig(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SavePerformanceConfig
// @Tags     PHP
// @Summary   保存php性能配置
// @Router    /extension_php/SavePerformanceConfig [post]
func (s *ExtensionPHPApi) SavePerformanceConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SavePerformanceConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.SavePerformanceConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//保存php性能配置
	err := ServiceGroupApp.ExtensionPHPServiceApp.SavePerformanceConfig(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// LoadStatus
// @Tags     PHP
// @Summary   获取php加载状态
// @Router    /extension_php/LoadStatus [post]
func (s *ExtensionPHPApi) LoadStatus(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.LoadStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取php加载状态
	data, err := ServiceGroupApp.ExtensionPHPServiceApp.LoadStatus(param.Version)
	if err != nil {
		global.Log.Errorf("LoadStatus->ServiceGroupApp.ExtensionPHPServiceApp.LoadStatus:%v \n", err)
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// FpmLog
// @Tags     PHP
// @Summary   获取php-fpm日志
// @Router    /extension_php/FpmLog [post]
func (s *ExtensionPHPApi) FpmLog(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.FpmLog.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取php-fpm日志
	data, err := ServiceGroupApp.ExtensionPHPServiceApp.FpmLog(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// FpmSlowLog
// @Tags     PHP
// @Summary   获取php-fpm慢日志
// @Router    /extension_php/FpmSlowLog [post]
func (s *ExtensionPHPApi) FpmSlowLog(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.FpmSlowLog.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取php-fpm慢日志
	data, err := ServiceGroupApp.ExtensionPHPServiceApp.FpmSlowLog(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// PHPInfo
// @Tags     PHP
// @Summary   获取php信息
// @Router    /extension_php/PHPInfo [post]
func (s *ExtensionPHPApi) PHPInfo(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.PHPInfo.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取php信息
	data, err := ServiceGroupApp.ExtensionPHPServiceApp.PHPInfo(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// PHPInfoHtml
// @Tags     PHP
// @Summary   获取php信息html
// @Router    /extension_php/PHPInfoHtml [post]
func (s *ExtensionPHPApi) PHPInfoHtml(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.PHPInfoHtml.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取php信息html
	data, err := ServiceGroupApp.ExtensionPHPServiceApp.PHPInfoHtml(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// CmdVersion
// @Tags     PHP
// @Summary   获取php命令行版本
// @Router    /extension_php/CmdVersion [post]
func (s *ExtensionPHPApi) CmdVersion(c *gin.Context) {
	response := app.NewResponse(c)
	version, versionList, err := ServiceGroupApp.ExtensionPHPServiceApp.CmdVersion()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(map[string]any{
		"version":     version,
		"versionList": versionList,
	})
}

// SetCmdVersion
// @Tags     PHP
// @Summary   设置php命令行版本
// @Router    /extension_php/SetCmdVersion [post]
func (s *ExtensionPHPApi) SetCmdVersion(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("php.SetCmdVersion.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.ExtensionPHPServiceApp.SetCmdVersion(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExtension, helper.MessageWithMap("php.SetCmdVersion", map[string]any{"Version": param.Version}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

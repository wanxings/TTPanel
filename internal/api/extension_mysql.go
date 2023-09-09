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

type ExtensionMysqlApi struct{}

// Info
// @Tags     Mysql
// @Summary   获取Mysql扩展的信息
// @Router    /extensions/mysql/Info [get]
func (s *ExtensionMysqlApi) Info(c *gin.Context) {
	response := app.NewResponse(c)
	//获取基本信息
	data, err := ServiceGroupApp.ExtensionMysqlServiceApp.Info()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// Install
// @Tags     Mysql
// @Summary   安装Mysql
// @Router    /extensions/mysql/Install [post]
func (s *ExtensionMysqlApi) Install(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Mysql.Install.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//安裝Mysql
	err := ServiceGroupApp.ExtensionMysqlServiceApp.Install(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Install", map[string]any{"Name": "mysql-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// Uninstall
// @Tags     Mysql
// @Summary   卸载Mysql
// @Router    /extensions/mysql/Uninstall [post]
func (s *ExtensionMysqlApi) Uninstall(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Mysql.Uninstall.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//卸载Mysql
	err := ServiceGroupApp.ExtensionMysqlServiceApp.Uninstall(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Uninstall", map[string]any{"Name": "mysql-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// SetStatus
// @Tags     Mysql
// @Summary   设置Mysql状态
// @Router    /extensions/mysql/SetStatus [post]
func (s *ExtensionMysqlApi) SetStatus(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.MysqlSetStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Mysql.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//设置Mysql状态
	err := ServiceGroupApp.ExtensionMysqlServiceApp.SetStatus(param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExtension, helper.MessageWithMap("extension.SetStatus", map[string]any{"Name": "mysql", "Status": param.Action}))
	response.ToResponseMsg(helper.Message("tips.SetStatusSuccess"))
}

// PerformanceConfig
// @Tags     Mysql
// @Summary   获取性能配置
// @Router    /extensions/mysql/PerformanceConfig [get]
func (s *ExtensionMysqlApi) PerformanceConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取性能配置
	data, err := ServiceGroupApp.ExtensionMysqlServiceApp.PerformanceConfig()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SavePerformanceConfig
// @Tags     Mysql
// @Summary   保存性能配置
// @Router    /extensions/mysql/SavePerformanceConfig [post]
func (s *ExtensionMysqlApi) SavePerformanceConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.MysqlSetPerformanceConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Mysql.SavePerformanceConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//保存性能配置
	err := ServiceGroupApp.ExtensionMysqlServiceApp.SavePerformanceConfig(param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDatabase, helper.Message("database.mysql.SavePerformanceConfig"))
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// LoadStatus
// @Tags     Mysql
// @Summary   获取负载状态
// @Router    /extensions/mysql/LoadStatus [get]
func (s *ExtensionMysqlApi) LoadStatus(c *gin.Context) {
	response := app.NewResponse(c)
	//获取性能配置
	data, err := ServiceGroupApp.ExtensionMysqlServiceApp.LoadStatus()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// ErrorLog
// @Tags     Mysql
// @Summary   获取错误日志
// @Router    /extensions/mysql/ErrorLog [get]
func (s *ExtensionMysqlApi) ErrorLog(c *gin.Context) {
	response := app.NewResponse(c)
	//获取错误日志
	data, err := ServiceGroupApp.ExtensionMysqlServiceApp.ErrorLog()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SlowLogs
// @Tags     Mysql
// @Summary   获取慢查询日志
// @Router    /extensions/mysql/SlowLogs [get]
func (s *ExtensionMysqlApi) SlowLogs(c *gin.Context) {
	response := app.NewResponse(c)
	//获取慢查询日志
	data, err := ServiceGroupApp.ExtensionMysqlServiceApp.SlowLogs()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

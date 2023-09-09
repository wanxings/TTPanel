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

type ExtensionNginxApi struct{}

// Info
// @Tags     Nginx
// @Summary   获取Nginx的信息
// @Router    /extensions/nginx/Info [get]
func (s *ExtensionNginxApi) Info(c *gin.Context) {
	response := app.NewResponse(c)
	//获取基本信息
	data, err := ServiceGroupApp.ExtensionNginxServiceApp.Info()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

//// GetStatus
//// @Tags     Nginx
//// @Summary   获取Nginx的状态
//// @Router    /extensions/nginx/GetStatus [get]
//func (s *ExtensionNginxApi) GetStatus(c *gin.Context) {
//	response := app.NewResponse(c)
//	ResponseData := make(map[string]interface{})
//	//ServiceGroupApp.ExtensionNginxServiceApp.Install()
//	response.ToResponse(&ResponseData)
//}

// Install
// @Tags     Nginx
// @Summary   安装Nginx
// @Router    /extensions/nginx/Install [get]
func (s *ExtensionNginxApi) Install(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("nginx.Install.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//安裝nginx
	err := ServiceGroupApp.ExtensionNginxServiceApp.Install(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Install", map[string]any{"Name": "Nginx-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// Uninstall
// @Tags     Nginx
// @Summary   卸载nginx
// @Router    /extensions/nginx/Uninstall [get]
func (s *ExtensionNginxApi) Uninstall(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("nginx.Uninstall.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//卸载nginx
	err := ServiceGroupApp.ExtensionNginxServiceApp.Uninstall(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Uninstall", map[string]any{"Name": "nginx-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// SetStatus
// @Tags     Nginx
// @Summary   设置Nginx的运行状态
// @Router    /extensions/nginx/SetStatus [get]
func (s *ExtensionNginxApi) SetStatus(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.NginxSetStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("nginx.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//设置nginx状态
	err := ServiceGroupApp.ExtensionNginxServiceApp.SetStatus(param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExtension, helper.MessageWithMap("extension.SetStatus", map[string]any{"Name": "nginx", "Status": param.Action}))
	response.ToResponseMsg(helper.Message("tips.SetStatusSuccess"))
}

// PerformanceConfig
// @Tags     Nginx
// @Summary   获取性能配置
// @Router    /extensions/nginx/PerformanceConfig [get]
func (s *ExtensionNginxApi) PerformanceConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取性能配置
	data, err := ServiceGroupApp.ExtensionNginxServiceApp.PerformanceConfig()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SavePerformanceConfig
// @Tags     Nginx
// @Summary   保存性能配置
// @Router    /extensions/nginx/SavePerformanceConfig [post]
func (s *ExtensionNginxApi) SavePerformanceConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.NginxSetPerformanceConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("nginx.SavePerformanceConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//保存性能配置
	err := ServiceGroupApp.ExtensionNginxServiceApp.SavePerformanceConfig(param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExtension, helper.Message("nginx.SavePerformanceConfig"))

	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// LoadStatus
// @Tags     Nginx
// @Summary   获取负载状态
// @Router    /extensions/nginx/LoadStatus [get]
func (s *ExtensionNginxApi) LoadStatus(c *gin.Context) {
	response := app.NewResponse(c)
	//获取负载状态
	data, err := ServiceGroupApp.ExtensionNginxServiceApp.LoadStatus()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

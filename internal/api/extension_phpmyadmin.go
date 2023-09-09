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

type ExtensionPhpmyadminApi struct{}

// Info
// @Tags     Phpmyadmin
// @Summary   获取Phpmyadmin的信息
// @Router    /extensions/phpmyadmin/Info [post]
func (s *ExtensionPhpmyadminApi) Info(c *gin.Context) {
	response := app.NewResponse(c)
	//获取基本信息
	data, err := ServiceGroupApp.ExtensionPhpmyadminServiceApp.Info()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// Install
// @Tags     Phpmyadmin
// @Summary   安装Phpmyadmin
// @Router    /extensions/phpmyadmin/Install [post]
func (s *ExtensionPhpmyadminApi) Install(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Phpmyadmin.Install.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//安裝Phpmyadmin
	err := ServiceGroupApp.ExtensionPhpmyadminServiceApp.Install(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Install", map[string]any{"Name": "phpmyadmin-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// Uninstall
// @Tags     Phpmyadmin
// @Summary   卸载Phpmyadmin
// @Router    /extensions/phpmyadmin/Uninstall [post]
func (s *ExtensionPhpmyadminApi) Uninstall(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Phpmyadmin.Uninstall.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//卸载Phpmyadmin
	err := ServiceGroupApp.ExtensionPhpmyadminServiceApp.Uninstall(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Uninstall", map[string]any{"Name": "phpmyadmin-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// GetConfig
// @Tags     Phpmyadmin
// @Summary   获取Phpmyadmin的配置
// @Router    /extensions/phpmyadmin/GetConfig [post]
func (s *ExtensionPhpmyadminApi) GetConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取配置
	data, err := ServiceGroupApp.ExtensionPhpmyadminServiceApp.GetConfig()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SetConfig
// @Tags     Phpmyadmin
// @Summary   设置Phpmyadmin的配置
// @Router    /extensions/phpmyadmin/SetConfig [post]
func (s *ExtensionPhpmyadminApi) SetConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Phpmyadmin.SetConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//设置配置
	err := ServiceGroupApp.ExtensionPhpmyadminServiceApp.SetConfig(param.Port, param.Php, false, false)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExtension, helper.MessageWithMap("phpmyadmin.SetConfig", map[string]any{"Version": param.Php, "Port": param.Port}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

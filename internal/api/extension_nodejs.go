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

type ExtensionNodejsApi struct{}

// Info
// @Tags     Nodejs
// @Summary   获取Nodejs管理器的信息
// @Router    /extensions/nodejs/Info [post]
func (s *ExtensionNodejsApi) Info(c *gin.Context) {
	response := app.NewResponse(c)
	//获取Nodejs管理器的信息
	data, err := ServiceGroupApp.ExtensionNodejsServiceApp.Info()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// Config
// @Tags     Nodejs
// @Summary   获取Nodejs管理器的配置
// @Router    /extensions/nodejs/Config [post]
func (s *ExtensionNodejsApi) Config(c *gin.Context) {
	response := app.NewResponse(c)
	//获取Nodejs管理器的配置
	data, err := ServiceGroupApp.ExtensionNodejsServiceApp.Config()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SetRegistrySources
// @Tags     Nodejs
// @Summary   设置Nodejs管理器的镜像源
// @Router    /extensions/nodejs/SetRegistrySources [post]
func (s *ExtensionNodejsApi) SetRegistrySources(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.NodejsSetRegistrySourcesR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Nodejs.SetRegistrySources.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//设置Nodejs管理器的镜像源
	err := ServiceGroupApp.ExtensionNodejsServiceApp.SetRegistrySources(param.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetVersionUrl
// @Tags     Nodejs
// @Summary   设置获取Nodejs版本地址
// @Router    /extensions/nodejs/SetVersionUrl [post]
func (s *ExtensionNodejsApi) SetVersionUrl(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.NodejsSetVersionUrlR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Nodejs.SetVersionUrl.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//设置Nodejs管理器的版本地址
	err := ServiceGroupApp.ExtensionNodejsServiceApp.SetVersionUrl(param.Url)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// VersionList
// @Tags     Nodejs
// @Summary   获取Nodejs的版本列表
// @Router    /extensions/nodejs/VersionList [post]
func (s *ExtensionNodejsApi) VersionList(c *gin.Context) {
	response := app.NewResponse(c)
	//获取Nodejs的版本列表
	data, err := ServiceGroupApp.ExtensionNodejsServiceApp.VersionList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// UpdateVersionList
// @Tags     Nodejs
// @Summary   更新Nodejs的版本列表
// @Router    /extensions/nodejs/UpdateVersionList [post]
func (s *ExtensionNodejsApi) UpdateVersionList(c *gin.Context) {
	response := app.NewResponse(c)
	//更新Nodejs的版本列表
	err := ServiceGroupApp.ExtensionNodejsServiceApp.UpdateVersionList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.updateSuccess"))
}

// Install
// @Tags     Nodejs
// @Summary   安装Nodejs
// @Router    /extensions/nodejs/Install [post]
func (s *ExtensionNodejsApi) Install(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.NodejsInstallR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Nodejs.Install.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//安装Nodejs
	err := ServiceGroupApp.ExtensionNodejsServiceApp.Install(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Install", map[string]any{"Name": "nodejs-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// Uninstall
// @Tags     Nodejs
// @Summary   卸载Nodejs
// @Router    /extensions/nodejs/Uninstall [post]
func (s *ExtensionNodejsApi) Uninstall(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.NodejsUninstallR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Nodejs.Uninstall.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//卸载Nodejs
	err := ServiceGroupApp.ExtensionNodejsServiceApp.Uninstall(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Uninstall", map[string]any{"Name": "nodejs-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// SetDefaultEnv
// @Tags     Nodejs
// @Summary   设置Nodejs的默认环境
// @Router    /extensions/nodejs/SetDefaultEnv [post]
func (s *ExtensionNodejsApi) SetDefaultEnv(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.NodejsSetDefaultEnvR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Nodejs.SetDefaultEnv.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//设置Nodejs的默认环境
	err := ServiceGroupApp.ExtensionNodejsServiceApp.SetDefaultEnv(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// NodeModulesList
// @Tags     Nodejs
// @Summary   获取Nodejs的模块列表
// @Router    /extensions/nodejs/NodeModulesList [post]
func (s *ExtensionNodejsApi) NodeModulesList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.NodejsNodeModulesListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Nodejs.NodeModulesList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//获取Nodejs的模块列表
	data, err := ServiceGroupApp.ExtensionNodejsServiceApp.NodeModulesList(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, 0, 0, 0)
}

// OperationNodeModules
// @Tags     Nodejs
// @Summary   操作Nodejs模块
// @Router    /extensions/nodejs/OperationNodeModules [post]
func (s *ExtensionNodejsApi) OperationNodeModules(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.NodejsOperationNodeModulesR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Nodejs.InstallNodeModules.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//安装Nodejs的模块
	err := ServiceGroupApp.ExtensionNodejsServiceApp.OperationNodeModules(param.Version, param.Modules, param.Operation)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExtension, helper.MessageWithMap("nodejs.OperationNodeModules", map[string]any{"Operation": param.Operation, "Modules": param.Modules, "Version": param.Version}))
	response.ToResponseMsg(helper.Message("tips.OperateSuccess"))
}

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

type LinuxToolsApi struct{}

// GetDnsConfig
// @Tags      System
// @Summary   获取dns配置
// @Router    /linux_tools/GetDnsConfig [post]
func (s *LinuxToolsApi) GetDnsConfig(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.LinuxToolsServiceApp.GetDnsConfig()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// SaveDnsConfig
// @Tags      System
// @Summary   保存dns配置
// @Router    /linux_tools/SaveDnsConfig [post]
func (s *LinuxToolsApi) SaveDnsConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SaveDnsConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("SaveDnsConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.LinuxToolsServiceApp.SaveDnsConfig(param.DNSList)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByLinuxTools, helper.Message("linux_tools.SaveDnsConfig"))

	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// TestDnsConfig
// @Tags      System
// @Summary   测试dns配置
// @Router    /linux_tools/TestDnsConfig [post]
func (s *LinuxToolsApi) TestDnsConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SaveDnsConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("TestDnsConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.LinuxToolsServiceApp.TestDnsConfig(param.DNSList)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseMsg(helper.Message("tips.TestSuccess"))
}

// RecoverDnsConfig
// @Tags      System
// @Summary   恢复dns配置
// @Router    /linux_tools/RecoverDnsConfig [post]
func (s *LinuxToolsApi) RecoverDnsConfig(c *gin.Context) {
	response := app.NewResponse(c)
	err := ServiceGroupApp.LinuxToolsServiceApp.RecoverDnsConfig()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByLinuxTools, helper.Message("linux_tools.RecoverDnsConfig"))

	response.ToResponseMsg(helper.Message("tips.RecoverySuccess"))
}

// GetTimeZoneConfig
// @Tags      System
// @Summary   获取时区配置
// @Router    /linux_tools/GetTimeZoneConfig [post]
func (s *LinuxToolsApi) GetTimeZoneConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.GetTimeZoneConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("GetTimeZoneConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.LinuxToolsServiceApp.GetTimeZoneConfig(param.MainZone)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// SetTimeZone
// @Tags      System
// @Summary   设置时区配置
// @Router    /linux_tools/SetTimeZone [post]
func (s *LinuxToolsApi) SetTimeZone(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SetTimeZoneR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("SetTimeZone.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.LinuxToolsServiceApp.SetTimeZone(param.MainZone, param.SubZone)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByLinuxTools, helper.MessageWithMap("linux_tools.SetTimeZone", map[string]interface{}{"MainZone": param.MainZone, "SubZone": param.SubZone}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SyncDate
// @Tags     SyncDate
// @Summary   同步服务器时间
// @Router    /linux_tools/SyncDate [post]
func (s *LinuxToolsApi) SyncDate(c *gin.Context) {
	response := app.NewResponse(c)
	err := ServiceGroupApp.LinuxToolsServiceApp.SyncDate()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByLinuxTools, helper.Message("linux_tools.SyncDate"))
	response.ToResponseMsg(helper.Message("tips.SyncSuccess"))
}

// GetNetworkConfig
// @Tags     Network
// @Summary   获取网卡配置
// @Router    /linux_tools/GetNetworkConfig [post]
func (s *LinuxToolsApi) GetNetworkConfig(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.LinuxToolsServiceApp.GetNetworkConfig()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// GetHostsConfig
// @Tags     Host
// @Summary   获取host配置
// @Router    /linux_tools/GetHostsConfig [post]
func (s *LinuxToolsApi) GetHostsConfig(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.LinuxToolsServiceApp.GetHostsConfig()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// AddHosts
// @Tags     Host
// @Summary   添加host
// @Router    /linux_tools/AddHosts [post]
func (s *LinuxToolsApi) AddHosts(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.AddHostsR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("AddHost.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.LinuxToolsServiceApp.AddHosts(param.Domain, param.IP)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByLinuxTools, helper.MessageWithMap("linux_tools.AddHosts", map[string]interface{}{"Hosts": param.IP + " " + param.Domain}))
	response.ToResponseMsg(helper.Message("tips.AddSuccess"))
}

// RemoveHosts
// @Tags     Host
// @Summary   删除host
// @Router    /linux_tools/RemoveHosts [post]
func (s *LinuxToolsApi) RemoveHosts(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.AddHostsR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("RemoveHosts.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.LinuxToolsServiceApp.RemoveHosts(param.Domain)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByLinuxTools, helper.MessageWithMap("linux_tools.RemoveHosts", map[string]interface{}{"Hosts": param.IP + " " + param.Domain}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

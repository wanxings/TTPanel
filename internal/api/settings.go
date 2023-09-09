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

type SettingsApi struct{}

// List
// @Tags     Settings
// @Summary   获取基础信息
// @Router    /api/settings/List [post]
func (s *SettingsApi) List(c *gin.Context) {
	response := app.NewResponse(c)
	ResponseData, err := ServiceGroupApp.SettingsServiceApp.List()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(&ResponseData)
}

// SetBasicAuth
// @Tags     Settings
// @Summary   设置基础认证
// @Router    /api/settings/SetBasicAuth [post]
func (s *SettingsApi) SetBasicAuth(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetBasicAuthR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetBasicAuth.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetBasicAuth(param.Status, param.Username, param.Password)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySettings, helper.MessageWithMap("settings.SetBasicAuth", map[string]any{"Status": param.Status}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetPanelPort
// @Tags     Settings
// @Summary   设置面板端口
// @Router    /api/settings/SetPanelPort [post]
func (s *SettingsApi) SetPanelPort(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetPanelPortR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetPanelPort.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetPanelPort(param.PanelPort)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//为了保证面板的访问，直接添加面板的端口
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	FirewallInfo := firewall.FirewallStatus()
	if FirewallInfo.Status {
		rules := []*request.CreatePortRuleR{
			{
				Port:     param.PanelPort,
				Strategy: constant.SystemFirewallStrategyAllow,
				Ps:       "面板端口",
				Protocol: "tcp",
			},
		}
		err = firewall.BatchCreatePortRule(rules)
		if err != nil {
			global.Log.Errorf("settings->SetPanelPort->BatchCreatePortRule Error:%s", err.Error())
		}
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySettings, helper.MessageWithMap("settings.SetPanelPort", map[string]any{"Port": param.PanelPort}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetEntrance
// @Tags     Settings
// @Summary   设置入口
// @Router    /api/settings/SetEntrance [post]
func (s *SettingsApi) SetEntrance(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetEntranceR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetEntrance.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetEntrance(param.Entrance)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySettings, helper.MessageWithMap("settings.SetEntrance", map[string]any{"Entrance": param.Entrance}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetEntranceErrorCode
// @Tags     Settings
// @Summary   设置入口错误码
// @Router    /api/settings/SetEntranceErrorCode [post]
func (s *SettingsApi) SetEntranceErrorCode(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetEntranceErrorCodeR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetEntranceErrorCode.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetEntranceErrorCode(param.ErrorCode)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySettings, helper.MessageWithMap("settings.SetEntranceErrorCode", map[string]any{"ErrorCode": param.ErrorCode}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetUser
// @Tags     Settings
// @Summary   设置用户
// @Router    /api/settings/SetUser [post]
func (s *SettingsApi) SetUser(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetUserR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetUser.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetUser(param.Username, param.Password)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySettings, helper.MessageWithMap("settings.ChangeUser", map[string]any{"User": param.Username}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetPanelName
// @Tags     Settings
// @Summary   设置面板名称
// @Router    /api/settings/SetPanelName [post]
func (s *SettingsApi) SetPanelName(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetPanelNameR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetPanelName.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetPanelName(param.PanelName)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySettings, helper.MessageWithMap("settings.SetPanelName", map[string]any{"Name": param.PanelName}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetPanelIP
// @Tags     Settings
// @Summary   设置面板IP
// @Router    /api/settings/SetPanelIP [post]
func (s *SettingsApi) SetPanelIP(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetPanelIPR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetPanelIP.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetPanelIP(param.PanelIP)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySettings, helper.MessageWithMap("settings.SetPanelIP", map[string]any{"IP": param.PanelIP}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetDefaultWebsiteDirectory
// @Tags     Settings
// @Summary   设置默认网站目录
// @Router    /api/settings/SetDefaultWebsiteDirectory [post]
func (s *SettingsApi) SetDefaultWebsiteDirectory(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetDefaultWebsiteDirectoryR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetDefaultWebsiteDirectory.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetDefaultWebsiteDirectory(param.DefaultWebsiteDirectory)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetDefaultBackupDirectory
// @Tags     Settings
// @Summary   设置默认备份目录
// @Router    /api/settings/SetDefaultBackupDirectory [post]
func (s *SettingsApi) SetDefaultBackupDirectory(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetDefaultBackupDirectoryR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetDefaultBackupDirectory.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetDefaultBackupDirectory(param.DefaultBackupDirectory)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetPanelApi
// @Tags     Settings
// @Summary   设置面板API
// @Router    /api/settings/SetPanelApi [post]
func (s *SettingsApi) SetPanelApi(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetPanelApiR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetPanelApi.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetPanelApi(param.Status, param.Key, param.Whitelist)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySettings, helper.MessageWithMap("settings.SetPanelApi", map[string]any{"Status": param.Status}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetLanguage
// @Tags     Settings
// @Summary   设置面板语言
// @Router    /api/settings/SetLanguage [post]
func (s *SettingsApi) SetLanguage(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetLanguageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("settings.SetLanguage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.SettingsServiceApp.SetLanguage(param.Language)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

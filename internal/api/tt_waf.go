package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model/request"
	ttwafResponse "TTPanel/internal/model/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type TTWafApi struct{}

// Config
// @Tags      tt_waf
// @Summary   配置
// @Router    ttwaf/Config [post]
func (m *TTWafApi) Config(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.TTWafServiceApp.Config()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// ProjectConfig
// @Tags      tt_waf
// @Summary   项目配置
// @Router    ttwaf/ProjectConfig [post]
func (m *TTWafApi) ProjectConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TTWafProjectConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("tt_waf.GlobalSet.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.TTWafServiceApp.ProjectConfig(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = ServiceGroupApp.ExtensionNginxServiceApp.SetStatus(constant.ProcessCommandByRestart)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// GlobalSet
// @Tags      tt_waf
// @Summary   全局设置
// @Router    ttwaf/GlobalSet [post]
func (m *TTWafApi) GlobalSet(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TTWafGlobalSetR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("tt_waf.GlobalSet.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.TTWafServiceApp.GlobalSet(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = ServiceGroupApp.ExtensionNginxServiceApp.SetStatus(constant.ProcessCommandByRestart)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByTTWaf, helper.Message("ttwaf.GlobalSet"))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// CountryList
// @Tags      tt_waf
// @Summary   城市列表
// @Router    ttwaf/CountryList [post]
func (m *TTWafApi) CountryList(c *gin.Context) {
	response := app.NewResponse(c)
	data := ServiceGroupApp.TTWafServiceApp.CountryList()
	response.ToResponse(data)
}

// SaveConfig
// @Tags      tt_waf
// @Summary   保存配置
// @Router    ttwaf/SaveConfig [post]
func (m *TTWafApi) SaveConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := ttwafResponse.TTWafConfig{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("tt_waf.SaveConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.TTWafServiceApp.SaveConfig(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByTTWaf, helper.Message("ttwaf.SaveConfig"))
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// SaveProjectConfig
// @Tags      tt_waf
// @Summary   保存项目配置
// @Router    ttwaf/SaveProjectConfig [post]
func (m *TTWafApi) SaveProjectConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TTWafSaveProjectConfig{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("tt_waf.GlobalSet.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.TTWafServiceApp.SaveProjectConfig(param.ProjectID, param.TTWafProjectConfig)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByTTWaf, helper.Message("ttwaf.SaveProjectConfig"))
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// BlockList
// @Tags      tt_waf
// @Summary   拦截列表
// @Router    ttwaf/BlockList [post]
func (m *TTWafApi) BlockList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TTWafBlockListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("go_project.ProjectList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	data, total, err := ServiceGroupApp.TTWafServiceApp.BlockList(&param, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// BanList
// @Tags      tt_waf
// @Summary   封禁列表
// @Router    ttwaf/BanList [post]
func (m *TTWafApi) BanList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TTWafBanListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("go_project.ProjectList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	data, total, err := ServiceGroupApp.TTWafServiceApp.BanList(&param, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// Overview
// @Tags      tt_waf
// @Summary   概览
// @Router    ttwaf/Overview [post]
func (m *TTWafApi) Overview(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.TTWafServiceApp.Overview()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// GetRegRule
// @Tags      tt_waf
// @Summary   获取正则规则
// @Router    ttwaf/GetRegRule [post]
func (m *TTWafApi) GetRegRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TTWafGetRegRuleR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("tt_waf.GetRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.TTWafServiceApp.GetRegRule(param.RuleName)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SaveRegRule
// @Tags      tt_waf
// @Summary   保存正则规则
// @Router    ttwaf/SaveRegRule [post]
func (m *TTWafApi) SaveRegRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TTWafSaveRegRuleR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("tt_waf.SaveRegRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.TTWafServiceApp.SaveRegRule(param.RuleName, param.RuleList)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByTTWaf, helper.Message("ttwaf.SaveRegRule"))
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// AllowIP
// @Tags      tt_waf
// @Summary   放行IP
// @Router    ttwaf/AllowIP [post]
func (m *TTWafApi) AllowIP(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TTWafIPR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("tt_waf.OperateIP.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	var errList []string

	for _, ip := range param.Ips {
		err := ServiceGroupApp.TTWafServiceApp.AllowIP(ip)
		if err != nil {
			errList = append(errList, fmt.Sprintf("IP:%s,ERROR:%s", ip, err.Error()))
		}
	}

	if len(errList) != 0 {
		response.ToErrorResponse(errcode.ServerError.WithDetails(errList...))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByTTWaf, helper.MessageWithMap("ttwaf.AllowIP", map[string]any{"IPS": strings.Join(param.Ips, " ")}))
	response.ToResponseMsg(helper.Message("tips.OperateSuccess"))
}

// AddIpBlackList
// @Tags      tt_waf
// @Summary   添加IP黑名单
// @Router    ttwaf/AddIpBlackList [post]
func (m *TTWafApi) AddIpBlackList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.TTWafIPBlackWhiteR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("tt_waf.AddIpBlackList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.TTWafServiceApp.AddIpBlackList(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByTTWaf, helper.MessageWithMap("ttwaf.AddIpBlackList", map[string]any{"Ipv4": param.IPV4, "Ipv6": param.IPV6}))
	response.ToResponseMsg(helper.Message("tips.AddSuccess"))
}

// AnalyticsOverview
// @Tags      tt_waf
// @Summary   统计分析概览
// @Router    ttwaf/AnalyticsOverview [post]
func (m *TTWafApi) AnalyticsOverview(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.AnalyticsOverviewR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("tt_waf.AnalyticsOverview.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectGet, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//获取项目统计分析信息
	data, err := ServiceGroupApp.TTWafServiceApp.AnalyticsOverview(projectGet.Name, param.StartTime, param.EndTime)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

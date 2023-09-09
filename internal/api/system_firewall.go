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

type SystemFirewallApi struct{}

// FirewallStatus
// @Tags      firewall
// @Summary   系统防火墙信息
// @Router    /firewall/FirewallStatus [get]
func (f *SystemFirewallApi) FirewallStatus(c *gin.Context) {
	response := app.NewResponse(c)
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	data := firewall.FirewallStatus()
	response.ToResponse(data)
}

// BatchCreatePortRule
// @Tags      firewall
// @Summary   创建端口规则-批量
// @Router    /firewall/BatchCreatePortRule [post]
func (f *SystemFirewallApi) BatchCreatePortRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchCreatePortRuleR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchCreatePortRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.BatchCreatePortRule(param.Rules)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.BatchCreatePortRule"))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// BatchDeletePortRule
// @Tags      firewall
// @Summary   删除端口规则-批量
// @Router    /firewall/BatchDeletePortRule [post]
func (f *SystemFirewallApi) BatchDeletePortRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchIDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchDeletePortRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.BatchDeletePortRule(param.IDs)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.BatchDeletePortRule"))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// PortRuleList
// @Tags      firewall
// @Summary   端口规则列表
// @Router    /firewall/PortRuleList [get]
func (f *SystemFirewallApi) PortRuleList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.FirewallListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("PortRuleList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	data, total, err := firewall.GetPortRules(param.Query, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// UpdatePortRule
// @Tags      firewall
// @Summary   更新端口规则
// @Router    /firewall/UpdatePortRule [post]
func (f *SystemFirewallApi) UpdatePortRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.UpdatePortRuleR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("UpdatePortRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.UpdatePortRule(param.ID, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.UpdatePortRule"))
	response.ToResponseMsg(helper.Message("tips.updateSuccess"))
}

// Close 关闭防火墙
// @Tags      firewall
// @Summary   关闭防火墙
// @Router    /firewall/Close [post]
func (f *SystemFirewallApi) Close(c *gin.Context) {
	response := app.NewResponse(c)
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.Close()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.Close"))
	response.ToResponseMsg(helper.Message("tips.CloseSuccess"))
}

// Open 开启防火墙
// @Tags      firewall
// @Summary   开启防火墙
// @Router    /firewall/Open [post]
func (f *SystemFirewallApi) Open(c *gin.Context) {
	response := app.NewResponse(c)
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//打开防火墙
	err = firewall.Open()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	//为了保证面板的访问，直接添加面板的端口
	rules := []*request.CreatePortRuleR{
		{
			Port:     global.Config.System.PanelPort,
			Strategy: constant.SystemFirewallStrategyAllow,
			Ps:       "面板端口",
			Protocol: "tcp",
		},
	}
	err = firewall.BatchCreatePortRule(rules)
	if err != nil {
		global.Log.Errorf("system_firewall->Open->BatchCreatePortRule Error:%s", err.Error())
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.Open"))
	response.ToResponseMsg(helper.Message("tips.OpenSuccess"))
}

// AllowPing 允许ping
// @Tags      firewall
// @Summary   允许ping
// @Router    /firewall/AllowPing [post]
func (f *SystemFirewallApi) AllowPing(c *gin.Context) {
	response := app.NewResponse(c)
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.AllowPing()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.AllowPing"))
	response.ToResponseMsg(helper.Message("tips.OperateSuccess"))
}

// DenyPing 禁止ping
// @Tags      firewall
// @Summary   禁止ping
// @Router    /firewall/DenyPing [post]
func (f *SystemFirewallApi) DenyPing(c *gin.Context) {
	response := app.NewResponse(c)
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.DenyPing()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.DenyPing"))
	response.ToResponseMsg(helper.Message("tips.OperateSuccess"))
}

// BatchCreateIPRule 批量创建IP规则
// @Tags      firewall
// @Summary   批量创建IP规则
// @Router    /firewall/BatchCreateIPRule [post]
func (f *SystemFirewallApi) BatchCreateIPRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchCreateIPRuleR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchCreateIPRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.BatchCreateIPRule(param.Rules)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.BatchCreateIPRule"))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// BatchDeleteIPRule 批量删除IP规则
// @Tags      firewall
// @Summary   批量删除IP规则
// @Router    /firewall/BatchDeleteIPRule [post]
func (f *SystemFirewallApi) BatchDeleteIPRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchDeleteIPRuleR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchDeleteIPRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.BatchDeleteIPRule(param.Ids)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.BatchDeleteIPRule"))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// IPRuleList
// @Tags      firewall
// @Summary   IP规则列表
// @Router    /firewall/IPRuleList [get]
func (f *SystemFirewallApi) IPRuleList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.FirewallListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Edit.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	data, total, err := firewall.GetIpRules(param.Query, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// UpdateIPRule 更新IP规则
// @Tags      firewall
// @Summary   更新IP规则
// @Router    /firewall/UpdateIPRule [post]
func (f *SystemFirewallApi) UpdateIPRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.UpdateIPRuleR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("UpdateIPRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.UpdateIPRule(param.ID, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.UpdateIPRule"))
	response.ToResponseMsg(helper.Message("tips.updateSuccess"))
}

// BatchCreateForwardRule 批量创建转发规则
// @Tags      firewall
// @Summary   批量创建转发规则
// @Router    /firewall/BatchCreateForwardRule [post]
func (f *SystemFirewallApi) BatchCreateForwardRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchCreateForwardRuleR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchCreateForwardRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.BatchCreateForwardRule(param.Rules)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.BatchCreateForwardRule"))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// BatchDeleteForwardRule 批量删除转发规则
// @Tags      firewall
// @Summary   批量删除转发规则
// @Router    /firewall/BatchDeleteForwardRule [post]
func (f *SystemFirewallApi) BatchDeleteForwardRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchDeleteForwardRuleR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchDeleteForwardRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.BatchDeleteForwardRule(param.Ids)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.BatchDeleteForwardRule"))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// ForwardRuleList
// @Tags      firewall
// @Summary   转发规则列表
// @Router    /firewall/ForwardRuleList [get]
func (f *SystemFirewallApi) ForwardRuleList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.FirewallListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Edit.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	data, total, err := firewall.GetForwardRules(param.Query, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// UpdateForwardRule 更新转发规则
// @Tags      firewall
// @Summary   更新转发规则
// @Router    /firewall/UpdateForwardRule [post]
func (f *SystemFirewallApi) UpdateForwardRule(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.UpdateForwardRuleR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("UpdateForwardRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = firewall.UpdateForwardRule(param.ID, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeBySystemFirewall, helper.Message("firewalld.UpdateForwardRule"))

	response.ToResponseMsg(helper.Message("tips.updateSuccess"))
}

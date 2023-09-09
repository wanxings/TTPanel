package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type SystemFirewallRouter struct{}

func (f *SystemFirewallRouter) Init(Router *gin.RouterGroup) {
	systemFirewallRouter := Router.Group("system_firewall")
	systemFirewallApi := api.GroupApp.SystemFirewallApiApp
	{
		systemFirewallRouter.POST("FirewallStatus", systemFirewallApi.FirewallStatus) // 获取系统防火墙信息
		systemFirewallRouter.POST("Close", systemFirewallApi.Close)                   //  关闭防火墙
		systemFirewallRouter.POST("Open", systemFirewallApi.Open)                     // 开启防火墙
		systemFirewallRouter.POST("AllowPing", systemFirewallApi.AllowPing)           // 允许ping
		systemFirewallRouter.POST("DenyPing", systemFirewallApi.DenyPing)             // 禁止ping
	}
	{
		systemFirewallRouter.POST("BatchCreatePortRule", systemFirewallApi.BatchCreatePortRule) // 创建端口规则-批量
		systemFirewallRouter.POST("PortRuleList", systemFirewallApi.PortRuleList)               // 端口规则列表
		systemFirewallRouter.POST("UpdatePortRule", systemFirewallApi.UpdatePortRule)           // 编辑端口规则
		systemFirewallRouter.POST("BatchDeletePortRule", systemFirewallApi.BatchDeletePortRule) // 删除端口规则-批量
	}
	{
		systemFirewallRouter.POST("BatchCreateIPRule", systemFirewallApi.BatchCreateIPRule) // 创建IP规则-批量
		systemFirewallRouter.POST("IPRuleList", systemFirewallApi.IPRuleList)               // IP规则列表
		systemFirewallRouter.POST("UpdateIPRule", systemFirewallApi.UpdateIPRule)           // 更新IP规则
		systemFirewallRouter.POST("BatchDeleteIPRule", systemFirewallApi.BatchDeleteIPRule) // 删除IP规则-批量
	}
	{
		systemFirewallRouter.POST("BatchCreateForwardRule", systemFirewallApi.BatchCreateForwardRule) // 创建转发规则-批量
		systemFirewallRouter.POST("ForwardRuleList", systemFirewallApi.ForwardRuleList)               // 转发规则列表
		systemFirewallRouter.POST("UpdateForwardRule", systemFirewallApi.UpdateForwardRule)           // 更新转发规则
		systemFirewallRouter.POST("BatchDeleteForwardRule", systemFirewallApi.BatchDeleteForwardRule) // 删除转发规则-批量
	}
}

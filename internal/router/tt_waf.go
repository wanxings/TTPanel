package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type TTWafRouter struct{}

func (m *TTWafRouter) Init(Router *gin.RouterGroup) {
	ttWafRouter := Router.Group("tt_waf")
	ttWafApi := api.GroupApp.TTWafApiApp
	{ //读操作
		ttWafRouter.POST("Config", ttWafApi.Config)               //信息
		ttWafRouter.POST("CountryList", ttWafApi.CountryList)     //国家列表
		ttWafRouter.POST("BlockList", ttWafApi.BlockList)         //阻止列表
		ttWafRouter.POST("BanList", ttWafApi.BanList)             //封禁列表
		ttWafRouter.POST("Overview", ttWafApi.Overview)           //概览
		ttWafRouter.POST("GetRegRule", ttWafApi.GetRegRule)       //获取规则
		ttWafRouter.POST("SaveRegRule", ttWafApi.SaveRegRule)     //保存规则
		ttWafRouter.POST("ProjectConfig", ttWafApi.ProjectConfig) //项目防火墙配置
	}
	{ //写操作,写操作需要防火墙是开启状态
		ttWafRouter.POST("GlobalSet", ttWafApi.GlobalSet)                 //全局应用cc和容忍攻击
		ttWafRouter.POST("SaveConfig", ttWafApi.SaveConfig)               //保存配置
		ttWafRouter.POST("AllowIP", ttWafApi.AllowIP)                     //解封IP
		ttWafRouter.POST("AddIpBlackList", ttWafApi.AddIpBlackList)       //添加IP黑名单
		ttWafRouter.POST("SaveProjectConfig", ttWafApi.SaveProjectConfig) //保存项目防火墙配置
	}
}

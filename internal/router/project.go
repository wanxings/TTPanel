package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ProjectRouter struct{}

func (s *ProjectRouter) Init(Router *gin.RouterGroup) {
	projectRouter := Router.Group("project")
	projectApi := api.GroupApp.ProjectApiApp
	{
		projectRouter.POST("CategoryList", projectApi.CategoryList)     // 获取分类
		projectRouter.POST("CreateCategory", projectApi.CreateCategory) // 创建分类
		projectRouter.POST("EditCategory", projectApi.EditCategory)     // 编辑分类
		projectRouter.POST("SetPs", projectApi.SetPs)                   // 设置项目备注
	}
	{
		projectRouter.POST("AddDomains", projectApi.AddDomains)               // 添加域名
		projectRouter.POST("DomainList", projectApi.DomainList)               // 域名列表
		projectRouter.POST("BatchDeleteDomain", projectApi.BatchDeleteDomain) // 批量删除域名
	}
	{
		projectRouter.POST("DefaultIndex", projectApi.DefaultIndex)         // 默认首页
		projectRouter.POST("SaveDefaultIndex", projectApi.SaveDefaultIndex) // 保存默认首页
	}
	{
		projectRouter.POST("RewriteTemplateList", projectApi.RewriteTemplateList) // 伪静态模板列表
	}
	{
		projectRouter.POST("SetSSL", projectApi.SetSSL)                 // 设置SSL
		projectRouter.POST("CloseSSL", projectApi.CloseSSL)             // 关闭SSL
		projectRouter.POST("AlwaysUseHttps", projectApi.AlwaysUseHttps) // 始终使用ssl
	}
	{
		projectRouter.POST("CreateRedirect", projectApi.CreateRedirect)           // 创建重定向
		projectRouter.POST("RedirectList", projectApi.RedirectList)               // 重定向列表
		projectRouter.POST("BatchEditRedirect", projectApi.BatchEditRedirect)     // 批量编辑重定向
		projectRouter.POST("BatchDeleteRedirect", projectApi.BatchDeleteRedirect) // 批量删除重定向

	}
	{
		projectRouter.POST("CreateAntiLeechConfig", projectApi.CreateAntiLeechConfig) // 创建防盗链
		projectRouter.POST("GetAntiLeechConfig", projectApi.GetAntiLeechConfig)       // 获取防盗链
		projectRouter.POST("CloseAntiLeech", projectApi.CloseAntiLeech)               // 关闭防盗链
	}
	{
		projectRouter.POST("CreateReverseProxyConfig", projectApi.CreateReverseProxyConfig)           // 创建反代
		projectRouter.POST("ReverseProxyList", projectApi.ReverseProxyList)                           // 反代列表
		projectRouter.POST("BatchEditReverseProxyConfig", projectApi.BatchEditReverseProxyConfig)     // 批量编辑反代
		projectRouter.POST("BatchDeleteReverseProxyConfig", projectApi.BatchDeleteReverseProxyConfig) // 批量删除反代
	}
	{
		projectRouter.POST("CreateAccessRuleConfig", projectApi.CreateAccessRuleConfig)           // 创建访问规则
		projectRouter.POST("AccessRuleConfigList", projectApi.AccessRuleConfigList)               // 访问规则列表
		projectRouter.POST("EditAccessRuleConfig", projectApi.EditAccessRuleConfig)               // 编辑访问规则
		projectRouter.POST("BatchDeleteAccessRuleConfig", projectApi.BatchDeleteAccessRuleConfig) // 批量删除访问规则
	}
}

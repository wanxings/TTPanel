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

type ProjectApi struct{}

// AddDomains
// @Tags      AddDomains
// @Summary   添加域名
// @Router    /project/AddDomains [post]
func (s *ProjectApi) AddDomains(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.AddDomainsR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.AddDomains.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//添加域名
	failedErrs := ServiceGroupApp.ProjectServiceApp.AddDomains(projectData, &param)
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.AddDomains", map[string]any{"Name": projectData.Name, "Count": len(param.Domains)}))
	if failedErrs != nil {
		var errS []string
		for _, err := range failedErrs {
			errS = append(errS, err.Error())
		}
		response.ToErrorResponse(errcode.ServerError.WithDetails(errS...))
		return
	}

	response.ToResponseMsg(helper.Message("tips.AddSuccess"))
}

// DomainList
// @Tags      DomainList
// @Summary   域名列表
// @Router    /project/DomainList [post]
func (s *ProjectApi) DomainList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.IDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.DomainList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	data, total, err := ServiceGroupApp.ProjectServiceApp.DomainList(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), 0, 0)
}

// BatchDeleteDomain
// @Tags      BatchDeleteDomain
// @Summary   删除域名
// @Router    /project/BatchDeleteDomain [post]
func (s *ProjectApi) BatchDeleteDomain(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchDeleteDomainR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.BatchDeleteDomain.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//删除域名
	failedErrs := ServiceGroupApp.ProjectServiceApp.BatchDeleteDomain(projectData, param.IDS)
	if failedErrs != nil {
		var errS []string
		for _, err := range failedErrs {
			errS = append(errS, err.Error())
		}
		response.ToErrorResponse(errcode.ServerError.WithDetails(errS...))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.DeleteDomains", map[string]any{"Name": projectData.Name, "Count": len(param.IDS)}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// RewriteTemplateList
// @Tags      RewriteTemplateList
// @Summary   伪静态模板列表
// @Router    /project/RewriteTemplateList [post]
func (s *ProjectApi) RewriteTemplateList(c *gin.Context) {
	response := app.NewResponse(c)

	data, err := ServiceGroupApp.ProjectServiceApp.RewriteTemplateList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// DefaultIndex
// @Tags      DefaultIndex
// @Summary   默认首页
// @Router    /project/DefaultIndex [post]
func (s *ProjectApi) DefaultIndex(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.IDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.DefaultIndex.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//获取默认首页配置
	data, err := ServiceGroupApp.ProjectServiceApp.DefaultIndex(projectData.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SaveDefaultIndex
// @Tags      SaveDefaultIndex
// @Summary   保存默认首页
// @Router    /project/SaveDefaultIndex [post]
func (s *ProjectApi) SaveDefaultIndex(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SaveDefaultIndexR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.SaveDefaultIndex.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//保存默认首页
	err = ServiceGroupApp.ProjectServiceApp.SaveDefaultIndex(projectData.Name, param.Index)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.SaveDefaultIndex", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// CategoryList
// @Tags      CategoryList
// @Summary   分类列表
// @Router    /project/CategoryList [post]
func (s *ProjectApi) CategoryList(c *gin.Context) {
	response := app.NewResponse(c)
	//获取分类
	data, total, err := ServiceGroupApp.ProjectServiceApp.CategoryList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), 0, 0)
}

// CreateCategory
// @Tags      CreateCategory
// @Summary   创建分类
// @Router    /project/CreateCategory [post]
func (s *ProjectApi) CreateCategory(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.CreateCategoryR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.CreateCategory.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//添加分类
	err := ServiceGroupApp.ProjectServiceApp.CreateCategory(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.CreateCategory", map[string]any{"Name": param.Name}))

	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// EditCategory
// @Tags      EditCategory
// @Summary   编辑分类
// @Router    /project/EditCategory [post]
func (s *ProjectApi) EditCategory(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.EditCategoryR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.EditCategory.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//编辑分类
	err := ServiceGroupApp.ProjectServiceApp.EditCategory(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.EditCategory", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// SetPs
// @Tags      SetPs
// @Summary   设置项目备注
// @Router    /project/SetPs [post]
func (s *ProjectApi) SetPs(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetPsR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.SetPs.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	projectData.Ps = param.Ps
	err = projectData.Update(global.PanelDB)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.SetPs", map[string]any{"Name": projectData.Name, "Ps": param.Ps}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetSSL
// @Tags      SetSSL
// @Summary   设置SSL
// @Router    /project/SetSSL [post]
func (s *ProjectApi) SetSSL(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetSslR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.SetSSL.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//设置SSL
	err = ServiceGroupApp.ProjectServiceApp.SetSSL(projectData.Name, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.SetSSL", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// CloseSSL
// @Tags      CloseSSL
// @Summary   关闭SSL
// @Router    /project/CloseSSL [post]
func (s *ProjectApi) CloseSSL(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.IDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.CloseSSL.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//设置SSL
	err = ServiceGroupApp.ProjectServiceApp.CloseSSL(projectData.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.CloseSSL", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.CloseSuccess"))
}

// AlwaysUseHttps
// @Tags      AlwaysUseHttps
// @Summary   强制HTTPS
// @Router    /project/AlwaysUseHttps [post]
func (s *ProjectApi) AlwaysUseHttps(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.AlwaysUseHttpsR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.AlwaysUseHttps.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//设置强制HTTPS
	err = ServiceGroupApp.ProjectServiceApp.AlwaysUseHttps(projectData.Name, param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.AlwaysUseHttps", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// CreateRedirect
// @Tags      CreateRedirect
// @Summary   创建重定向
// @Router    /project/CreateRedirect [post]
func (s *ProjectApi) CreateRedirect(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.CreateRedirectR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.CreateRedirect.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//添加重定向
	err = ServiceGroupApp.ProjectServiceApp.CreateRedirect(projectData.Name, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.CreateRedirect", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// RedirectList
// @Tags      RedirectList
// @Summary   重定向列表
// @Router    /project/RedirectList [post]
func (s *ProjectApi) RedirectList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.IDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.RedirectList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//重定向列表
	data, err := ServiceGroupApp.ProjectServiceApp.RedirectList(projectData.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// BatchEditRedirect
// @Tags      BatchEditRedirect
// @Summary   批量编辑重定向
// @Router    /project/BatchEditRedirect [post]
func (s *ProjectApi) BatchEditRedirect(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchEditRedirectR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.BatchEditRedirect.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//批量编辑重定向
	err = ServiceGroupApp.ProjectServiceApp.BatchEditRedirect(projectData.Name, param.List)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.BatchEditRedirect", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// BatchDeleteRedirect
// @Tags      BatchDeleteRedirect
// @Summary   批量删除重定向
// @Router    /project/BatchDeleteRedirect [post]
func (s *ProjectApi) BatchDeleteRedirect(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchDeleteRedirectR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.BatchDeleteRedirect.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//批量删除重定向
	err = ServiceGroupApp.ProjectServiceApp.BatchDeleteRedirect(projectData.Name, param.Keys)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.BatchDeleteRedirect", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// GetAntiLeechConfig
// @Tags      GetAntiLeechConfig
// @Summary   获取防盗链
// @Router    /project/GetAntiLeechConfig [post]
func (s *ProjectApi) GetAntiLeechConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.IDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.GetAntiLeech.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//获取防盗链
	data, err := ServiceGroupApp.ProjectServiceApp.GetAntiLeechConfig(projectData)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// CreateAntiLeechConfig
// @Tags      CreateAntiLeechConfig
// @Summary   创建防盗链
// @Router    /project/CreateAntiLeechConfig [post]
func (s *ProjectApi) CreateAntiLeechConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.CreateAntiLeechConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.CreateAntiLeechConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	err = ServiceGroupApp.ProjectServiceApp.CreateAntiLeechConfig(projectData.Name, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.CreateAntiLeechConfig", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// CloseAntiLeech
// @Tags      CloseAntiLeech
// @Summary   关闭防盗链
// @Router    /project/CloseAntiLeech [post]
func (s *ProjectApi) CloseAntiLeech(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.IDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.CloseAntiLeech.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//关闭防盗链
	_, err = ServiceGroupApp.ProjectServiceApp.CloseAntiLeechConfig(projectData.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.CloseAntiLeech", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.CloseSuccess"))
}

// CreateReverseProxyConfig
// @Tags      CreateReverseProxyConfig
// @Summary   创建反向代理
// @Router    /project/CreateReverseProxyConfig [post]
func (s *ProjectApi) CreateReverseProxyConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.CreateReverseProxyConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.CreateReverseProxy.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//创建反向代理
	err = ServiceGroupApp.ProjectServiceApp.CreateReverseProxyConfig(projectData.Name, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.CreateReverseProxyConfig", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// ReverseProxyList
// @Tags      ReverseProxyList
// @Summary   反向代理列表
// @Router    /project/ReverseProxyList [post]
func (s *ProjectApi) ReverseProxyList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.IDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.ReverseProxyList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//获取反向代理列表
	data, err := ServiceGroupApp.ProjectServiceApp.ReverseProxyList(projectData.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// BatchEditReverseProxyConfig
// @Tags      BatchEditReverseProxyConfig
// @Summary   批量编辑反向代理
// @Router    /project/BatchEditReverseProxyConfig [post]
func (s *ProjectApi) BatchEditReverseProxyConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchEditReverseProxyConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.BatchEditReverseProxyConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	err = ServiceGroupApp.ProjectServiceApp.BatchEditReverseProxyConfig(projectData.Name, param.List)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.BatchEditReverseProxyConfig", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// BatchDeleteReverseProxyConfig
// @Tags      BatchDeleteReverseProxyConfig
// @Summary   批量删除反向代理
// @Router    /project/BatchDeleteReverseProxyConfig [post]
func (s *ProjectApi) BatchDeleteReverseProxyConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchDeleteReverseProxyConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.BatchDeleteReverseProxyConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//批量删除反向代理
	err = ServiceGroupApp.ProjectServiceApp.BatchDeleteReverseProxyConfig(projectData.Name, param.Keys)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.BatchDeleteReverseProxyConfig", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// CreateAccessRuleConfig
// @Tags      CreateAccessRuleConfig
// @Summary   创建访问控制
// @Router    /project/CreateAccessRuleConfig [post]
func (s *ProjectApi) CreateAccessRuleConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.CreateAccessRuleConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.CreateAccessRule.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//创建访问控制
	err = ServiceGroupApp.ProjectServiceApp.CreateAccessRuleConfig(projectData.Name, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.CreateAccessRuleConfig", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// AccessRuleConfigList
// @Tags      AccessRuleConfigList
// @Summary   访问控制列表
// @Router    /project/AccessRuleConfigList [post]
func (s *ProjectApi) AccessRuleConfigList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.IDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.AccessRuleConfigList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//获取访问控制列表
	data, err := ServiceGroupApp.ProjectServiceApp.AccessRuleConfigList(projectData.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// EditAccessRuleConfig
// @Tags      EditAccessRuleConfig
// @Summary   编辑访问控制
// @Router    /project/EditAccessRuleConfig [post]
func (s *ProjectApi) EditAccessRuleConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.EditAccessRuleConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.EditAccessRuleConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//编辑访问控制
	err = ServiceGroupApp.ProjectServiceApp.EditAccessRuleConfig(projectData.Name, param.Key, param.Config)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.EditAccessRuleConfig", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// BatchDeleteAccessRuleConfig
// @Tags      BatchDeleteAccessRuleConfig
// @Summary   批量删除访问控制
// @Router    /project/BatchDeleteAccessRuleConfig [post]
func (s *ProjectApi) BatchDeleteAccessRuleConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.BatchDeleteAccessRuleConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("project.BatchDeleteAccessRuleConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//批量删除访问控制
	err = ServiceGroupApp.ProjectServiceApp.BatchDeleteAccessRuleConfig(projectData.Name, param.Keys)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project.BatchDeleteAccessRuleConfig", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

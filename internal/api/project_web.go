package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	request2 "TTPanel/internal/model/request"
	"TTPanel/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ProjectPHPApi struct{}

// Create
// @Tags      Create
// @Summary   创建php项目
// @Router    /php_project/Create [post]
func (s *ProjectPHPApi) Create(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.CreatePHPProjectR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.Create.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//创建php项目
	err := ServiceGroupApp.ProjectPHPServiceApp.Create(&param)
	//
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_web.Create", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// List
// @Tags      List
// @Summary   项目列表
// @Router    /php_project/List [post]
func (s *ProjectPHPApi) List(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.PHPProjectListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.List.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)

	//获取项目列表
	data, total, err := ServiceGroupApp.ProjectPHPServiceApp.ProjectList(&param, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// ProjectInfo
// @Tags      ProjectInfo
// @Summary   项目详情
// @Router    /php_project/ProjectInfo [post]
func (s *ProjectPHPApi) ProjectInfo(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.IDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.ProjectInfo.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目详情
	data, err := ServiceGroupApp.ProjectPHPServiceApp.ProjectInfo(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SetStatus
// @Tags      SetStatus
// @Summary   设置项目状态
// @Router    /php_project/SetStatus [post]
func (s *ProjectPHPApi) SetStatus(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.SetPHPProjectStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//设置项目状态
	err = ServiceGroupApp.ProjectPHPServiceApp.SetStatus(projectData, param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_web.SetStatus", map[string]any{"Name": projectData.Name, "Status": param.Action}))
	response.ToResponseMsg(helper.Message("tips.SetStatusSuccess"))
}

// SetExpireTime
// @Tags      SetExpireTime
// @Summary   设置项目到期时间
// @Router    /php_project/SetExpireTime [post]
func (s *ProjectPHPApi) SetExpireTime(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.SetPHPProjectExpireTimeR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.SetExpireTime.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//设置项目到期时间
	err = ServiceGroupApp.ProjectPHPServiceApp.SetExpireTime(projectData, param.ExpireTime)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_web.SetExpireTime", map[string]any{"Name": projectData.Name, "ExpireTime": param.ExpireTime}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetRunPath
// @Tags      SetRunPath
// @Summary   设置项目运行路径
// @Router    /php_project/SetRunPath [post]
func (s *ProjectPHPApi) SetRunPath(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.SetPHPProjectRunPathR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.SetRunPath.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//设置项目运行路径
	err = ServiceGroupApp.ProjectPHPServiceApp.SetRunPath(projectData.Name, projectData.Path, param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_web.SetRunPath", map[string]any{"Name": projectData.Name, "RunPath": param.Path}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetPath
// @Tags      SetPath
// @Summary   设置项目路径
// @Router    /php_project/SetPath [post]
func (s *ProjectPHPApi) SetPath(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.SetPHPProjectPathR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.SetPath.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	if param.Path == "/" {
		param.Path = projectData.Path
	}
	if !util.PathExists(param.Path) {
		response.ToErrorResponse(errcode.ServerError.WithDetails(fmt.Sprintf("项目路径不存在: %s", param.Path)))
		return
	}
	//设置项目运行路径
	err = ServiceGroupApp.ProjectPHPServiceApp.SetRunPath(projectData.Name, projectData.Path, "/")
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	projectData.Path = param.Path
	err = projectData.Update(global.PanelDB)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_web.SetPath", map[string]any{"Name": projectData.Name, "Path": param.Path}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetUserIni
// @Tags      SetUserIni
// @Summary   设置项目用户配置
// @Router    /php_project/SetUserIni [post]
func (s *ProjectPHPApi) SetUserIni(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.SetPHPProjectUserIniR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.SetUserIni.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//设置项目用户配置
	err = ServiceGroupApp.ProjectPHPServiceApp.SetUserIni(projectData.Name, param.Status)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_web.SetUserIni", map[string]any{"Name": projectData.Name, "Status": param.Status}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// Delete
// @Tags      Delete
// @Summary   删除项目
// @Router    /php_project/Delete [post]
func (s *ProjectPHPApi) Delete(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DeletePHPProjectR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.Delete.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//删除项目
	err = ServiceGroupApp.ProjectPHPServiceApp.Delete(projectData, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_web.Delete", map[string]any{"Name": projectData.Name}))

	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// UsingPHPVersion
// @Tags      UsingPHPVersion
// @Summary   使用中的php版本
// @Router    /php_project/UsingPHPVersion [post]
func (s *ProjectPHPApi) UsingPHPVersion(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.IDR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.EditCategory.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//获取使用中的php版本
	data, err := ServiceGroupApp.ProjectPHPServiceApp.UsingPHPVersion(projectData.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SwitchUsingPHPVersion
// @Tags      SwitchUsingPHPVersion
// @Summary   切换使用的php版本
// @Router    /php_project/SwitchUsingPHPVersion [post]
func (s *ProjectPHPApi) SwitchUsingPHPVersion(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.SwitchUsingPHPVersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("php_project.SwitchUsingPHPVersion.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//切换使用的php版本
	err = ServiceGroupApp.ProjectPHPServiceApp.SwitchUsingPHPVersion(projectData.Name, param.Version, param.Customize)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_web.SwitchUsingPHPVersion", map[string]any{"Name": projectData.Name, "Version": param.Version}))
	response.ToResponseMsg(helper.Message("tips.SwitchSuccess"))
}

package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model/request"
	modelResponse "TTPanel/internal/model/response"
	"github.com/gin-gonic/gin"
)

type ProjectGeneralApi struct{}

// Create
// @Tags      GeneralProject
// @Summary   Create
// @Router    /general_project/Create [post]
func (s *ProjectGeneralApi) Create(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.CreateGeneralProjectR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("general_project.Create.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//创建项目
	err := ServiceGroupApp.ProjectGeneralServiceApp.Create(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_general.Create", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// List
// @Tags      GeneralProject
// @Summary   List
// @Router    /general_project/List [post]
func (s *ProjectGeneralApi) List(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.GeneralProjectListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("general_project.ProjectList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)

	//获取项目列表
	projectBaseInfoList, total, err := ServiceGroupApp.ProjectGeneralServiceApp.List(&param, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//获取进程信息
	processList, err := ServiceGroupApp.TaskManagerServiceApp.ProcessList()
	if err != nil {
		return
	}

	projectInfoList := make([]*modelResponse.GeneralProject, 0)
	//完善项目信息
	for _, projectBaseInfo := range projectBaseInfoList {
		projectInfo, _ := ServiceGroupApp.ProjectGeneralServiceApp.GetDetails(projectBaseInfo, processList)
		projectInfoList = append(projectInfoList, projectInfo)
	}

	response.ToResponseList(projectInfoList, int(total), param.Limit, param.Page)
}

// Delete
// @Tags      GeneralProject
// @Summary   Delete
// @Router    /general_project/Delete [post]
func (s *ProjectGeneralApi) Delete(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.DeleteGeneralProjectR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("general_project.Delete.BindAndValid errs: %v", errs)
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
	err = ServiceGroupApp.ProjectGeneralServiceApp.Delete(projectData, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_general.Delete", map[string]any{"Name": projectData.Name}))

	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// Info
// @Tags      GoProject
// @Summary   Info
// @Router    /general_project/Info [post]
func (s *ProjectGeneralApi) Info(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.ProjectInfoR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("general_project.Info.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//获取进程信息
	processList, err := ServiceGroupApp.TaskManagerServiceApp.ProcessList()
	if err != nil {
		return
	}

	//完善项目信息
	projectInfo, err := ServiceGroupApp.ProjectGeneralServiceApp.GetDetails(projectData, processList)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(projectInfo)
}

// SaveProjectConfig
// @Tags      GoProject
// @Summary   SaveProjectConfig
// @Router    /general_project/SaveProjectConfig [post]
func (s *ProjectGeneralApi) SaveProjectConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SaveProjectConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("general_project.SaveProjectConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//保存go项目配置
	err = ServiceGroupApp.ProjectGeneralServiceApp.SaveProjectConfig(projectData, &param)

	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_general.SaveConfig", map[string]any{"Name": projectData.Name}))
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// SetStatus
// @Tags      GoProject
// @Summary   SetStatus
// @Router    /general_project/SetStatus [post]
func (s *ProjectGeneralApi) SetStatus(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.SetStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("general_project.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取项目信息
	projectData, err := ServiceGroupApp.ProjectServiceApp.GetProjectInfoByID(param.ProjectId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	err = ServiceGroupApp.ProjectGeneralServiceApp.SetStatus(projectData, param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("project_general.SetStatus", map[string]any{"Name": projectData.Name, "Status": param.Action}))
	response.ToResponseMsg(helper.Message("tips.SetStatusSuccess"))
}

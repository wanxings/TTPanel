package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"github.com/gin-gonic/gin"
)

type BackupApi struct{}

// Database
// @Tags      Backup
// @Summary   备份数据库
// @Router    /backup/Database [post]
func (s *BackupApi) Database(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BackupDatabaseR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BackupDatabase.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//查询数据库信息
	databaseInfo, err := (&model.Databases{ID: param.Id}).Get(global.PanelDB)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	err = ServiceGroupApp.BackupServiceApp.BackupMysqlDatabase(param.StorageId, param.KeepLocalFile, databaseInfo, 0)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByBackup, helper.MessageWithMap("backup.database", map[string]any{"Name": databaseInfo.Name}))

	response.ToResponseMsg(helper.Message("tips.BackupSuccess"))
}

// Project
// @Tags      Backup
// @Summary   备份项目
// @Router    /backup/Project [post]
func (s *BackupApi) Project(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BackupProjectR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BackupProject.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//查询项目信息
	projectInfo, err := (&model.Project{ID: param.Id}).Get(global.PanelDB)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	err = ServiceGroupApp.BackupServiceApp.BackupProject(param.StorageId, param.KeepLocalFile, projectInfo, 0, param.ExclusionRules)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByBackup, helper.MessageWithMap("backup.project", map[string]any{"Name": projectInfo.Name}))
	response.ToResponseMsg(helper.Message("tips.BackupSuccess"))
}

// Dir
// @Tags      Backup
// @Summary   备份目录
// @Router    /backup/Dir [post]
func (s *BackupApi) Dir(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BackupDirR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BackupDir.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.BackupServiceApp.BackupDir(param.StorageId, param.KeepLocalFile, param.Path, 0, param.ExclusionRules, param.Description)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByBackup, helper.MessageWithMap("backup.dir", map[string]any{"Name": param.Path}))
	response.ToResponseMsg(helper.Message("tips.BackupSuccess"))
}

// Panel
// @Tags      Backup
// @Summary   备份面板
// @Router    /backup/Panel [post]
func (s *BackupApi) Panel(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BackupPanelR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BackupPanel.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.BackupServiceApp.Panel(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.BackupSuccess"))
}

// List
// @Tags      Backup
// @Summary   备份列表
// @Router    /backup/List [post]
func (s *BackupApi) List(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BackupListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BackupList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	list, total, err := ServiceGroupApp.BackupServiceApp.List(param.Category, param.Pid, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(list, int(total), limit, param.Page)
}

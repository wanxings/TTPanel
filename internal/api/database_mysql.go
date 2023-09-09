package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model/request"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DatabaseMysqlApi struct{}

// Create
// @Tags      Create
// @Summary   创建数据库
// @Router    /database_mysql/Create [post]
func (s *DatabaseMysqlApi) Create(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.CreateMysqlR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DatabaseMysqlApi.Create.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//
	err := ServiceGroupApp.DatabaseMysqlServiceApp.Create(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDatabase, helper.MessageWithMap("database.mysql.Create", map[string]any{"Name": param.DatabaseName}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// List
// @Tags      List
// @Summary   数据库列表
// @Router    /database_mysql/List [post]
func (s *DatabaseMysqlApi) List(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.ListMysqlR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DatabaseMysqlApi.List.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	//
	list, err := ServiceGroupApp.DatabaseMysqlServiceApp.List(&param, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(list)
}

// ServerList
// @Tags      ServerList
// @Summary   数据库服务列表
// @Router    /database_mysql/ServerList [post]
func (s *DatabaseMysqlApi) ServerList(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.ListMysqlServerR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DatabaseMysqlApi.ServerList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	list, err := ServiceGroupApp.DatabaseMysqlServiceApp.ServerList(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(list)
}

// GetRootPwd
// @Tags      GetRootPwd
// @Summary   获取数据库root密码
// @Router    /database_mysql/GetRootPwd [post]
func (s *DatabaseMysqlApi) GetRootPwd(c *gin.Context) {
	response := app.NewResponse(c)
	pwd, err := ServiceGroupApp.DatabaseMysqlServiceApp.GetRootPwd()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDatabase, helper.Message("database.mysql.GetRootPassword"))
	response.ToResponse(pwd)
}

// SetRootPwd
// @Tags      SetRootPwd
// @Summary   设置数据库root密码
// @Router    /database_mysql/SetRootPwd [post]
func (s *DatabaseMysqlApi) SetRootPwd(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SetRootPwdR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DatabaseMysqlApi.SetRootPwd.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.DatabaseMysqlServiceApp.SetRootPwd(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDatabase, helper.Message("database.mysql.SetRootPassword"))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SyncGetDB
// @Tags      SyncGetDB
// @Summary   从服务器同步数据库
// @Router    /database_mysql/SyncGetDB [post]
func (s *DatabaseMysqlApi) SyncGetDB(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SyncGetDBR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DatabaseMysqlApi.SyncGetDB.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	errList := ServiceGroupApp.DatabaseMysqlServiceApp.SyncGetDB(&param)
	if errList != nil {
		var errS []string
		for _, err := range errList {
			errS = append(errS, err.Error())
		}
		response.ToErrorResponse(errcode.ServerError.WithDetails(errS...))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDatabase, helper.Message("database.mysql.SyncDatabaseFromServer"))
	response.ToResponseMsg(helper.Message("tips.SyncSuccess"))
}

// SyncToDB
// @Tags      SyncToDB
// @Summary   同步数据库到服务器
// @Router    /database_mysql/SyncToDB [post]
func (s *DatabaseMysqlApi) SyncToDB(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数

	param := request.SyncToDBR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DatabaseMysqlApi.SyncToDB.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	errList := ServiceGroupApp.DatabaseMysqlServiceApp.SyncToDB(param.Ids)
	if errList != nil {
		var errS []string
		for _, err := range errList {
			errS = append(errS, err.Error())
		}
		response.ToErrorResponse(errcode.ServerError.WithDetails(errS...))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDatabase, helper.Message("database.mysql.SyncDatabaseToServer"))
	response.ToResponseMsg(helper.Message("tips.SyncSuccess"))
}

// SetAccessPermission
// @Tags      SetAccessPermission
// @Summary   设置数据库访问权限
// @Router    /database_mysql/SetAccessPermission [post]
func (s *DatabaseMysqlApi) SetAccessPermission(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SetAccessPermissionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DatabaseMysqlApi.SetAccessPermission.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.DatabaseMysqlServiceApp.SetAccessPermission(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c,
		constant.OperationLogTypeByDatabase,
		helper.MessageWithMap("database.mysql.SetDatabaseUserAccess", map[string]any{"User": param.UserName, "Permission": param.AccessPermission}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// GetAccessPermission
// @Tags      GetAccessPermission
// @Summary   获取数据库访问权限
// @Router    /database_mysql/GetAccessPermission [post]
func (s *DatabaseMysqlApi) GetAccessPermission(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.GetAccessPermissionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DatabaseMysqlApi.GetAccessPermission.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	result, err := ServiceGroupApp.DatabaseMysqlServiceApp.GetAccessPermission(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(result)
}

// SetPwd
// @Tags      SetPwd
// @Summary   设置数据库密码
// @Router    /database_mysql/SetPwd [post]
func (s *DatabaseMysqlApi) SetPwd(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SetPwdR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DatabaseMysqlApi.SetPwd.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.DatabaseMysqlServiceApp.SetPwd(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDatabase, helper.MessageWithMap("database.mysql.ChangeUserPassword", map[string]any{"User": param.UserName}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// CheckDeleteDatabase
// @Tags      CheckDeleteDatabase
// @Summary   检查是否可以删除数据库
// @Router    /database_mysql/CheckDeleteDatabase [post]
func (s *DatabaseMysqlApi) CheckDeleteDatabase(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.CheckDeleteDatabaseR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DatabaseMysqlApi.CheckDeleteDatabase.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.DatabaseMysqlServiceApp.CheckDeleteDatabase(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// DeleteDatabase
// @Tags      DeleteDatabase
// @Summary   删除数据库
// @Router    /database_mysql/DeleteDatabase [post]
func (s *DatabaseMysqlApi) DeleteDatabase(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.DeleteDatabaseR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("DatabaseMysqlApi.DeleteDatabase.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	databaseName, err := ServiceGroupApp.DatabaseMysqlServiceApp.DeleteDatabase(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDatabase, helper.MessageWithMap("database.mysql.Delete", map[string]any{"Name": databaseName}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// ImportDatabase
// @Tags      ImportDatabase
// @Summary   导入数据库
// @Router    /database_mysql/ImportDatabase [post]
func (s *DatabaseMysqlApi) ImportDatabase(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.ImportDatabaseR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		logrus.Errorf("DatabaseMysqlApi.ImportDatabase.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	databaseInfo, err := ServiceGroupApp.DatabaseMysqlServiceApp.GetDatabase(param.Id)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	logPath, err := ServiceGroupApp.DatabaseMysqlServiceApp.ImportDatabase(databaseInfo, param.FilePath, param.BackupId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDatabase, helper.MessageWithMap("database.mysql.ImportDatabase", map[string]any{"Name": databaseInfo.Name}))
	response.ToResponse(logPath)
}

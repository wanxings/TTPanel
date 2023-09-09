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

type StorageApi struct{}

// AddStorage
// @Tags     Storage
// @Summary 添加存储
// @Router    /api/storage/AddStorage [post]
func (s *StorageApi) AddStorage(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.AddStorageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Storage.AddStorage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.StorageServiceApp.AddStorage(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByStorage, helper.MessageWithMap("storage.Add", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.AddSuccess"))
}

// StorageBucketList
// @Tags     Storage
// @Summary 获取存储桶列表
// @Router    /api/storage/StorageBucketList [post]
func (s *StorageApi) StorageBucketList(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.StorageBucketListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Storage.StorageBucketList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	list, err := ServiceGroupApp.StorageServiceApp.StorageBucketList(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseList(list, len(list), 0, 0)
}

// StorageList
// @Tags     Storage
// @Summary 获取存储列表
// @Router    /api/storage/StorageList [post]
func (s *StorageApi) StorageList(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.StorageListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Storage.StorageList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	list, total, err := ServiceGroupApp.StorageServiceApp.StorageList(param.Query, param.Category, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(list, int(total), param.Limit, param.Page)
}

// EditStorage
// @Tags     Storage
// @Summary 编辑存储
// @Router    /api/storage/EditStorage [post]
func (s *StorageApi) EditStorage(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.EditStorageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Storage.EditStorage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.StorageServiceApp.EditStorage(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByStorage, helper.MessageWithMap("storage.Edit", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// LocalStorage
// @Tags     Storage
// @Summary 本地存储
// @Router    /api/storage/LocalStorage [post]
func (s *StorageApi) LocalStorage(c *gin.Context) { //Todo: 本地存储
	response := app.NewResponse(c)
	path := ServiceGroupApp.StorageServiceApp.GetLocalStoragePath()
	response.ToResponse(path)
}

// EditLocalStorage
// @Tags     Storage
// @Summary 编辑本地存储
// @Router    /api/storage/EditLocalStorage [post]
func (s *StorageApi) EditLocalStorage(c *gin.Context) { //Todo: 本地存储
	response := app.NewResponse(c)
	//获取参数
	param := request.EditLocalStorageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Storage.EditLocalStorage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.StorageServiceApp.EditLocalStoragePath(param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	request2 "TTPanel/internal/model/request"
	"github.com/gin-gonic/gin"
)

type ExtensionRedisApi struct{}

// Info
// @Tags     Redis
// @Summary   获取Redis的信息
// @Router    /extensions/redis/Info [post]
func (s *ExtensionRedisApi) Info(c *gin.Context) {
	response := app.NewResponse(c)
	//获取Redis的信息
	data, err := ServiceGroupApp.ExtensionRedisServiceApp.Info()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SetStatus
// @Tags     Redis
// @Summary   设置Redis的状态
// @Router    /extensions/redis/SetStatus [post]
func (s *ExtensionRedisApi) SetStatus(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.RedisSetStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Redis.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//设置Redis的状态
	err := ServiceGroupApp.ExtensionRedisServiceApp.SetStatus(param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExtension, helper.MessageWithMap("extension.SetStatus", map[string]any{"Name": "Redis", "Status": param.Action}))
	response.ToResponseMsg(helper.Message("tips.SetStatusSuccess"))
}

// Install
// @Tags     Redis
// @Summary   安装Redis
// @Router    /extensions/redis/Install [post]
func (s *ExtensionRedisApi) Install(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Redis.Install.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//安装Redis
	err := ServiceGroupApp.ExtensionRedisServiceApp.Install(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Install", map[string]any{"Name": "redis-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// Uninstall
// @Tags     Redis
// @Summary   卸载Redis
// @Router    /extensions/redis/Uninstall [post]
func (s *ExtensionRedisApi) Uninstall(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.VersionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Redis.Uninstall.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//卸载Redis
	err := ServiceGroupApp.ExtensionRedisServiceApp.Uninstall(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Uninstall", map[string]any{"Name": "redis-" + param.Version}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// PerformanceConfig
// @Tags     Redis
// @Summary   获取性能配置
// @Router    /extensions/redis/PerformanceConfig [post]
func (s *ExtensionRedisApi) PerformanceConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//性能配置
	data, err := ServiceGroupApp.ExtensionRedisServiceApp.PerformanceConfig()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SavePerformanceConfig
// @Tags     Redis
// @Summary   保存性能配置
// @Router    /extensions/redis/SavePerformanceConfig [post]
func (s *ExtensionRedisApi) SavePerformanceConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.RedisSavePerformanceConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Redis.SavePerformanceConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//保存性能配置
	err := ServiceGroupApp.ExtensionRedisServiceApp.SavePerformanceConfig(param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// LoadStatus
// @Tags     Redis
// @Summary   获取负载状态
// @Router    /extensions/redis/LoadStatus [post]
func (s *ExtensionRedisApi) LoadStatus(c *gin.Context) {
	response := app.NewResponse(c)
	//获取负载状态
	data, err := ServiceGroupApp.ExtensionRedisServiceApp.LoadStatus()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// PersistentConfig
// @Tags     Redis
// @Summary   获取持久化配置
// @Router    /extensions/redis/PersistentConfig [post]
func (s *ExtensionRedisApi) PersistentConfig(c *gin.Context) {
	response := app.NewResponse(c)
	//获取持久化配置
	data, err := ServiceGroupApp.ExtensionRedisServiceApp.PersistentConfig()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SavePersistentConfig
// @Tags     Redis
// @Summary   保存持久化配置
// @Router    /extensions/redis/SavePersistentConfig [post]
func (s *ExtensionRedisApi) SavePersistentConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.RedisSavePersistentConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Redis.SavePersistentConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//保存持久化配置
	err := ServiceGroupApp.ExtensionRedisServiceApp.SavePersistentConfig(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

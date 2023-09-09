package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model/request"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

type PanelApi struct{}

// OperatePanel
// @Tags     Panel
// @Summary   操作面板
// @Router    /system/panel/OperatePanel [post]
func (s *PanelApi) OperatePanel(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.OperateR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("panel.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.PanelServiceApp.OperatePanel(param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByPanel, helper.MessageWithMap("panel.OperatePanel", map[string]any{"Action": param.Action}))
	response.ToResponseMsg(helper.Message("tips.OperateSuccess"))
}

// OperateServer
// @Tags     Panel
// @Summary   操作服务器
// @Router    /system/panel/OperateServer [post]
func (s *PanelApi) OperateServer(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.OperateR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("panel.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.PanelServiceApp.OperateServer(param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByPanel, helper.MessageWithMap("panel.OperateServer", map[string]any{"Action": param.Action}))
	response.ToResponseMsg(helper.Message("tips.OperateSuccess"))
}

// ExtensionList
// @Tags     Panel
// @Summary   获取基础信息
// @Router    /system/panel/ExtensionList [post]
func (s *PanelApi) ExtensionList(c *gin.Context) {
	response := app.NewResponse(c)
	extensionList, err := ServiceGroupApp.PanelServiceApp.ExtensionList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(extensionList)
}

// Base
// @Tags     Panel
// @Summary   获取基础信息
// @Router    /system/panel/Base [post]
func (s *PanelApi) Base(c *gin.Context) {
	response := app.NewResponse(c)
	base, err := ServiceGroupApp.PanelServiceApp.Base()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(base)
}

// CheckUpdate
// @Tags     Panel
// @Summary   检查更新
// @Router    /system/panel/CheckUpdate [post]
func (s *PanelApi) CheckUpdate(c *gin.Context) {
	response := app.NewResponse(c)
	update, err := ServiceGroupApp.PanelServiceApp.CheckUpdate()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(update)
}

// Update
// @Tags     Panel
// @Summary   更新
// @Router    /system/panel/Update [post]
func (s *PanelApi) Update(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.UpdateR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("panel.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	logPath, err := ServiceGroupApp.PanelServiceApp.Update(param.Version)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByPanel, helper.Message("panel.Update"))
	response.ToResponse(logPath)
}

// Language
// @Tags     Panel
// @Summary   获取语言列表
// @Router    /system/panel/Language [post]
func (s *PanelApi) Language(c *gin.Context) {
	response := app.NewResponse(c)
	fmt.Println(language.English)
	fmt.Println(language.TraditionalChinese)
	fmt.Println(language.SimplifiedChinese)

	languageData := make(map[string]any)
	languageData["use"] = global.Config.System.Language
	languageData["list"] = global.I18n.LanguageTags()

	response.ToResponse(languageData)
}

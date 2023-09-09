package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model/request"
	"github.com/gin-gonic/gin"
	"strings"
)

type SSLApi struct{}

// CreateDnsAccount
// @Tags      CreateDnsAccount
// @Summary   创建DNS账号
// @Router    /ssl/CreateDnsAccount [post]
func (s *SSLApi) CreateDnsAccount(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.CreateDnsAccountR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ssl.CreateDnsAccount.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.SSLServiceApp.CreateDnsAccount(param.Name, param.Type, param.Authorization)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("ssl.CreateDnsAccount", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// EditDnsAccount
// @Tags      EditDnsAccount
// @Summary   编辑DNS账号
// @Router    /ssl/EditDnsAccount [post]
func (s *SSLApi) EditDnsAccount(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.CreateDnsAccountR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ssl.EditDnsAccount.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.SSLServiceApp.EditDnsAccount(param.Name, param.Authorization)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("ssl.EditDnsAccount", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// DnsTypeList
// @Tags      DnsTypeList
// @Summary   获取DNS类型列表
// @Router    /ssl/DnsTypeList [post]
func (s *SSLApi) DnsTypeList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.SSLServiceApp.DnsTypeList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// DnsAccountList
// @Tags      DnsAccountList
// @Summary   获取DNS账号列表
// @Router    /ssl/DnsAccountList [post]
func (s *SSLApi) DnsAccountList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.SSLServiceApp.DnsAccountList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// CreateAcmeAccount
// @Tags      CreateAcmeAccount
// @Summary   创建ACME账号
// @Router    /ssl/CreateAcmeAccount [post]
func (s *SSLApi) CreateAcmeAccount(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.CreateAcmeAccountR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ssl.CreateAcmeAccount.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.SSLServiceApp.CreateAcmeAccount(param.Email)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("ssl.CreateAcmeAccount", map[string]any{"Name": param.Email}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// AcmeAccountList
// @Tags      AcmeAccountList
// @Summary   获取ACME账号列表
// @Router    /ssl/AcmeAccountList [post]
func (s *SSLApi) AcmeAccountList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.SSLServiceApp.AcmeAccountList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// ApplyCertificate
// @Tags      ApplyCertificate
// @Summary   申请证书
// @Router    /ssl/ApplyCertificate [post]
func (s *SSLApi) ApplyCertificate(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.ApplyCertificateR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ssl.ApplyCertificate.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	logPath, err := ServiceGroupApp.SSLServiceApp.ApplyCertificate(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByProjectManager, helper.MessageWithMap("ssl.ApplyCertificate", map[string]any{"Domains": strings.Join(param.Domains, " "), "LogPath": logPath}))
	response.ToResponse(logPath)
}

// GetResolve
// @Tags      GetResolve
// @Summary   获取解析记录
// @Router    /ssl/GetResolve [post]
func (s *SSLApi) GetResolve(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.GetResolveR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ssl.GetResolve.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.SSLServiceApp.GetResolve(param.AcmeAccount, param.Domains)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// CheckDnsRecords
// @Tags      GetResolve
// @Summary   检查dns记录
// @Router    /ssl/CheckDnsRecords [post]
func (s *SSLApi) CheckDnsRecords(c *gin.Context) {
	response := app.NewResponse(c)
	param := request.CheckDnsRecordsR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ssl.CheckDnsRecords.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data := ServiceGroupApp.SSLServiceApp.CheckDnsRecords(param.List)
	response.ToResponse(data)
}

// CertList
// @Tags      CertList
// @Summary   获取证书列表
// @Router    /ssl/CertList [post]
func (s *SSLApi) CertList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.SSLServiceApp.CertList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

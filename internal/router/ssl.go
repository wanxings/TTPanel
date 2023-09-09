package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type SSLRouter struct{}

func (s *SSLRouter) Init(Router *gin.RouterGroup) {
	sslRouter := Router.Group("ssl")
	sslApi := api.GroupApp.SSLApiApp
	{

		sslRouter.POST("DnsTypeList", sslApi.DnsTypeList)             // 获取DNS类型列表
		sslRouter.POST("CreateDnsAccount", sslApi.CreateDnsAccount)   // 创建DNS账号
		sslRouter.POST("EditDnsAccount", sslApi.EditDnsAccount)       // 编辑DNS账号
		sslRouter.POST("DnsAccountList", sslApi.DnsAccountList)       // DNS账号列表
		sslRouter.POST("CreateAcmeAccount", sslApi.CreateAcmeAccount) // 创建ACME账号
		sslRouter.POST("AcmeAccountList", sslApi.AcmeAccountList)     // ACME账号列表
		sslRouter.POST("ApplyCertificate", sslApi.ApplyCertificate)   // 申请证书
		sslRouter.POST("GetResolve", sslApi.GetResolve)               // 获取解析记录
		sslRouter.POST("CheckDnsRecords", sslApi.CheckDnsRecords)     //检查dns记录
		sslRouter.POST("CertList", sslApi.CertList)                   // 证书列表
	}
}

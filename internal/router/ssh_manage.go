package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type SSHManageRouter struct{}

func (s *SSHManageRouter) Init(Router *gin.RouterGroup) {
	sshManageRouter := Router.Group("ssh_manage")
	sshManageApi := api.GroupApp.SSHManageApiApp
	{
		sshManageRouter.POST("GetSSHInfo", sshManageApi.GetSSHInfo)                       //获取SSH信息
		sshManageRouter.POST("SetSSHStatus", sshManageApi.SetSSHStatus)                   //设置SSH状态
		sshManageRouter.POST("OperateSSHKeyLogin", sshManageApi.OperateSSHKeyLogin)       //操作SSH密钥登录
		sshManageRouter.POST("OperatePasswordLogin", sshManageApi.OperatePasswordLogin)   //操作密码登录
		sshManageRouter.POST("GetSSHLoginStatistics", sshManageApi.GetSSHLoginStatistics) //获取SSH登录统计
	}
}

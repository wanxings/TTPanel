package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (u *UserRouter) Init(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	userApi := api.GroupApp.UserApiApp
	{
		userRouter.POST("info", userApi.Info)     // 用户信息
		userRouter.POST("logout", userApi.Logout) // 退出登录
	}
	{
		userRouter.POST("CreateTemporaryUser", userApi.CreateTemporaryUser) // 创建临时用户
		//userRouter.POST("TemporaryUserList", userApi.TemporaryUserList)               // 临时用户列表
		//userRouter.POST("ForceLogoutTemporaryUser", userApi.ForceLogoutTemporaryUser) // 强制退出临时用户
	}
}

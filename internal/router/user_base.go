package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type UserBaseRouter struct{}

func (u *UserBaseRouter) InitBase(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	userApi := api.GroupApp.UserApiApp
	{
		baseRouter.POST("login", userApi.Login) // 登陆
		//baseRouter.POST("captcha", baseApi.GetCaptcha) // 获取验证码
	}
	panelApi := api.GroupApp.PanelApiApp
	{
		baseRouter.POST("language", panelApi.Language) // 获取语言列表
	}

}

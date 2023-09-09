package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type WebSSHRouter struct{}

func (s *WebSSHRouter) Init(Router *gin.RouterGroup) {
	webSSHApi := api.GroupApp.WebSSHApiApp
	{
		Router.GET("webSSH", webSSHApi.WsSSH) // webSSH
	}
}

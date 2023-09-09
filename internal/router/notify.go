package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type NotifyRouter struct{}

func (s *NotifyRouter) Init(Router *gin.RouterGroup) {
	notifyRouter := Router.Group("notify")
	notifyApi := api.GroupApp.NotifyApiApp
	{
		notifyRouter.POST("AddNotifyChannel", notifyApi.AddNotifyChannel)       // 添加通知通道
		notifyRouter.POST("NotifyChannelList", notifyApi.NotifyChannelList)     // 通知通道列表
		notifyRouter.POST("TestNotifyChannel", notifyApi.TestNotifyChannel)     // 测试通知通道
		notifyRouter.POST("EditNotifyChannel", notifyApi.EditNotifyChannel)     // 编辑通知通道
		notifyRouter.POST("DeleteNotifyChannel", notifyApi.DeleteNotifyChannel) // 删除通知通道
	}

}

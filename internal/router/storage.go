package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type StorageRouter struct{}

func (s *StorageRouter) Init(Router *gin.RouterGroup) {
	storageRouter := Router.Group("storage")
	storageApi := api.GroupApp.StorageApiApp
	{
		storageRouter.POST("AddStorage", storageApi.AddStorage)               // 添加存储
		storageRouter.POST("StorageBucketList", storageApi.StorageBucketList) // 获取存储桶列表
		storageRouter.POST("StorageList", storageApi.StorageList)             // 获取存储列表
		storageRouter.POST("EditStorage", storageApi.EditStorage)             // 编辑存储
		storageRouter.POST("LocalStorage", storageApi.LocalStorage)           // 本地存储
		storageRouter.POST("EditLocalStorage", storageApi.EditLocalStorage)   // 编辑本地存储
	}
}

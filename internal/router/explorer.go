package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ExplorerRouter struct{}

func (s *ExplorerRouter) Init(Router *gin.RouterGroup) {
	filesRouter := Router.Group("explorer")
	explorerApi := api.GroupApp.ExplorerApiApp
	{
		filesRouter.POST("GetDir", explorerApi.GetDir)                                               // 获取目录文件
		filesRouter.POST("GetFileBody", explorerApi.GetFileBody)                                     // 获取文件内容
		filesRouter.POST("GetAttribute", explorerApi.GetAttribute)                                   // 获取文件（夹）属性
		filesRouter.POST("SaveFileBody", explorerApi.SaveFileBody)                                   // 写入文件内容
		filesRouter.POST("BatchSetFileAccess", explorerApi.BatchSetFileAccess)                       // 批量设置文件（夹）权限
		filesRouter.POST("BatchDeleteDirFile", explorerApi.BatchDeleteDirFile)                       // 批量删除文件（夹）
		filesRouter.POST("Rename", explorerApi.Rename)                                               // 重命名文件（夹）
		filesRouter.POST("BatchCheckExistsFiles", explorerApi.BatchCheckExistsFiles)                 //批量检查文件是否存在
		filesRouter.POST("BatchCopy", explorerApi.BatchCopy)                                         //批量复制文件（夹）
		filesRouter.POST("BatchMove", explorerApi.BatchMove)                                         //批量移动文件（夹）
		filesRouter.POST("CreateDir", explorerApi.CreateDir)                                         //创建文件夹
		filesRouter.POST("CreateFile", explorerApi.CreateFile)                                       //创建文件
		filesRouter.POST("CreateSymlink", explorerApi.CreateSymlink)                                 //创建符号链接
		filesRouter.POST("CreateDuplicate", explorerApi.CreateDuplicate)                             //创建副本
		filesRouter.POST("GetPathSize", explorerApi.GetPathSize)                                     //获取文件夹大小
		filesRouter.POST("CheckFileExists", explorerApi.CheckFileExists)                             //检查文件是否存在
		filesRouter.POST("Upload", explorerApi.Upload)                                               //上传文件
		filesRouter.POST("Compress", explorerApi.Compress)                                           //压缩文件(夹)
		filesRouter.POST("Decompress", explorerApi.Decompress)                                       //解压文件(夹)
		filesRouter.POST("GetLogContent", explorerApi.GetLogContent)                                 //获取日志
		filesRouter.POST("ClearLogContent", explorerApi.ClearLogContent)                             // 清空日志
		filesRouter.POST("SearchFileContent", explorerApi.SearchFileContent)                         // 搜索文件内容
		filesRouter.POST("RemoteDownload", explorerApi.RemoteDownload)                               // 远程下载文件
		filesRouter.POST("RemoteDownloadProcess", explorerApi.RemoteDownloadProcess)                 // 远程下载文件进度
		filesRouter.GET("Download", explorerApi.Download)                                            // 下载文件
		filesRouter.POST("BatchOperateSpecialPermission", explorerApi.BatchOperateSpecialPermission) // 批量操作特殊权限
		filesRouter.POST("SetRemark", explorerApi.SetRemark)                                         // 设置备注
		filesRouter.POST("FavoritesList", explorerApi.FavoritesList)                                 // 收藏列表
		filesRouter.POST("OperateFavorites", explorerApi.OperateFavorites)                           // 操作收藏
		filesRouter.POST("GenerateDownloadExternalLink", explorerApi.GenerateDownloadExternalLink)   // 生成下载外链
		filesRouter.POST("DownloadExternalLinkList", explorerApi.DownloadExternalLinkList)           // 下载外链列表
		filesRouter.POST("DeleteDownloadExternalLink", explorerApi.DeleteDownloadExternalLink)       // 删除下载外链
		filesRouter.POST("GetFileTemporaryDownloadLink", explorerApi.GetFileTemporaryDownloadLink)   // 获取文件临时下载链接
	}
}
func (s *ExplorerRouter) InitExternal(Router *gin.RouterGroup) {
	explorerApi := api.GroupApp.ExplorerApiApp
	{
		Router.GET("ExternalDownload", explorerApi.ExternalDownload) // 下载文件
	}
}

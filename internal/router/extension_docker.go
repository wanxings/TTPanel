package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type ExtensionDockerRouter struct{}

func (s *ExtensionDockerRouter) Init(Router *gin.RouterGroup) {
	dockerR := Router.Group("extension_docker")
	dockerApi := api.GroupApp.ExtensionDockerApiApp
	{
		dockerR.POST("Install", dockerApi.Install)     // 安装docker
		dockerR.POST("Uninstall", dockerApi.Uninstall) // 卸载docker
		dockerR.POST("Info", dockerApi.Info)           // 获取docker信息
		dockerR.POST("SetStatus", dockerApi.SetStatus) // 设置docker状态
		dockerR.POST("RunInfo", dockerApi.RunInfo)     // 获取docker运行信息
	}
	//APP
	{
		dockerR.POST("AppList", dockerApi.AppList)             //app列表
		dockerR.POST("UpdateAppList", dockerApi.UpdateAppList) // 更新app列表
		dockerR.POST("AppConfig", dockerApi.AppConfig)         //app配置
		dockerR.POST("DeployApp", dockerApi.DeployApp)         //部署app
	}
	//Image
	{
		dockerR.POST("PullImage", dockerApi.PullImage)               // 拉取镜像
		dockerR.POST("ImageList", dockerApi.ImageList)               // 获取镜像列表
		dockerR.POST("ImportImage", dockerApi.ImportImage)           // 导入镜像
		dockerR.POST("BuildImage", dockerApi.BuildImage)             // 构建镜像
		dockerR.POST("PushImage", dockerApi.PushImage)               // 推送镜像
		dockerR.POST("ExportImage", dockerApi.ExportImage)           // 导出镜像
		dockerR.POST("BatchDeleteImage", dockerApi.BatchDeleteImage) // 批量删除镜像
	}
	//repository
	{
		dockerR.POST("AddRepository", dockerApi.AddRepository)       // 添加仓库
		dockerR.POST("RepositoryList", dockerApi.RepositoryList)     // 获取仓库列表
		dockerR.POST("EditRepository", dockerApi.EditRepository)     // 编辑仓库
		dockerR.POST("DeleteRepository", dockerApi.DeleteRepository) // 删除仓库
	}
	//composeTemplate
	{
		dockerR.POST("CreateComposeTemplate", dockerApi.CreateComposeTemplate)           // 创建模板
		dockerR.POST("ComposeTemplateList", dockerApi.ComposeTemplateList)               // 获取模板列表
		dockerR.POST("BatchDeleteComposeTemplate", dockerApi.BatchDeleteComposeTemplate) // 批量删除模板
		dockerR.POST("EditComposeTemplate", dockerApi.EditComposeTemplate)               // 编辑模板
		dockerR.POST("ComposeTemplatePullImage", dockerApi.ComposeTemplatePullImage)     // 拉取镜像
	}
	//Compose
	{
		dockerR.POST("CreateCompose", dockerApi.CreateCompose)           // 创建Compose项目
		dockerR.POST("ComposeList", dockerApi.ComposeList)               // 获取Compose项目列表
		dockerR.POST("ComposeConfig", dockerApi.ComposeConfig)           // 获取Compose应用配置
		dockerR.POST("SaveComposeConfig", dockerApi.SaveComposeConfig)   // 保存Compose应用配置
		dockerR.POST("GetComposeServices", dockerApi.GetComposeServices) // 获取Compose服务列表
		dockerR.POST("OperateCompose", dockerApi.OperateCompose)         // 操作Compose项目
		dockerR.POST("DeleteCompose", dockerApi.DeleteCompose)           // 删除Compose项目
	}
	//Container
	{
		dockerR.POST("CreateContainer", dockerApi.CreateContainer)   // 创建容器
		dockerR.POST("ContainerList", dockerApi.ContainerList)       // 获取容器列表
		dockerR.POST("OperateContainer", dockerApi.OperateContainer) // 操作容器
		dockerR.GET("ContainerSSH", dockerApi.ContainerSSH)          // 容器SSH
		dockerR.POST("ContainerLogs", dockerApi.ContainerLogs)       // 容器日志
		dockerR.POST("ContainerMonitor", dockerApi.ContainerMonitor) // 容器监控
	}
	//Networking
	{
		dockerR.POST("CreateNetworking", dockerApi.CreateNetworking) // 创建网络
		dockerR.POST("NetworkingList", dockerApi.NetworkingList)     // 获取网络列表
		dockerR.POST("DeleteNetworking", dockerApi.DeleteNetworking) // 删除网络
	}
	//Volume
	{
		dockerR.POST("CreateVolume", dockerApi.CreateVolume) // 创建存储
		dockerR.POST("VolumeList", dockerApi.VolumeList)     // 获取存储列表
		dockerR.POST("DeleteVolume", dockerApi.DeleteVolume) // 删除存储
	}
}

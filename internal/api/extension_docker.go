package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	request2 "TTPanel/internal/model/request"
	"TTPanel/pkg/util"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

type ExtensionDockerApi struct{}

// Install
// @Tags     Docker
// @Summary   安装docker
// @Router    /extension_docker/Install [post]
func (s *ExtensionDockerApi) Install(c *gin.Context) {
	response := app.NewResponse(c)
	//安裝Docker
	err := ServiceGroupApp.ExtensionDockerServiceApp.Install()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Install", map[string]any{"Name": "docker"}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// Uninstall
// @Tags     Docker
// @Summary   卸载docker
// @Router    /extension_docker/Uninstall [post]
func (s *ExtensionDockerApi) Uninstall(c *gin.Context) {
	response := app.NewResponse(c)

	//卸载docker
	err := ServiceGroupApp.ExtensionDockerServiceApp.Uninstall()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByQueueTask, helper.MessageWithMap("queuetask.Uninstall", map[string]any{"Name": "docker"}))
	response.ToResponseMsg(helper.Message("tips.AddedToQueueTask"))
}

// Info
// @Tags     Docker
// @Summary   获取docker信息
// @Router    /extension_docker/Info [post]
func (s *ExtensionDockerApi) Info(c *gin.Context) {
	response := app.NewResponse(c)

	//获取docker信息
	data, err := ServiceGroupApp.ExtensionDockerServiceApp.Info()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// SetStatus
// @Tags     Docker
// @Summary   设置docker状态
// @Router    /extension_docker/SetStatus [post]
func (s *ExtensionDockerApi) SetStatus(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DockerSetStatusR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.SetStatus.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//设置docker状态
	_, err := ServiceGroupApp.ExtensionDockerServiceApp.NewDockerService()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = ServiceGroupApp.ExtensionDockerServiceApp.SetStatus(param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExtension, helper.MessageWithMap("extension.SetStatus", map[string]any{"Name": "docker", "Status": param.Action}))
	response.ToResponseMsg(helper.Message("tips.SetStatusSuccess"))
}

// RunInfo
// @Tags     Docker
// @Summary   获取docker运行信息
// @Router    /extension_docker/RunInfo [post]
func (s *ExtensionDockerApi) RunInfo(c *gin.Context) {
	response := app.NewResponse(c)
	//获取docker运行信息
	data, err := ServiceGroupApp.ExtensionDockerServiceApp.RunInfo()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// AppList
// @Tags     Docker
// @Summary  APP列表
// @Router    /extensions_docker/AppList [post]
func (s *ExtensionDockerApi) AppList(c *gin.Context) {
	response := app.NewResponse(c)
	//APP列表
	data, err := ServiceGroupApp.ExtensionDockerServiceApp.AppList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// UpdateAppList
// @Tags     Docker
// @Summary   更新APP列表
// @Router    /extension_docker/UpdateAppList [post]
func (s *ExtensionDockerApi) UpdateAppList(c *gin.Context) {
	response := app.NewResponse(c)
	//更新APP列表
	err := ServiceGroupApp.ExtensionDockerServiceApp.UpdateAppList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.Message("docker.app.UpdateAppList"))
	response.ToResponseMsg(helper.Message("tips.updateSuccess"))
}

// AppConfig
// @Tags     Docker
// @Summary   APP配置
// @Router    /extension_docker/ComposeConfig [post]
func (s *ExtensionDockerApi) AppConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DockerAppConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.ComposeConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//APP详情
	data, err := ServiceGroupApp.ExtensionDockerServiceApp.ComposeConfig(global.Config.System.PanelPath + "/data/docker_app/" + param.ServerName)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// DeployApp
// @Tags     Docker
// @Summary   部署APP
// @Router    /extension_docker/DeployApp [post]
func (s *ExtensionDockerApi) DeployApp(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DockerDeployAppR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.DeployApp.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//部署APP
	_, err := ServiceGroupApp.ExtensionDockerServiceApp.NewDockerService()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	logPath, err := ServiceGroupApp.ExtensionDockerServiceApp.DeployApp(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.app.DeployApp", map[string]any{"Name": param.ServerName}))
	response.ToResponse(logPath)
}

// PullImage
// @Tags     Docker
// @Summary   拉取镜像
// @Router    /extension_docker/PullImage [post]
func (s *ExtensionDockerApi) PullImage(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DockerPullImageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.PullImage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//拉取镜像
	_, err := ServiceGroupApp.ExtensionDockerServiceApp.NewDockerService()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	logPath, err := ServiceGroupApp.ExtensionDockerServiceApp.PullImage(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.image.Pull", map[string]any{"Name": param.Name, "LogPath": logPath}))
	response.ToResponse(logPath)
}

// ImageList
// @Tags     Docker
// @Summary   镜像列表
// @Router    /extension_docker/ImageList [post]
func (s *ExtensionDockerApi) ImageList(c *gin.Context) {
	response := app.NewResponse(c)

	data, total, err := ServiceGroupApp.ExtensionDockerServiceApp.ImageList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponseList(data, total, 0, 0)
}

// ImportImage
// @Tags     Docker
// @Summary   导入镜像
// @Router    /extension_docker/ImportImage [post]
func (s *ExtensionDockerApi) ImportImage(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DockerImportImageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.ImportImage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//导入镜像
	err := ServiceGroupApp.ExtensionDockerServiceApp.ImportImage(param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.image.Import", map[string]any{"Path": param.Path}))
	response.ToResponseMsg(helper.Message("tips.ImportSuccess"))
}

// BuildImage
// @Tags     Docker
// @Summary   构建镜像
// @Router    /extension_docker/BuildImage [post]
func (s *ExtensionDockerApi) BuildImage(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DockerBuildImageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.BuildImage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//构建镜像
	logPath, err := ServiceGroupApp.ExtensionDockerServiceApp.BuildImage(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.image.Build", map[string]any{"Name": param.Tags, "LogPath": logPath}))
	response.ToResponse(logPath)
}

// PushImage
// @Tags     Docker
// @Summary   推送镜像
// @Router    /extension_docker/PushImage [post]
func (s *ExtensionDockerApi) PushImage(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DockerPushImageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.PushImage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//推送镜像
	logPath, err := ServiceGroupApp.ExtensionDockerServiceApp.PushImage(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.image.Push", map[string]any{"ID": param.ImageId, "LogPath": logPath}))
	response.ToResponse(logPath)
}

// ExportImage
// @Tags     Docker
// @Summary   导出镜像
// @Router    /extension_docker/ExportImage [post]
func (s *ExtensionDockerApi) ExportImage(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DockerExportImageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.ExportImage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//导出镜像
	savePath, err := ServiceGroupApp.ExtensionDockerServiceApp.ExportImage(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.image.Export", map[string]any{"ID": param.ImageId, "SavePath": savePath}))
	response.ToResponseMsg(helper.Message("tips.ExportSuccess"))
}

// BatchDeleteImage
// @Tags     Docker
// @Summary   批量删除镜像
// @Router    /extension_docker/BatchDeleteImage [post]
func (s *ExtensionDockerApi) BatchDeleteImage(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DockerBatchDeleteImageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.BatchDeleteImage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//批量删除镜像
	for _, imageId := range param.ImageIds {
		err := ServiceGroupApp.ExtensionDockerServiceApp.DeleteImage(imageId)
		if err != nil {
			response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
			return
		}

		go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
			c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.image.Delete", map[string]any{"ID": imageId}))
	}

	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// AddRepository
// @Tags     Docker
// @Summary   添加仓库
// @Router    /extension_docker/AddRepository [post]
func (s *ExtensionDockerApi) AddRepository(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.AddRepositoryR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.AddRepository.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//添加仓库
	err := ServiceGroupApp.ExtensionDockerServiceApp.AddRepository(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.repository.Add", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.AddSuccess"))
}

// RepositoryList
// @Tags     Docker
// @Summary   仓库列表
// @Router    /extension_docker/RepositoryList [post]
func (s *ExtensionDockerApi) RepositoryList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.RepositoryListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.RepositoryList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//仓库列表
	data, total, err := ServiceGroupApp.ExtensionDockerServiceApp.RepositoryList(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// EditRepository
// @Tags     Docker
// @Summary   编辑仓库
// @Router    /extension_docker/EditRepository [post]
func (s *ExtensionDockerApi) EditRepository(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.EditRepositoryR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.EditRepository.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//编辑仓库
	err := ServiceGroupApp.ExtensionDockerServiceApp.EditRepository(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.repository.Edit", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// DeleteRepository
// @Tags     Docker
// @Summary   删除仓库
// @Router    /extension_docker/DeleteRepository [post]
func (s *ExtensionDockerApi) DeleteRepository(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DeleteRepositoryR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.DeleteRepository.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//删除仓库
	data, err := ServiceGroupApp.ExtensionDockerServiceApp.DeleteRepository(param.Id)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.repository.Delete", map[string]any{"Name": data.Name}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// CreateComposeTemplate
// @Tags     Docker
// @Summary   创建compose模板
// @Router    /extension_docker/CreateComposeTemplate [post]
func (s *ExtensionDockerApi) CreateComposeTemplate(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.CreateComposeTemplateR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.CreateComposeTemplate.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//创建compose模板
	err := ServiceGroupApp.ExtensionDockerServiceApp.CreateComposeTemplate(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.compose_template.Create", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// ComposeTemplateList
// @Tags     Docker
// @Summary   compose模板列表
// @Router    /extension_docker/ComposeTemplateList [post]
func (s *ExtensionDockerApi) ComposeTemplateList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.ComposeTemplateListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.ComposeTemplateList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)

	//compose模板列表
	data, total, err := ServiceGroupApp.ExtensionDockerServiceApp.ComposeTemplateList(param.Query, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// BatchDeleteComposeTemplate
// @Tags     Docker
// @Summary   批量删除compose模板
// @Router    /extension_docker/BatchDeleteComposeTemplate [post]
func (s *ExtensionDockerApi) BatchDeleteComposeTemplate(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.BatchDeleteComposeTemplateR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.BatchDeleteComposeTemplate.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//批量删除compose模板
	for _, id := range param.Ids {
		template, err := ServiceGroupApp.ExtensionDockerServiceApp.DeleteComposeTemplate(id)
		if err != nil {
			response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
			return
		}

		go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
			c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.compose_template.Delete", map[string]any{"Name": template.Name}))
	}

	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// EditComposeTemplate
// @Tags     Docker
// @Summary   编辑compose模板
// @Router    /extension_docker/EditComposeTemplate [post]
func (s *ExtensionDockerApi) EditComposeTemplate(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.EditComposeTemplateR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.EditComposeTemplate.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//模板是否存在
	template, err := ServiceGroupApp.ExtensionDockerServiceApp.GetComposeTemplateByID(param.Id)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	//编辑compose模板
	err = ServiceGroupApp.ExtensionDockerServiceApp.EditComposeTemplate(template, &param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.compose_template.Edit", map[string]any{"Name": template.Name, "ID": template.ID}))
	response.ToResponseMsg(helper.Message("tips.EditSuccess"))
}

// ComposeTemplatePullImage
// @Tags     Docker
// @Summary   拉取compose模板镜像
// @Router    /extension_docker/ComposeTemplatePullImage [post]
func (s *ExtensionDockerApi) ComposeTemplatePullImage(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.ComposeTemplatePullImageR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.ComposeTemplatePullImage.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//拉取compose模板镜像
	logPath, err := ServiceGroupApp.ExtensionDockerServiceApp.ComposeTemplatePullImage(param.Id)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.compose_template.PullImage", map[string]any{"ID": param.Id, "LogPath": logPath}))
	response.ToResponse(logPath)
}

// CreateCompose
// @Tags     Docker
// @Summary   创建compose
// @Router    /extension_docker/CreateCompose [post]
func (s *ExtensionDockerApi) CreateCompose(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.CreateComposeR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.CreateCompose.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//创建compose
	logPath, err := ServiceGroupApp.ExtensionDockerServiceApp.CreateCompose(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.compose.Create", map[string]any{"Name": param.Name, "LogPath": logPath}))
	response.ToResponse(logPath)
}

// ComposeList
// @Tags     Docker
// @Summary   compose列表
// @Router    /extension_docker/ComposeList [post]
func (s *ExtensionDockerApi) ComposeList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.ComposeListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.ComposeList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)

	//compose列表
	data, total, err := ServiceGroupApp.ExtensionDockerServiceApp.ComposeList(param.Query, offset, limit)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseList(data, int(total), param.Limit, param.Page)
}

// ComposeConfig
// @Tags     Docker
// @Summary   compose配置
// @Router    /extension_docker/ComposeConfig [post]
func (s *ExtensionDockerApi) ComposeConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.ComposeConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.ComposeConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.ExtensionDockerServiceApp.ComposeConfig(param.ComposePath)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	envFileBody, _ := util.ReadFileStringBody(param.ComposePath + "/.env")

	envVariables := make(map[string]any)

	lines := strings.Split(envFileBody, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]

		envVariables[key] = value
	}
	for k, v := range data.Params {
		if value, ok := envVariables[v.Key]; ok {
			v.Value = value
			data.Params[k] = v
		}
	}

	response.ToResponse(data)
}

// SaveComposeConfig
// @Tags     Docker
// @Summary   保存compose配置
// @Router    /extension_docker/SaveComposeConfig [post]
func (s *ExtensionDockerApi) SaveComposeConfig(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.SaveComposeConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.SaveComposeConfig.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExtensionDockerServiceApp.SaveComposeConfig(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.compose.SaveConfig", map[string]any{"ComposePath": param.ComposePath}))

	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))
}

// GetComposeServices
// @Tags     Docker
// @Summary   获取compose服务
// @Router    /extension_docker/GetComposeServices [post]
func (s *ExtensionDockerApi) GetComposeServices(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.ComposeConfigR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.GetComposeServices.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.ExtensionDockerServiceApp.GetComposeServices(param.ComposePath)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// OperateCompose
// @Tags     Docker
// @Summary   操作compose
// @Router    /extension_docker/OperateCompose [post]
func (s *ExtensionDockerApi) OperateCompose(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.OperateComposeR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.OperateCompose.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//编排项目是否存在
	composeProject, err := ServiceGroupApp.ExtensionDockerServiceApp.GetComposeByID(param.ID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	//操作compose
	logData, err := ServiceGroupApp.ExtensionDockerServiceApp.OperateCompose(composeProject.Path+"/docker-compose.yml", composeProject.Name, param.Action, param.Services)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.compose.Operate", map[string]any{"Name": composeProject.Name, "Action": param.Action, "Service": strings.Join(param.Services, " ")}))

	response.ToResponse(logData)
}

// DeleteCompose
// @Tags     Docker
// @Summary   删除compose
// @Router    /extension_docker/DeleteCompose [post]
func (s *ExtensionDockerApi) DeleteCompose(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DeleteComposeR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.DeleteCompose.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//compose是否存在
	compose, err := ServiceGroupApp.ExtensionDockerServiceApp.GetComposeByID(param.Id)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	//删除compose
	err = ServiceGroupApp.ExtensionDockerServiceApp.DeleteCompose(compose)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	//判断是否删除目录
	if param.DelPath {
		_ = os.RemoveAll(compose.Path)
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.compose.Delete", map[string]any{"Name": compose.Name}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// CreateContainer
// @Tags     Docker
// @Summary   创建容器
// @Router    /extension_docker/CreateContainer [post]
func (s *ExtensionDockerApi) CreateContainer(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.CreateContainerR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.CreateContainer.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//创建容器
	err := ServiceGroupApp.ExtensionDockerServiceApp.CreateContainer(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.container.Create", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// ContainerList
// @Tags     Docker
// @Summary   容器列表
// @Router    /extension_docker/ContainerList [post]
func (s *ExtensionDockerApi) ContainerList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.ContainerListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.ContainerList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//容器列表
	data, err := ServiceGroupApp.ExtensionDockerServiceApp.ContainerList(param.Filters, param.GetStats)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// OperateContainer
// @Tags     Docker
// @Summary   操作容器
// @Router    /extension_docker/OperateContainer [post]
func (s *ExtensionDockerApi) OperateContainer(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.OperateContainerR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.OperateContainer.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//操作容器
	err := ServiceGroupApp.ExtensionDockerServiceApp.OperateContainer(param.Name, param.NewName, param.Action)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.container.Operate", map[string]any{"Name": param.Name, "Action": param.Action}))
	response.ToResponseMsg(helper.Message("tips.OperateSuccess"))
}

// ContainerLogs
// @Tags     Docker
// @Summary   容器日志
// @Router    /extension_docker/ContainerLogs [post]
func (s *ExtensionDockerApi) ContainerLogs(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.ContainerLogR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.ContainerLog.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//容器日志
	data, err := ServiceGroupApp.ExtensionDockerServiceApp.ContainerLogs(param.Name, param.Since, param.Until, param.Tail)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// ContainerMonitor
// @Tags     Docker
// @Summary   容器监控
// @Router    /extension_docker/ContainerMonitor [post]
func (s *ExtensionDockerApi) ContainerMonitor(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.ContainerMonitorR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.ContainerMonitor.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//容器监控
	data, err := ServiceGroupApp.ExtensionDockerServiceApp.ContainerMonitor(param.Id)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// ContainerSSH
// @Tags     Docker
// @Summary   容器SSH
// @Router    /extension_docker/ContainerSSH [post]
func (s *ExtensionDockerApi) ContainerSSH(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.ContainerSSHR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Host.Terminal.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.container.ConnectSSH", map[string]any{"ContainerId": param.ContainerId, "User": param.User, "Command": param.Command}))
	global.Log.Debugf("容器SSH:%v\n", param)
	ServiceGroupApp.ExtensionDockerServiceApp.ContainerSSH(c, &param)
}

// CreateNetworking
// @Tags     Docker
// @Summary   创建网络
// @Router    /extension_docker/CreateNetworking [post]
func (s *ExtensionDockerApi) CreateNetworking(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.CreateNetworkingR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.CreateNetworking.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//创建网络
	err := ServiceGroupApp.ExtensionDockerServiceApp.CreateNetworking(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.networking.Create", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// NetworkingList
// @Tags     Docker
// @Summary   网络列表
// @Router    /extension_docker/NetworkingList [post]
func (s *ExtensionDockerApi) NetworkingList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.NetworkingListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.NetworkingList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//网络列表
	data, total, err := ServiceGroupApp.ExtensionDockerServiceApp.NetworkingList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseList(data, total, 0, 0)
}

// DeleteNetworking
// @Tags     Docker
// @Summary   删除网络
// @Router    /extension_docker/DeleteNetworking [post]
func (s *ExtensionDockerApi) DeleteNetworking(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DeleteNetworkingR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.DeleteNetworking.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//删除网络
	err := ServiceGroupApp.ExtensionDockerServiceApp.DeleteNetworking(param.Id)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.networking.Delete", map[string]any{"ID": param.Id}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// CreateVolume
// @Tags     Docker
// @Summary   创建存储
// @Router    /extension_docker/CreateVolume [post]
func (s *ExtensionDockerApi) CreateVolume(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.CreateVolumeR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.CreateVolume.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//创建存储
	err := ServiceGroupApp.ExtensionDockerServiceApp.CreateVolume(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.volume.Create", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// VolumeList
// @Tags     Docker
// @Summary   存储列表
// @Router    /extension_docker/VolumeList [post]
func (s *ExtensionDockerApi) VolumeList(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.VolumeListR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.VolumeList.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//存储列表
	data, total, err := ServiceGroupApp.ExtensionDockerServiceApp.VolumeList(param.Name)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseList(data, total, 0, 0)
}

// DeleteVolume
// @Tags     Docker
// @Summary   删除存储
// @Router    /extension_docker/DeleteVolume [post]
func (s *ExtensionDockerApi) DeleteVolume(c *gin.Context) {
	response := app.NewResponse(c)
	param := request2.DeleteVolumeR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("docker.DeleteVolume.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//删除存储
	err := ServiceGroupApp.ExtensionDockerServiceApp.DeleteVolume(param.Name, param.Force)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByDockerManager, helper.MessageWithMap("docker.volume.Delete", map[string]any{"Name": param.Name}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

package service

import (
	"TTPanel/internal/core/terminal"
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ExtensionDockerService struct {
}

// NewDockerService Docker服务
func (s *ExtensionDockerService) NewDockerService() (*docker.Client, error) {
	var err error
	err = s.InitCheck()
	if err != nil {
		return nil, err
	}
	return docker.NewClientWithOpts(docker.WithHost("unix:///var/run/docker.sock"), docker.WithAPIVersionNegotiation())
}

// Info 获取docker信息
func (s *ExtensionDockerService) Info() (*response.ExtensionsInfoResponse, error) {
	var err error
	var dockerInfo response.ExtensionsInfoResponse
	err = ReadExtensionsInfo(constant.ExtensionDockerName, &dockerInfo)
	if err != nil {
		global.Log.Errorf("ReadDockerExtensionsInfo  Error：%v", err.Error())
		return nil, err
	}
	dockerInfo.Description.Install = s.CheckDockerInstalled()
	dockerInfo.Description.Status = s.IsRunning()
	if dockerInfo.Description.Install && !util.PathExists("/etc/docker/daemon.json") {
		// /etc/docker/daemon.json
		_ = util.WriteFile("/etc/docker/daemon.json", []byte(`{}`), 0755)
	}

	return &dockerInfo, nil
}

// BaseStatistics 获取基础统计信息
func (s *ExtensionDockerService) BaseStatistics() (*response.DockerBaseStatistics, error) {
	stats := &response.DockerBaseStatistics{}
	stats.Install = s.CheckDockerInstalled()
	stats.Status = s.IsRunning()
	// 创建Docker客户端
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return stats, err
	}

	// 获取编排列表
	composeTotal, err := (&model.DockerCompose{}).Count(global.PanelDB, &model.ConditionsT{})
	if err != nil {
		global.Log.Errorf("BaseStatistics->DockerCompose.Count Error:%v", err.Error())
		return stats, err
	}

	// 获取镜像列表
	imageList, err := dockerClient.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		global.Log.Errorf("BaseStatistics->ImageList Error:%v", err.Error())
		return stats, err
	}

	// 获取容器列表
	containerList, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		global.Log.Errorf("BaseStatistics->ContainerList Error:%v", err.Error())
		return stats, err
	}

	// 获取存储卷列表
	volumeList, err := dockerClient.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		global.Log.Errorf("BaseStatistics->VolumeList Error:%v", err.Error())
		return stats, err
	}

	// 获取网络列表
	networkList, err := dockerClient.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		global.Log.Errorf("BaseStatistics->NetworkList Error:%v", err.Error())
		return stats, err
	}

	//填充Container
	for _, containerInfo := range containerList {
		stats.Container.Total++
		switch containerInfo.State {
		case "running":
			stats.Container.Running++
		case "exited":
			stats.Container.Exited++
		case "paused":
			stats.Container.Paused++
		case "removing":
			stats.Container.Removing++
		case "restarting":
			stats.Container.Restarting++
		case "created":
			stats.Container.Created++
		}
	}

	//填充Compose
	stats.Compose.Total = composeTotal

	// 填充Image
	stats.Image.Total = len(imageList)
	for _, img := range imageList {
		stats.Image.Size += img.Size
	}

	//填充Volume
	stats.Volume.Total = len(volumeList.Volumes)

	//填充Network
	stats.Network.Total = len(networkList)

	return stats, nil
}

// AppList 获取应用列表
func (s *ExtensionDockerService) AppList() (*response.DockerAppListP, error) {
	fileBody, err := util.ReadFileStringBody(global.Config.System.PanelPath + "/data/docker_app/app_list.json")
	if err != nil {
		return nil, err
	}
	var list response.DockerAppListP
	err = util.JsonStrToStruct(fileBody, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// UpdateAppList 更新应用列表
func (s *ExtensionDockerService) UpdateAppList() error {
	downloadUrl := getDownloadNode(global.Config.System.CloudNodes[0], global.Config.System.CloudNodes)
	dataDir := global.Config.System.PanelPath + "/data/docker_app"
	_ = os.MkdirAll(dataDir, 0655)
	tarName := "docker_app_v1.tar.gz"
	Shell := `
#!/bin/bash

# 设置变量
DOWNLOAD_URL="{{downloadUrl}}/docker_app/{{tarName}}"
TEMP_DIR="/tmp"
DATA_DIR="{{dataDir}}"
TEMP_ARCHIVE="$TEMP_DIR/{{tarName}}"

# 下载压缩包
echo "正在下载压缩包..."
if ! wget -q "$DOWNLOAD_URL" -O "$TEMP_ARCHIVE"; then
    echo "下载失败，请检查链接或网络连接。"
    exit 1
fi

# 校验文件完整性
echo "正在校验文件完整性..."
if ! tar tzf "$TEMP_ARCHIVE" >/dev/null; then
    echo "压缩包文件损坏或格式不正确。"
    rm -f "$TEMP_ARCHIVE"
    exit 1
fi

# 清空目标目录
echo "清空目标目录 $DATA_DIR ..."
rm -rf "$DATA_DIR"/*

# 解压压缩包到目标目录
echo "正在解压压缩包到 $DATA_DIR ..."
tar zxvf "$TEMP_ARCHIVE" -C "$DATA_DIR"

# 删除压缩包
echo "删除压缩包..."
rm -f "$TEMP_ARCHIVE"

echo "更新完成。"

`
	Shell = strings.Replace(Shell, "{{downloadUrl}}", downloadUrl, -1)
	Shell = strings.Replace(Shell, "{{dataDir}}", dataDir, -1)
	Shell = strings.Replace(Shell, "{{tarName}}", tarName, -1)

	_, err := util.ExecShellScript(Shell)
	if err != nil {
		return err
	}
	return nil
}

// ComposeConfig 获取应用配置
func (s *ExtensionDockerService) ComposeConfig(composePath string) (*response.DockerAppConfigP, error) {
	configFileBody, err := util.ReadFileStringBody(composePath + "/env_params.json")
	if err != nil {
		return nil, err
	}
	dockerComposeFileBody, err := util.ReadFileStringBody(composePath + "/docker-compose.yml")
	if err != nil {
		return nil, err
	}
	var params []response.DockerAppEnvParam
	err = util.JsonStrToStruct(configFileBody, &params)
	if err != nil {
		return nil, err
	}
	data := &response.DockerAppConfigP{
		DockerCompose: dockerComposeFileBody,
		Params:        params,
	}
	return data, nil
}

// SaveComposeConfig 保存应用配置
func (s *ExtensionDockerService) SaveComposeConfig(param *request.SaveComposeConfigR) error {
	//检查端口是否被占用
	if _, ok := param.Params["HTTP_PORT"]; ok {
		if util.CheckPortOccupied("tcp", int(param.Params["HTTP_PORT"].(float64))) {
			return errors.New(fmt.Sprintf("port %v is occupied", param.Params["HTTP_PORT"]))
		}
	}
	//生成新的配置
	return s.WriteComposeConfig(param.ComposePath, param.DockerCompose, param.Params)
}

// GetComposeServices 获取compose服务列表
func (s *ExtensionDockerService) GetComposeServices(composePath string) ([]string, error) {
	// 初始化 Viper
	viper.SetConfigFile(composePath + "/docker-compose.yml")
	viper.SetConfigType("yaml")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading the configuration file: %v", err))
	}

	// 获取服务列表
	var servicesList []string
	services := viper.GetStringMap("services")
	for serviceName := range services {
		servicesList = append(servicesList, serviceName)
	}
	return servicesList, nil
}

// WriteComposeConfig 写入应用配置
func (s *ExtensionDockerService) WriteComposeConfig(composePath string, DockerCompose string, params map[string]any) error {
	//生成新的变量文件
	var envStr string
	for k, v := range params {
		envStr += fmt.Sprintf("%s=%v\n", k, v)
	}
	err := util.WriteFile(composePath+"/.env", []byte(envStr), 0755)
	if err != nil {
		return err
	}

	//重写docker-compose.yml
	err = util.WriteFile(composePath+"/docker-compose.yml", []byte(DockerCompose), 0755)

	return nil
}

// DeployApp 部署应用
func (s *ExtensionDockerService) DeployApp(param *request.DockerDeployAppR) (string, error) {
	//检查compose是否存在
	compose, err := (&model.DockerCompose{Name: param.AppName}).Get(global.PanelDB)
	if err != nil {
		return "", err
	}
	if compose.ID > 0 {
		return "", errors.New("compose already exists")
	}
	//检查端口是否被占用
	if _, ok := param.Params["HTTP_PORT"]; ok {
		if util.CheckPortOccupied("tcp", int(param.Params["HTTP_PORT"].(float64))) {
			return "", errors.New(fmt.Sprintf("port %v is occupied", param.Params["HTTP_PORT"]))
		}
	}

	param.AppPath = filepath.Clean(param.AppPath)
	_ = os.MkdirAll(param.AppPath, 0755)

	//复制配置信息到compose_app目录
	cmdStr := fmt.Sprintf("cp -rf %s/data/docker_app/%s/* %s/", global.Config.System.PanelPath, param.ServerName, param.AppPath)
	_, err = util.ExecShell(cmdStr)
	if err != nil {
		return "", err
	}

	//生成新的配置
	err = s.WriteComposeConfig(param.AppPath, param.DockerCompose, param.Params)
	if err != nil {
		return "", err
	}

	//创建compose
	compose = &model.DockerCompose{
		Name:       param.AppName,
		ServerName: param.ServerName,
		Path:       param.AppPath,
		Remark:     "From Docker App",
	}
	if _, err = compose.Create(global.PanelDB); err != nil {
		return "", err
	}
	return s.BuildCompose(param.AppName, param.AppPath+"/docker-compose.yml")
}

// PullImage 拉取镜像
func (s *ExtensionDockerService) PullImage(param *request.DockerPullImageR) (string, error) {
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return "", err
	}

	//检查是否存在
	get, err := (&model.DockerRepository{ID: param.SourceId}).Get(global.PanelDB)
	if err != nil {
		return "", err
	}
	if get.ID == 0 {
		return "", errors.New("not found DockerRepository")
	}

	//如果有用户名密码则构造登录信息
	imagePullOptions := types.ImagePullOptions{}
	if !util.StrIsEmpty(get.Username) {
		authConfig := registry.AuthConfig{
			Username: get.Username,
			Password: get.Password,
		}
		imagePullOptions.RegistryAuth = util.EncodeAuthToBase64(authConfig)
	}

	//开始拉取镜像
	logPath := fmt.Sprintf("%s/docker/pull_image_%s.log", global.Config.Logger.RootPath, time.Now().Format("20060102150405"))
	err = os.MkdirAll(path.Dir(logPath), 0777)
	if err != nil {
		return "", err
	}
	imageName := fmt.Sprintf("%s/%s", get.Url, param.Name)
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		return "", err
	}
	go func() {
		defer func(logFile *os.File) {
			_ = logFile.Close()
		}(logFile)

		_, _ = logFile.WriteString(fmt.Sprintf("Start pulling image: [%v]\n", imageName) + util.GetCmdDelimiter2())
		out, err := dockerClient.ImagePull(context.Background(), imageName, imagePullOptions)
		if err != nil {
			_, _ = logFile.WriteString(fmt.Sprintf("Pull image [%v] failed, error: [%v]\n", imageName, err.Error()))
			global.Log.Errorf("pull image %s failed! error: %s", imageName, err.Error())
			return
		}
		defer func(out io.ReadCloser) {
			_ = out.Close()

		}(out)
		_, _ = io.Copy(logFile, out)

		_, _ = logFile.WriteString(fmt.Sprintf("%sPull image [%v] successfully \n %s", util.GetCmdDelimiter2(), imageName, util.GetCmdDelimiter2()))
		global.Log.Infof("pull image %s successfully!", imageName)
	}()
	return logPath, nil
}

// ImportImage 导入镜像
func (s *ExtensionDockerService) ImportImage(path string) error {
	//检查文件是否存在
	if !util.PathExists(path) {
		return errors.New(fmt.Sprintf("Local image file [%v] does not exist", path))
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return err
	}
	res, err := dockerClient.ImageLoad(context.Background(), file, true)
	if err != nil {
		return err
	}
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if strings.Contains(string(content), "Error") {
		return errors.New(string(content))
	}
	return nil
}

// BuildImage 构建镜像
func (s *ExtensionDockerService) BuildImage(param *request.DockerBuildImageR) (string, error) {
	var buildDir string
	var buildFileName string
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return "", err
	}

	if util.StrIsEmpty(param.Path) { //如果没有指定构建目录则创建一个
		//创建构建目录
		buildDirPath := fmt.Sprintf("%s/data/extensions/docker/build/%s", global.Config.System.PanelPath, time.Now().Format("20060102150405"))
		err = os.MkdirAll(buildDirPath, 0777)
		if err != nil {
			return "", err
		}

		buildFilePath := fmt.Sprintf("%s/Dockerfile", buildDirPath)
		file, err := os.OpenFile(buildFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return "", err
		}
		defer func(file *os.File) {
			_ = file.Close()

		}(file)
		write := bufio.NewWriter(file)
		_, _ = write.WriteString(param.DockerFileBody)
		_ = write.Flush()

		buildDir = buildDirPath
		buildFileName = "Dockerfile"
	} else {
		buildDir = path.Dir(param.Path)
		buildFileName = path.Base(param.Path)
	}

	arch, err := archive.TarWithOptions(buildDir+"/", &archive.TarOptions{})
	if err != nil {
		return "", err
	}

	buildOptions := types.ImageBuildOptions{
		Dockerfile: buildFileName,
		Tags:       param.Tags,
		Remove:     true,
		Labels:     param.Labels,
	}
	logPath := fmt.Sprintf("%s/docker/build_image_%s.log", global.Config.Logger.RootPath, time.Now().Format("20060102150405"))
	err = os.MkdirAll(path.Dir(logPath), 0777)
	if err != nil {
		return "", err
	}
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", err
	}

	_, _ = file.WriteString(fmt.Sprintf("Start building image Tags: [%v], Dockerfile: [%v/%v] \n %s", param.Tags, buildDir, buildFileName, util.GetCmdDelimiter2()))
	go func() {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		defer func(arch io.ReadCloser) {
			_ = arch.Close()
		}(arch)
		res, err := dockerClient.ImageBuild(context.Background(), arch, buildOptions)
		if err != nil {
			global.Log.Errorf("build image %s failed, err: %v", param.Tags, err)
			_, _ = file.WriteString(util.GetCmdDelimiter2())
			_, _ = file.WriteString(fmt.Sprintf("Build image failed Tags: %v, Dockerfile: %v, Error: %v \n", param.Tags, buildDir+buildFileName, err.Error()))
			return
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			global.Log.Errorf("build image %s failed, err: %v", param.Tags, err)
			_, _ = file.WriteString(util.GetCmdDelimiter2())
			_, _ = file.WriteString(fmt.Sprintf("Build image failed Tags: %v, Dockerfile: %v, Error: %v \n", param.Tags, buildDir+buildFileName, err.Error()))
			return
		}

		// 判断构建镜像是否成功
		if strings.Contains(string(body), "Successfully built") {
			// 获取构建镜像的ID
			re := regexp.MustCompile(`Successfully built ([a-z0-9]+)`)
			matches := re.FindStringSubmatch(string(body))
			if len(matches) >= 2 {
				imageID := matches[1]
				// 构建成功
				global.Log.Infof("build image %s Successfully,imageID:%v", param.Tags, imageID)
				_, _ = file.WriteString(fmt.Sprintf("%sBuild image successfully imageID：[%v] \n", util.GetCmdDelimiter2(), imageID))
				return
			}
		} else {
			global.Log.Errorf("build image %s failed, err: %v", param.Tags, err)
			_, _ = file.WriteString(util.GetCmdDelimiter2())
			_, _ = file.WriteString(fmt.Sprintf("Build image failed Tags: %v, Dockerfile: %v, Error: %v \n", param.Tags, buildDir+buildFileName, err.Error()))
			return
		}
	}()

	return logPath, nil
}

// PushImage 推送镜像
func (s *ExtensionDockerService) PushImage(param *request.DockerPushImageR) (string, error) {
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return "", err
	}

	//检查是否是默认仓库
	if param.RepositoryId == 1 {
		return "", errors.New("do not allow pushing to the default repository")
	}

	//获取仓库信息
	RepositoryInfo, err := (&model.DockerRepository{ID: param.RepositoryId}).Get(global.PanelDB)
	if err != nil {
		return "", err
	}
	if RepositoryInfo.ID == 0 {
		return "", errors.New("not found DockerRepository")
	}

	//如果有用户名密码则构造登录信息
	imagePushOptions := types.ImagePushOptions{}
	if !util.StrIsEmpty(RepositoryInfo.Username) && !util.StrIsEmpty(RepositoryInfo.Password) {
		authConfig := registry.AuthConfig{
			Username: RepositoryInfo.Username,
			Password: RepositoryInfo.Password,
		}
		imagePushOptions.RegistryAuth = util.EncodeAuthToBase64(authConfig)
	}

	logPath := fmt.Sprintf("%s/docker/push_image_%s.log", global.Config.Logger.RootPath, time.Now().Format("20060102150405"))
	err = os.MkdirAll(path.Dir(logPath), 0777)
	if err != nil {
		return "", err
	}
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", err
	}
	image := fmt.Sprintf("%s/%s/%s", RepositoryInfo.Url, RepositoryInfo.Namespace, param.Tag)

	_, _ = file.WriteString(fmt.Sprintf("Start pushing image image_id: [%v] \n%s", param.ImageId, util.GetCmdDelimiter2()))
	go func() {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		output, err := dockerClient.ImagePush(context.Background(), image, imagePushOptions)
		if err != nil {
			global.Log.Errorf("push image %s failed, err: %v", param.ImageId, err)
			_, _ = file.WriteString(fmt.Sprintf("%sPush image failed image_id: %v, Error: %v \n%s", util.GetCmdDelimiter2(), param.ImageId, err.Error(), util.GetCmdDelimiter2()))
			return
		}
		defer func(output io.ReadCloser) {
			_ = output.Close()
		}(output)
		_, _ = io.Copy(file, output)
		// 推送成功
		global.Log.Infof("push image  Successfully,imageID:%v", param.ImageId)
		_, _ = file.WriteString(fmt.Sprintf("%s推送镜像成功 image_id: %v \n%s", util.GetCmdDelimiter2(), param.ImageId, util.GetCmdDelimiter2()))
		return

	}()

	return logPath, nil
}

// ExportImage 导出镜像
func (s *ExtensionDockerService) ExportImage(param *request.DockerExportImageR) (string, error) {
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return "", err
	}
	savePath := fmt.Sprintf("%s/%s.tar", param.Path, param.Name)

	//检查路径是否存在
	if !util.PathExists(param.Path) {
		return "", errors.New("does not exist")
	}
	//检查文件是否存在
	if util.PathExists(savePath) {
		savePath = fmt.Sprintf("%s/%s_%s.tar", param.Path, param.Name, time.Now().Format("20060102150405"))
	}

	output, err := dockerClient.ImageSave(context.Background(), []string{param.ImageId})
	if err != nil {
		return "", err
	}
	defer func(output io.ReadCloser) {
		_ = output.Close()
	}(output)
	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0755)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if _, err = io.Copy(file, output); err != nil {
		return "", err
	}
	return savePath, nil
}

// DeleteImage 删除镜像
func (s *ExtensionDockerService) DeleteImage(imageID string) error {
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return err
	}
	_, err = dockerClient.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
	if err != nil {
		return err
	}
	return nil
}

// ImageList 镜像列表
func (s *ExtensionDockerService) ImageList() ([]*response.DockerImageListP, int, error) {
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return nil, 0, err
	}

	//获取镜像列表
	images, err := dockerClient.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil, 0, err
	}

	//获取镜像数量
	total := len(images)

	//返回数据
	var list []*response.DockerImageListP
	for _, image := range images {
		var tags []string
		for _, tag := range image.RepoTags {
			tags = append(tags, strings.Split(tag, ":")[1])
		}
		list = append(list, &response.DockerImageListP{
			Id:         image.ID,
			CreateTime: image.Created,
			Size:       image.Size,
			Tags:       tags,
			Detail:     image,
		})
	}
	return list, total, nil
}

// AddRepository 添加docker仓库
func (s *ExtensionDockerService) AddRepository(param *request.AddRepositoryR) error {

	//如果有用户名密码则验证登录信息
	if !util.StrIsEmpty(param.Username) {
		_, err := s.LoginRepository(param.Url, param.Username, param.Password)
		if err != nil {
			return err
		}
	}

	//检查是否存在
	getData, err := (&model.DockerRepository{Name: param.Name}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if getData.ID > 0 {
		return errors.New(helper.MessageWithMap("docker.repository.NameExist", map[string]any{"Name": param.Name}))
	}

	//添加仓库
	_, err = (&model.DockerRepository{
		Name:      param.Name,
		Url:       param.Url,
		Username:  param.Username,
		Password:  param.Password,
		Namespace: param.Namespace,
		Remark:    param.Remark,
	}).Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// RepositoryList 仓库列表
func (s *ExtensionDockerService) RepositoryList(param *request.RepositoryListR) ([]*model.DockerRepository, int64, error) {
	whereT := model.ConditionsT{"ORDER": "create_time DESC"}
	whereOrT := model.ConditionsT{}
	if !util.StrIsEmpty(param.Query) {
		param.Query = "%" + param.Query + "%"
		whereT["name LIKE ?"] = param.Query
		//whereOrT["FIXED"] = map[string]interface{}{"url LIKE ?": param.Query, "remark LIKE ?": param.Query}
		whereOrT["url LIKE ?"] = param.Query
		whereOrT["remark LIKE ?"] = param.Query
	}
	list, total, err := (&model.DockerRepository{}).List(global.PanelDB, &whereT, &whereOrT, 0, 0)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// EditRepository 编辑仓库
func (s *ExtensionDockerService) EditRepository(param *request.EditRepositoryR) error {
	//检查是否是默认仓库
	if param.Id == 1 {
		return errors.New("default repository is not allowed to be edited")
	}

	//如果有用户名密码则验证登录信息
	if !util.StrIsEmpty(param.Username) {
		_, err := s.LoginRepository(param.Url, param.Username, param.Password)
		if err != nil {
			return err
		}
	}

	//检查是否存在
	get, err := (&model.DockerRepository{
		ID: param.Id,
	}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if get.ID == 0 {
		return errors.New("not found DockerRepository")
	}

	//更新仓库
	err = (&model.DockerRepository{
		ID:        param.Id,
		Name:      param.Name,
		Url:       param.Url,
		Username:  param.Username,
		Password:  param.Password,
		Namespace: param.Namespace,
		Remark:    param.Remark,
	}).Update(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRepository 删除仓库
func (s *ExtensionDockerService) DeleteRepository(id int64) (*model.DockerRepository, error) {
	//检查是否是默认仓库
	if id == 1 {
		return nil, errors.New("default repository is not allowed to be deleted")
	}

	//检查是否存在
	get, err := (&model.DockerRepository{
		ID: id,
	}).Get(global.PanelDB)
	if err != nil {
		return nil, err
	}
	if get.ID == 0 {
		return nil, errors.New("not found DockerRepository")
	}

	//删除仓库
	err = (&model.DockerRepository{
		ID: get.ID,
	}).Delete(global.PanelDB)
	if err != nil {
		return nil, err
	}
	return get, nil
}

// Install Docker
func (s *ExtensionDockerService) Install() error {
	//检查是否在等待或者进行队列中
	_, total, err := (&model.QueueTask{}).List(global.PanelDB, &model.ConditionsT{
		"FIXED": "status = " + fmt.Sprintf("%d", constant.QueueTaskStatusProcessing) + " OR status = " + fmt.Sprintf("%d", constant.QueueTaskStatusWait),
		"name":  "安装[Docker]",
	}, 0, 0)
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New(helper.MessageWithMap("queuetask.InstallingOrWaiting", map[string]any{"Name": "安装[Docker]"}))
	}

	//添加面板队列任务
	_, err = (&model.QueueTask{
		Name:    "安装[Docker]",
		Type:    1,
		Status:  constant.QueueTaskStatusWait,
		ExecStr: fmt.Sprintf(`cd %s && /bin/bash %s install`, s.GetShellPath(), s.GetShellName()),
	}).Create(global.PanelDB)
	if err != nil {
		global.Log.Errorf("添加任務失敗->createTaskQueue()->ds.CreateTaskQueue()  Error:%s", err)
		return err
	}
	return nil
}

// Uninstall Docker
func (s *ExtensionDockerService) Uninstall() error {
	//检查是否在等待或者进行队列中
	taskName := "卸载[Docker]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.UninstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash %s uninstall`, s.GetShellPath(), s.GetShellName())
	err := AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}

	return nil
}

// SetStatus 设置docker状态
func (s *ExtensionDockerService) SetStatus(action string) error {
	err := s.InitCheck()
	if err != nil {
		return err
	}
	var cmdStr string
	switch action {
	case constant.ProcessCommandByStart:
		cmdStr = "systemctl start docker"
	case constant.ProcessCommandByStop:
		cmdStr = "systemctl stop docker && systemctl stop docker.socket"
	case constant.ProcessCommandByRestart:
		cmdStr = "systemctl restart docker"
	default:
		return errors.New("action error")
	}
	_, err = util.ExecShell(cmdStr)
	if err != nil {
		return err
	}
	return nil
}

// RunInfo 运行信息
func (s *ExtensionDockerService) RunInfo() (types.Info, error) {
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return types.Info{}, err
	}
	return dockerClient.Info(context.Background())
}

// GetDockerVersion 获取docker的版本号
func (s *ExtensionDockerService) GetDockerVersion() string {
	dockerVersion, err := util.ExecShell("docker version --format '{{.Server.Version}}' | tr -d '\\n'")
	if err != nil {
		return ""
	}
	return dockerVersion
}

// GetDockerComposeVersion 获取docker-compose的版本号
func (s *ExtensionDockerService) GetDockerComposeVersion() string {
	dockerComposeVersion, err := util.ExecShell("docker-compose version --short | tr -d '\\n'")
	if err != nil {
		return ""
	}
	return dockerComposeVersion
}

func (s *ExtensionDockerService) GetShellPath() string {
	return global.Config.System.PanelPath + "/data/extensions/docker/install"
}

func (s *ExtensionDockerService) GetShellName() any {
	return "install_docker.sh"
}

// IsRunning 获取docker运行状态
func (s *ExtensionDockerService) IsRunning() bool {
	output, _ := util.ExecShell("systemctl status docker")
	if strings.Contains(output, "Active: active (running)") {
		return true
	} else if strings.Contains(output, "Active: inactive (dead)") {
		return false
	}
	return false
}

//// CheckDockerInstalled 检查docker是否安装
//func (s *ExtensionDockerService) CheckDockerInstalled() error {
//
//	return nil
//}
//
//// CheckDockerComposeInstalled 检查docker-compose是否安装
//func (s *ExtensionDockerService) CheckDockerComposeInstalled() error {
//
//	return nil
//}

// InitCheck 初始化检查
func (s *ExtensionDockerService) InitCheck() error {
	if !s.CheckDockerInstalled() {
		return errors.New("docker Not Installed")
	}
	if !s.CheckDockerComposeInstalled() {
		return errors.New("docker-compose Not Installed")
	}
	return nil
}

// CheckDockerInstalled 检查docker是否已经安装
func (s *ExtensionDockerService) CheckDockerInstalled() bool {
	_, err := util.ExecShell("docker --version")
	if err != nil {
		return false
	}
	return true
}

// CheckDockerComposeInstalled 检查docker-compose是否已经安装
func (s *ExtensionDockerService) CheckDockerComposeInstalled() bool {
	_, err := util.ExecShell("docker-compose --version")
	if err != nil {
		return false
	}
	return true
}

// LoginRepository 登录仓库
func (s *ExtensionDockerService) LoginRepository(url, username, password string) (*registry.AuthenticateOKBody, error) {
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return nil, err
	}
	authConfig := registry.AuthConfig{
		ServerAddress: url,
		Username:      username,
		Password:      password,
	}
	body, err := dockerClient.RegistryLogin(context.Background(), authConfig)
	if err != nil {
		return nil, errors.New(helper.MessageWithMap("docker.repository.LoginError", map[string]any{"Err": err.Error()}))
	}
	return &body, nil
}

// CreateComposeTemplate 创建docker-compose模板
func (s *ExtensionDockerService) CreateComposeTemplate(param *request.CreateComposeTemplateR) error {
	//校验名称
	if !util.IsValidFileName(param.Name) {
		return errors.New("name is invalid")
	}

	//检查模板名称是否存在
	get, err := (&model.DockerComposeTemplate{Name: param.Name}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if get.ID > 0 {
		return errors.New("name already exists")
	}

	var insert *model.DockerComposeTemplate

	switch param.Type {
	case "create": //新建模板
		saveDir := fmt.Sprintf("%s/data/extensions/docker/compose/template", global.Config.System.PanelPath)
		savePath := fmt.Sprintf("%s/%s.yaml", saveDir, param.Name)
		err = os.MkdirAll(saveDir, 0755)
		if err != nil {
			return err
		}
		err = util.WriteFile(savePath, []byte(param.Body), 0755)
		if err != nil {
			return err
		}

		//检查模板
		err = s.CheckComposeTemplate(savePath)
		if err != nil {
			//删除文件
			_ = os.Remove(savePath)
			return err
		}
		insert = &model.DockerComposeTemplate{
			Name:   param.Name,
			Remark: param.Remark,
			Path:   savePath,
		}
	case "local": //加载本地模板
		//检查path是否存在
		if !util.PathExists(param.Path) {
			return errors.New(fmt.Sprintf("not found path:%s", param.Path))
		}

		//检查模板
		err = s.CheckComposeTemplate(param.Path)
		if err != nil {
			return err
		}

		// 读取源文件的内容
		input, err := os.ReadFile(param.Path)
		if err != nil {
			log.Fatal(err)
		}

		//写入新文件
		saveDir := fmt.Sprintf("%s/data/extensions/docker/compose/template", global.Config.System.PanelPath)
		savePath := fmt.Sprintf("%s/%s.yaml", saveDir, param.Name)
		err = os.MkdirAll(saveDir, 0755)
		if err != nil {
			return err
		}
		err = util.WriteFile(savePath, input, 0755)
		if err != nil {
			return err
		}
		insert = &model.DockerComposeTemplate{
			Name:   param.Name,
			Remark: param.Remark,
			Path:   savePath,
		}
	default:
		return errors.New("six six six")
	}

	_, err = (insert).Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// ComposeTemplateList 获取模板列表
func (s *ExtensionDockerService) ComposeTemplateList(query string, offset, limit int) ([]*model.DockerComposeTemplate, int64, error) {
	where := model.ConditionsT{"ORDER": "create_time DESC"}
	whereOr := model.ConditionsT{}

	if !util.StrIsEmpty(query) {
		query = fmt.Sprintf("%%%s%%", query)
		where = model.ConditionsT{
			"name LIKE ?": query,
		}
		whereOr = model.ConditionsT{
			"remark LIKE ?": query,
			"path LIKE ?":   query,
		}
	}
	list, total, err := (&model.DockerComposeTemplate{}).List(global.PanelDB, &where, &whereOr, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// DeleteComposeTemplate 删除模板
func (s *ExtensionDockerService) DeleteComposeTemplate(id int64) (*model.DockerComposeTemplate, error) {
	//检查模板是否存在
	get, err := (&model.DockerComposeTemplate{ID: id}).Get(global.PanelDB)
	if err != nil {
		return nil, err
	}
	if get.ID == 0 {
		return nil, errors.New("not found Template")
	}

	//删除文件
	_ = os.Remove(get.Path)

	//删除数据库记录
	err = (&model.DockerComposeTemplate{ID: id}).Delete(global.PanelDB)
	if err != nil {
		return nil, err
	}
	return get, nil
}

// GetComposeTemplateByID 通过ID获取模板信息
func (s *ExtensionDockerService) GetComposeTemplateByID(id int64) (*model.DockerComposeTemplate, error) {
	template, err := (&model.DockerComposeTemplate{ID: id}).Get(global.PanelDB)
	if err != nil {
		return nil, err
	}
	if template.ID == 0 {
		return nil, errors.New("not found Template")
	}
	return template, nil
}

// EditComposeTemplate 编辑模板
func (s *ExtensionDockerService) EditComposeTemplate(template *model.DockerComposeTemplate, param *request.EditComposeTemplateR) error {
	//校验名称
	if !util.IsValidFileName(param.Name) {
		return errors.New("name is invalid")
	}

	//检查模板名称是否存在
	get, err := (&model.DockerComposeTemplate{Name: param.Name}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if get.ID > 0 && get.ID != param.Id {
		return errors.New("name already exists")
	}

	//备份文件内容
	backUpBody, err := util.ReadFileStringBody(template.Path)
	if err != nil {
		return err
	}

	//写入新内容
	err = util.WriteFile(template.Path, []byte(param.Body), 0755)
	if err != nil {
		return err
	}

	//检查模板
	err = s.CheckComposeTemplate(template.Path)
	if err != nil {
		//恢复文件内容
		_ = util.WriteFile(template.Path, []byte(backUpBody), 0755)
		return err
	}

	//更新数据库
	template.Name = param.Name
	template.Remark = param.Remark
	err = template.Update(global.PanelDB)
	if err != nil {
		//恢复文件内容
		_ = util.WriteFile(template.Path, []byte(backUpBody), 0755)
		return err
	}
	return nil
}

// ComposeTemplatePullImage 拉取镜像
func (s *ExtensionDockerService) ComposeTemplatePullImage(id int64) (string, error) {
	//检查模板是否存在
	template, err := (&model.DockerComposeTemplate{ID: id}).Get(global.PanelDB)
	if err != nil {
		return "", err
	}
	if template.ID == 0 {
		return "", errors.New("not found Template")
	}

	logPath := fmt.Sprintf("%s/docker/docker_compose_pull_image_%d_%s.log", global.Config.Logger.RootPath, template.ID, time.Now().Format("20060102150405"))
	err = os.MkdirAll(path.Dir(logPath), 0777)
	if err != nil {
		return "", err
	}
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		return "", err
	}

	go func() {
		defer func(logFile *os.File) {
			_ = logFile.Close()
		}(logFile)
		_, _ = logFile.WriteString(fmt.Sprintf("Start pulling compose template image, template location:[%v] \n%s", template.Path, util.GetCmdDelimiter2()))
		cmdStr := fmt.Sprintf("docker-compose -f %s pull >> %s 2>&1 ", template.Path, logPath)
		_, err = util.ExecShell(cmdStr)
		if err != nil {
			_, _ = logFile.WriteString(fmt.Sprintf("DockerComposeTemplate pull image [%v] failed, error: [%v]\n", template.Path, err.Error()))
			global.Log.Errorf("ComposeTemplatePullImage %s failed! error: %s", template.Path, err.Error())
			return
		}
		_, _ = logFile.WriteString(fmt.Sprintf("%sComposeTemplate Pull image [%v] successfully \n %s", util.GetCmdDelimiter2(), template.Path, util.GetCmdDelimiter2()))
		global.Log.Infof("ComposeTemplatePullImage %s successfully!", template.Path)
	}()

	return logPath, nil
}

// CheckComposeTemplate 检查模板是否正确
func (s *ExtensionDockerService) CheckComposeTemplate(path string) error {
	cmdStr := fmt.Sprintf("docker-compose -f %s config", path)
	out, err := util.ExecShell(cmdStr)
	if err != nil {
		return errors.New(out)
	}
	return nil
}

// CreateContainer	创建容器
func (s *ExtensionDockerService) CreateContainer(param *request.CreateContainerR) error {
	//是否暴露端口
	if len(param.Port) > 0 {
		for _, port := range param.Port {
			//检测服务器端口是否被占用
			if util.CheckPortOccupied("tcp", port.HostPort) {
				return errors.New(helper.MessageWithMap("PortOccupied", map[string]any{"Port": port.HostPort}))
			}
		}
	}

	//检查cpu和内存是否超出限制
	if c, _ := cpu.Counts(false); c < param.NanoCPUs {
		return errors.New(helper.MessageWithMap("docker.container.CpuExceed", map[string]any{"Set": param.NanoCPUs, "Get": c}))
	}
	if m, _ := mem.VirtualMemory(); int64(m.Total) < param.Memory {
		return errors.New(helper.MessageWithMap("docker.container.MemoryExceed", map[string]any{"Set": param.NanoCPUs, "Get": m}))
	}

	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return err
	}

	config := &container.Config{
		Image:  param.Image,
		Cmd:    param.Cmd,
		Env:    param.Env,
		Labels: param.Labels,
	}
	hostConfig := &container.HostConfig{
		Resources:       container.Resources{NanoCPUs: int64(param.NanoCPUs) * 1e9, Memory: param.Memory},
		AutoRemove:      param.AutoRemove,
		PublishAllPorts: len(param.Port) == 0,
		RestartPolicy:   container.RestartPolicy{Name: param.RestartPolicy},
	}
	//如果是on-failure策略，最大重试次数为5
	if param.RestartPolicy == "on-failure" {
		hostConfig.RestartPolicy.MaximumRetryCount = 5
	}

	if len(param.Port) > 0 {
		hostConfig.PortBindings = make(nat.PortMap)
		for _, port := range param.Port {
			bindItem := nat.PortBinding{HostPort: strconv.Itoa(port.HostPort)}
			hostConfig.PortBindings[nat.Port(fmt.Sprintf("%d/%s", port.ContainerPort, port.Protocol))] = []nat.PortBinding{bindItem}
		}
	}
	if len(param.Volumes) > 0 {
		config.Volumes = make(map[string]struct{})
		for _, volumeInfo := range param.Volumes {
			config.Volumes[volumeInfo.ContainerDir] = struct{}{}
			hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%s:%s:%s", volumeInfo.HostDir, volumeInfo.ContainerDir, volumeInfo.HostDirPermissions))
		}
	}

	containers, err := dockerClient.ContainerCreate(
		context.Background(),
		config,
		hostConfig,
		&network.NetworkingConfig{},
		&v1.Platform{},
		param.Name,
	)
	if err != nil {
		_ = dockerClient.ContainerRemove(
			context.Background(),
			param.Name,
			types.ContainerRemoveOptions{RemoveVolumes: true, Force: true},
		)
		return err
	}
	global.Log.Infof("create container %s successful!", param.Name)

	//尝试启动容器
	if err = dockerClient.ContainerStart(context.Background(), containers.ID, types.ContainerStartOptions{}); err != nil {
		_ = dockerClient.ContainerRemove(
			context.Background(),
			param.Name,
			types.ContainerRemoveOptions{RemoveVolumes: true, Force: true},
		)
		return errors.New(fmt.Sprintf("Failed to start container[%s],Error:%s", param.Image, err.Error()))
	}
	return nil
}

// ContainerList	获取容器列表
func (s *ExtensionDockerService) ContainerList(filterS string, getStats bool) ([]response.ContainerInfo, error) {
	//[]types.Container
	dataList := make([]response.ContainerInfo, 0)
	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return nil, err
	}

	options := types.ContainerListOptions{All: true}
	fmt.Println(filterS)
	if !util.StrIsEmpty(filterS) {
		options.Filters = filters.NewArgs()
		options.Filters.Add("label", filterS)
	}
	list, err := dockerClient.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	for _, containerInfo := range list {
		var data response.ContainerInfo
		data.Container = containerInfo
		if getStats {
			res, err := dockerClient.ContainerStats(context.Background(), containerInfo.ID, false)
			if err != nil {
				return nil, err
			}
			body, err := io.ReadAll(res.Body)
			if err != nil {
				_ = res.Body.Close()
				return nil, err
			}
			_ = res.Body.Close()
			var stats types.StatsJSON
			if err := json.Unmarshal(body, &stats); err != nil {
				return nil, err
			}
			data.Stats = &stats
		}
		dataList = append(dataList, data)
	}
	return dataList, nil
}

// OperateContainer	操作容器
func (s *ExtensionDockerService) OperateContainer(name string, newName string, action string) error {
	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return err
	}
	//start stop reboot kill pause recover remove
	ctx := context.Background()
	switch action {
	case "start":
		return dockerClient.ContainerStart(ctx, name, types.ContainerStartOptions{})
	case "stop":
		return dockerClient.ContainerStop(ctx, name, container.StopOptions{})
	case "reboot":
		return dockerClient.ContainerRestart(ctx, name, container.StopOptions{})
	case "kill":
		return dockerClient.ContainerKill(ctx, name, "SIGKILL")
	case "pause":
		return dockerClient.ContainerPause(ctx, name)
	case "recover":
		return dockerClient.ContainerUnpause(ctx, name)
	case "remove":
		return dockerClient.ContainerRemove(ctx, name, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true})
	case "rename":
		if util.StrIsEmpty(newName) {
			return errors.New("newName is empty")
		}
		return dockerClient.ContainerRename(context.Background(), name, newName)
	default:
		return errors.New("action is empty")
	}
}

// ContainerLogs	获取容器日志
func (s *ExtensionDockerService) ContainerLogs(name, since, until string, tail int) (string, error) {
	cmdStr := fmt.Sprintf("docker logs %s", name)
	if !util.StrIsEmpty(since) {
		cmdStr = fmt.Sprintf("%s --since %s", cmdStr, since)
	}
	if !util.StrIsEmpty(until) {
		cmdStr = fmt.Sprintf("%s --until %s", cmdStr, until)
	}
	if tail > 0 {
		cmdStr = fmt.Sprintf("%s --tail %d", cmdStr, tail)
	}
	output, err := util.ExecShell(cmdStr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("ERROR: %s ,ERROR_MSG: %s", output, err))
	}
	return output, nil
}

// ContainerMonitor	容器监控
func (s *ExtensionDockerService) ContainerMonitor(id string) (*types.StatsJSON, error) {
	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return nil, err
	}
	stats, err := dockerClient.ContainerStats(context.Background(), id, false)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(stats.Body)

	body, err := io.ReadAll(stats.Body)
	if err != nil {
		return nil, err
	}
	var statsJson types.StatsJSON
	if err = json.Unmarshal(body, &statsJson); err != nil {
		return nil, err
	}
	return &statsJson, nil

}

// CreateCompose	创建compose
func (s *ExtensionDockerService) CreateCompose(param *request.CreateComposeR) (string, error) {
	//检查compose是否存在
	compose, err := (&model.DockerCompose{Name: param.Name}).Get(global.PanelDB)
	if err != nil {
		return "", err
	}
	if compose.ID > 0 {
		return "", errors.New("compose already exists")
	}
	var dockerComposeFilePath string
	if util.StrIsEmpty(param.DockerCompose) { //如果为空则是使用已有项目目录进行创建
		//判断是否含有docker-compose.yml文件
		dockerComposeFilePath = fmt.Sprintf("%s/docker-compose.yml", param.Path)
		if !util.PathExists(dockerComposeFilePath) {
			return "", errors.New(fmt.Sprintf("not found docker-compose.yml in %s", param.Path))
		}
	} else {
		//写入docker_compose.yaml
		dockerComposeFilePath = fmt.Sprintf("%s/docker-compose.yml", param.Path)
		err = os.MkdirAll(path.Dir(dockerComposeFilePath), 0755)
		if err != nil {
			return "", err
		}
		err = util.WriteFile(dockerComposeFilePath, []byte(param.DockerCompose), 0755)
		if err != nil {
			return "", err
		}
	}

	//创建compose
	compose = &model.DockerCompose{
		Name:   param.Name,
		Path:   param.Path,
		Remark: param.Remark,
	}
	if _, err = compose.Create(global.PanelDB); err != nil {
		return "", err
	}

	return s.BuildCompose(param.Name, dockerComposeFilePath)
}

// BuildCompose 构建compose
func (s *ExtensionDockerService) BuildCompose(projectName, dockerComposeFilePath string) (string, error) {
	logPath := fmt.Sprintf("%s/docker/docker_compose_up_%s.log", global.Config.Logger.RootPath, time.Now().Format("20060102150405"))
	err := os.MkdirAll(path.Dir(logPath), 0755)
	if err != nil {
		return "", err
	}
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0755)
	if err != nil {
		return "", err
	}
	go func() {
		defer func(logFile *os.File) {
			_ = logFile.Close()
		}(logFile)
		_, _ = logFile.WriteString(fmt.Sprintf("Start creating compose, docker-compose.yml location:[%v] \n", dockerComposeFilePath))

		cmdStr := fmt.Sprintf("docker-compose -f %s -p %s up -d >> %s 2>&1", dockerComposeFilePath, projectName, logPath)
		_, err = util.ExecShell(cmdStr)
		if err != nil {
			_, _ = logFile.WriteString(fmt.Sprintf("Create compose failed, compose name:[%v], docker-compose.yml location:[%v], error:[%v] \n", projectName, dockerComposeFilePath, err.Error()))
			return
		}

		_, _ = logFile.WriteString(fmt.Sprintf("Create compose successfully, compose name:[%v], docker-compose.yml location:[%v] \n", projectName, dockerComposeFilePath))
		global.Log.Infof("docker-compose up %s successful!", projectName)
	}()

	return logPath, nil
}

// OperateCompose 操作compose
func (s *ExtensionDockerService) OperateCompose(path, projectName, action string, services []string) (string, error) {
	if !util.PathExists(path) {
		return "", errors.New(fmt.Sprintf("not found docker-compose.yml in %s", path))
	}
	//up down start stop restart operate
	cmdStr := fmt.Sprintf("docker-compose -f %s -p %s ", path, projectName)
	switch action {
	case "up":
		cmdStr = fmt.Sprintf("%s up -d ", cmdStr)
	case "down":
		cmdStr = fmt.Sprintf("%s down ", cmdStr)
	case "start":
		cmdStr = fmt.Sprintf("%s start ", cmdStr)
	case "stop":
		cmdStr = fmt.Sprintf("%s stop ", cmdStr)
	case "restart":
		cmdStr = fmt.Sprintf("%s restart ", cmdStr)
	default:
		return "", errors.New("action is empty")
	}
	if len(services) > 0 {
		cmdStr += strings.Join(services, " ")
	}
	//if util.StrIsEmpty(logPath) {
	//	logPath = "/dev/null"
	//}
	//cmdStr = fmt.Sprintf("%s >> %s 2>&1", cmdStr, logPath)
	result, err := util.ExecShell(cmdStr)
	if err != nil {
		global.Log.Errorf("docker-compose %s failed! cmd:%s error: %s", action, cmdStr, err.Error())
		return result, nil
	}

	return result, nil
}

// ComposeList	获取compose列表
func (s *ExtensionDockerService) ComposeList(query string, offset, limit int) ([]*model.DockerCompose, int64, error) {
	whereT := model.ConditionsT{"ORDER": "create_time DESC"}
	whereOrT := model.ConditionsT{}
	if !util.StrIsEmpty(query) {
		query = "%" + query + "%"
		whereT["name LIKE ?"] = query
		whereOrT["path LIKE ?"] = query
		whereOrT["remark LIKE ?"] = query
	}
	list, total, err := (&model.DockerCompose{}).List(global.PanelDB, &whereT, &whereOrT, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	backData := make([]*model.DockerCompose, 0)

	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return nil, 0, err
	}

	options := types.ContainerListOptions{All: true}
	options.Filters = filters.NewArgs()
	options.Filters.Add("label", "com.docker.compose.project")

	ContainerList, err := dockerClient.ContainerList(context.Background(), options)
	if err != nil {
		return nil, 0, err
	}
	composeMap := make(map[string]model.DockerCompose)
	for _, containerInfo := range ContainerList {
		labelsValue := containerInfo.Labels["com.docker.compose.project"]
		composeContainers := model.ComposeContainer{
			ContainerID: containerInfo.ID,
			Name:        containerInfo.Names[0][1:],
			State:       containerInfo.State,
			CreateTime:  time.Unix(containerInfo.Created, 0).Format("2006-01-02 15:04:05"),
		}
		if compose, has := composeMap[labelsValue]; has {
			compose.ContainerNumber++
			compose.Containers = append(compose.Containers, composeContainers)
			composeMap[labelsValue] = compose
		} else {
			composeItem := model.DockerCompose{
				ContainerNumber: 1,
				Containers:      []model.ComposeContainer{composeContainers},
			}
			composeMap[labelsValue] = composeItem
		}
	}

	for _, compose := range list {
		if composeItem, has := composeMap[compose.Name]; has {
			compose.ContainerNumber = composeItem.ContainerNumber
			compose.Containers = composeItem.Containers
		}
		backData = append(backData, compose)
	}

	return backData, total, nil

}

// GetComposeByID 根据ID获取compose
func (s *ExtensionDockerService) GetComposeByID(id int64) (compose *model.DockerCompose, err error) {
	compose, err = (&model.DockerCompose{ID: id}).Get(global.PanelDB)
	if err != nil {
		return
	}
	if compose.ID == 0 {
		err = errors.New("compose is not found")
		return
	}
	return
}

// DeleteCompose 删除compose
func (s *ExtensionDockerService) DeleteCompose(compose *model.DockerCompose) error {
	//删除compose
	if err := compose.Delete(global.PanelDB); err != nil {
		return err
	}
	_, err := s.OperateCompose(compose.Path+"/docker-compose.yml", compose.Name, "down", nil)
	if err != nil {
		return err
	}

	return nil
}

// CreateNetworking 创建网络
func (s *ExtensionDockerService) CreateNetworking(param *request.CreateNetworkingR) error {
	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return err
	}

	IPA := network.IPAMConfig{
		Subnet:  param.Subnet,
		Gateway: param.Gateway,
		IPRange: param.IpRange,
	}

	options := types.NetworkCreate{
		Driver:  param.Driver,
		Options: param.Options,
		Labels:  param.Labels,
	}
	options.IPAM = &network.IPAM{Config: []network.IPAMConfig{IPA}}
	if _, err = dockerClient.NetworkCreate(context.Background(), param.Name, options); err != nil {
		return err
	}
	return nil

}

// NetworkingList 获取网络列表
func (s *ExtensionDockerService) NetworkingList() ([]types.NetworkResource, int, error) {
	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return nil, 0, err
	}

	list, err := dockerClient.NetworkList(context.TODO(), types.NetworkListOptions{})
	if err != nil {
		return nil, 0, err
	}
	return list, len(list), nil
}

// DeleteNetworking 删除网络
func (s *ExtensionDockerService) DeleteNetworking(id string) error {
	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return err
	}

	err = dockerClient.NetworkRemove(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

// CreateVolume 创建存储
func (s *ExtensionDockerService) CreateVolume(param *request.CreateVolumeR) error {
	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return err
	}

	//检查存储是否存在
	volumeList, _, err := s.VolumeList(param.Name)
	if len(volumeList) > 0 {
		return errors.New("volume already exists")
	}

	options := volume.CreateOptions{
		Name:       param.Name,
		Driver:     param.Driver,
		DriverOpts: param.DriverOpts,
		Labels:     param.Labels,
	}
	if _, err = dockerClient.VolumeCreate(context.Background(), options); err != nil {
		return err
	}
	return nil
}

// VolumeList 获取存储列表
func (s *ExtensionDockerService) VolumeList(name string) ([]*volume.Volume, int, error) {
	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return nil, 0, err
	}
	var initialArgs []filters.KeyValuePair
	if !util.StrIsEmpty(name) {
		initialArgs = []filters.KeyValuePair{filters.Arg("name", name)}
	} else {
		initialArgs = []filters.KeyValuePair{filters.Arg("name", "")}
	}
	vos, _ := dockerClient.VolumeList(context.TODO(), volume.ListOptions{Filters: filters.NewArgs(initialArgs...)})

	return vos.Volumes, len(vos.Volumes), nil
}

// DeleteVolume 删除存储
func (s *ExtensionDockerService) DeleteVolume(name string, force bool) error {
	//初始化docker服务
	dockerClient, err := s.NewDockerService()
	if err != nil {
		return err
	}

	err = dockerClient.VolumeRemove(context.Background(), name, force)
	if err != nil {
		return err
	}
	return nil
}

// ContainerSSH 连接容器
func (s *ExtensionDockerService) ContainerSSH(c *gin.Context, param *request.ContainerSSHR) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.Log.Errorf("upGrader.Upgrade:gin context http handler failed, err: %v", err)
		return
	}
	defer func(wsConn *websocket.Conn) {
		_ = wsConn.Close()
	}(wsConn)

	cmdStr := fmt.Sprintf("docker exec -u %s %s %s", param.User, param.ContainerId, param.Command)
	_, err = util.ExecShell(cmdStr)
	if wsHandleError(wsConn, err) {
		return
	}

	commands := fmt.Sprintf("docker exec -it -u %s %s %s", param.User, param.ContainerId, param.Command)
	slave, err := terminal.NewCommand(commands)
	if wsHandleError(wsConn, err) {
		return
	}
	defer func(slave *terminal.LocalCommand) {
		_ = slave.Close()
	}(slave)

	tty, err := terminal.NewLocalWsSession(param.Cols, param.Rows, wsConn, slave)
	if wsHandleError(wsConn, err) {
		return
	}

	quitChan := make(chan bool, 3)
	tty.Start(quitChan)
	go slave.Wait(quitChan)

	<-quitChan

	global.Log.Info("websocket finished")
	if wsHandleError(wsConn, err) {
		return
	}

}

package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ProjectPHPService struct {
}

var (
	checkProjectExpirationWg         sync.WaitGroup
	checkProjectExpirationGisRunning atomic.Value
	quitCheckProjectExpirationG      chan string
)

func CheckProjectExpirationInit() {
	if checkProjectExpirationGisRunning.Load() == true {
		// 已在运行中,直接返回
		return
	}
	// 运行检测项目过期监控
	checkProjectExpirationGisRunning.Store(true)
	quitCheckProjectExpirationG = make(chan string)
	checkProjectExpirationWg.Add(1)
	go startCheckProjectExpiration(quitCheckProjectExpirationG)
	_, _ = fmt.Fprintf(color.Output, "check project Expiration....   %s\n",
		color.GreenString("done"),
	)
}
func startCheckProjectExpiration(quit <-chan string) {
	defer checkProjectExpirationWg.Done()
	for {
		select {
		case <-quit:
			_, _ = fmt.Fprintf(color.Output, "check project Expiration quit....   %s\n",
				color.GreenString("done"),
			)
			checkProjectExpirationGisRunning.Store(false)
			return //必须return，否则goroutine不会结束
		default:
			time.Sleep(60 * time.Second)
			//查询是否有过期的项目
			projects, _, err := (&model.Project{}).List(global.PanelDB, &model.ConditionsT{
				"status = ?":                           constant.ProjectStatusByRunning,
				"expire_time <> 0 AND expire_time < ?": time.Now().Unix(),
			}, 0, 0)
			if err != nil {
				global.Log.Errorf("startCheckProjectExpiration->Project.List  Error:%s", err)
				return
			}
			var projectList []string
			for _, project := range projects {
				projectList = append(projectList, project.Name)
				err = (&ProjectPHPService{}).SetStatus(project, "stop")
				if err != nil {
					global.Log.Errorf("startCheckProjectExpiration->ProjectPHPService.SetStatus  Error:%s", err)
					return
				}
			}
			err = (&MonitorService{}).ProjectExpirationTimeEvent(projectList)
			if err != nil {
				global.Log.Errorf("startCheckProjectExpiration->MonitorService.ProjectExpirationTimeEvent  Error:%s", err)
			}
		}
	}

}
func (s *ProjectPHPService) Create(param *request.CreatePHPProjectR) error {
	//检查默认配置
	CheckConfigPath()
	defer func() {
		err := (&ExtensionNginxService{}).SetStatus("reload")
		if err != nil {
			global.Log.Warnf("ProjectPHPService->Create->SetStatus,nginx reload failed:%s", err)
		}
	}()
	//检查nginx配置
	nginxService := ExtensionNginxService{}
	err := nginxService.CheckConfig()
	if err != nil {
		return err
	}

	//检查主域名是否合法
	err = CheckDomainItem(param.Domain)
	if err != nil {
		return err
	}
	global.Log.Debugf("Create->projectName:%s \n", param.Name)
	projectPath := util.TrimPath(param.Path)
	phpVersion := param.PHPVersion
	//检查项目名称是否重复
	get, err := (&model.Project{Name: param.Name}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if get.ID > 0 {
		return errors.New(helper.MessageWithMap("project.ProjectNameAlreadyExists", map[string]any{"Name": param.Name}))
	}

	//如果是选择了php版本检查对应的php版本是否存在
	if phpVersion != "00" {
		phpService := ExtensionPHPService{}
		phpInfos, err := phpService.Info()
		if err != nil {
			return err
		}
		phpIsExist := false
		for _, phpInfo := range phpInfos {
			if phpInfo.Description.Install && phpInfo.VersionList[0].MVersion == phpVersion {
				phpIsExist = true
				break
			}
		}
		if !phpIsExist {
			return errors.New(helper.Message("extension.NotInstalled"))
		}
	}

	//敏感路径校验
	if util.CheckSensitivePath(projectPath) {
		return errors.New(helper.MessageWithMap("explorer.CannotOperateSensitivePath", map[string]any{"Path": projectPath}))
	}

	//如果根目录不存在创建根目录
	if !util.PathExists(projectPath) {
		err = os.MkdirAll(projectPath, 0755)
		if err != nil {
			return err
		}
		_, _ = util.ExecShell(fmt.Sprintf("chown -R 755 %s && chown -R www:www %s", projectPath, projectPath))
	}
	//创建user.ini
	userIniPath := fmt.Sprintf("%s/.user.ini", projectPath)
	if !util.PathExists(userIniPath) {
		_ = util.WriteFile(userIniPath, []byte(fmt.Sprintf("open_basedir=%s/:/tmp/", projectPath)), 0755)
		_, _ = util.ExecShell(fmt.Sprintf("chown  644 %s  && chown  root:root %s  && chattr  +i %s", userIniPath, userIniPath, userIniPath))
	}

	//创建nginx basedir,可能是防跨站攻击设置，暂时不用

	//创建默认文档
	defaultIndexDocument := fmt.Sprintf("%s/index.html", projectPath)
	if !util.PathExists(defaultIndexDocument) {
		defaultIndexContentStr, _ := util.ReadFileStringBody(fmt.Sprintf("%s/template/index.html", GetExtensionsPath(constant.ExtensionNginxName)))
		_ = util.WriteFile(defaultIndexDocument, []byte(defaultIndexContentStr), 0755)
		_, _ = util.ExecShell(fmt.Sprintf("chown  755 %s  &&  chown  www:www %s", defaultIndexDocument, defaultIndexDocument))
	}

	//创建404页面
	default404Document := fmt.Sprintf("%s/404.html", projectPath)
	if !util.PathExists(default404Document) {
		notfoundContentStr, _ := util.ReadFileStringBody(fmt.Sprintf("%s/template/404.html", GetExtensionsPath(constant.ExtensionNginxName)))
		_ = util.WriteFile(default404Document, []byte(notfoundContentStr), 0755)
		_, _ = util.ExecShell(fmt.Sprintf("chown  755 %s && chown  www:www %s", default404Document, default404Document))
	}

	//写入nginx配置
	err = GenerateNginxConfig(&NginxConfig{
		ProjectName: param.Name,
		ProjectPath: projectPath,
		PHPVersion:  phpVersion,
		Domain:      []*request.DomainItem{param.Domain},
		SSL:         false,
	})
	if err != nil {
		return err
	}

	//生成写入ttwaf项目配置
	err = GroupApp.TTWafServiceApp.GenerateProjectConfig(param.Name)
	if err != nil {
		return err
	}

	//写入项目信息至数据库
	var createData model.Project
	createData.Name = param.Name
	createData.Path = param.Path
	createData.CategoryId = 0
	createData.Ps = param.Ps
	createData.Status = constant.ProjectStatusByRunning
	createData.ExpireTime = 0
	createData.ProjectType = constant.ProjectTypeByPHP
	createProject, err := (&createData).Create(global.PanelDB)
	if err != nil || createProject == nil {
		return err
	}

	err = InsertDomain(param.Domain, createProject.ID)
	if err != nil {
		return err
	}

	//更新ttwaf域名配置文件
	err = GroupApp.TTWafServiceApp.OperateDomainConfig(createProject.ID, createProject.Name, createProject.Path, constant.OperateTTWafDomainConfigByUpdate)
	if err != nil {
		return err
	}

	//如果有数据库信息则添加
	if param.DataBase != nil {
		info := request.CreateMysqlR{
			DatabaseName:     param.DataBase.DBName,
			UserName:         param.DataBase.User,
			Password:         param.DataBase.Password,
			Coding:           param.DataBase.Coding,
			AccessPermission: "127.0.0.1", //%表示允许所有人访问，本地填 127.0.0.1,特定ip:多个ip使用 , 分割
			Ps:               param.Name,
			Sid:              0,
			Pid:              0,
		}
		mysqlService := DatabaseMysqlService{}
		err = mysqlService.Create(&info)
		if err != nil {
			return err
		}
	}
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
				return
			}
		}()
	}

	//添加更多域名
	if len(param.MoreDomains) > 0 {
		addDomainsParam := &request.AddDomainsR{
			ProjectID: createProject.ID,
			Domains:   param.MoreDomains,
		}
		_ = GroupApp.ProjectServiceApp.AddDomains(createProject, addDomainsParam)
	}

	return nil

}

// ProjectList php项目列表
func (s *ProjectPHPService) ProjectList(param *request.PHPProjectListR, offset, limit int) (list []*response.PHPProject, total int64, err error) {
	condition := model.ConditionsT{
		"project_type = ?": constant.ProjectTypeByPHP,
		"ORDER":            "project.create_time DESC",
	}
	if !util.StrIsEmpty(param.Query) {
		searchPidS, err := (&model.Project{}).Search(global.PanelDB, &condition, param.Query)
		if err != nil {
			return nil, 0, err
		}
		if len(searchPidS) > 0 {
			condition["id IN ?"] = searchPidS
		} else {
			return list, total, nil
		}
	}
	projects, total, err := (&model.Project{}).List(global.PanelDB, &condition, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	//读取全局配置文件
	ttWafProjectConfigs, err := GroupApp.TTWafServiceApp.GetProjectConfig()
	for _, projectInfo := range projects {
		//获取php版本
		phpVersion, _ := GetProjectPHPVersion(projectInfo.Name)
		//获取ssl信息
		sslInfo, _ := GetProjectSSLInfo(projectInfo.Name)
		list = append(list, &response.PHPProject{
			Project:     projectInfo,
			PHPVersion:  phpVersion,
			TTWafStatus: ttWafProjectConfigs[projectInfo.Name].Status,
			SSL:         sslInfo,
		})
	}

	if err != nil {
		return nil, 0, err
	}
	return list, total, nil

}

// ProjectInfo php项目详情
func (s *ProjectPHPService) ProjectInfo(projectID int64) (*response.PHPProject, error) {
	projectGet, err := (&model.Project{
		ID:          projectID,
		ProjectType: constant.ProjectTypeByPHP,
	}).Get(global.PanelDB)
	if err != nil {
		return nil, err
	}
	if projectGet.ID == 0 {
		return nil, errors.New("not found project or project type is not php")
	}

	//读取全局配置文件
	ttWafProjectConfigs, err := GroupApp.TTWafServiceApp.GetProjectConfig()
	//获取php版本
	phpVersion, _ := GetProjectPHPVersion(projectGet.Name)
	//获取ssl信息
	sslInfo, err := GetProjectSSLInfo(projectGet.Name)
	if err != nil {
		return nil, err
	}
	//获取运行目录
	var runPath string
	rootPath, _ := GetProjectRootPath(projectGet.Name)
	if rootPath == projectGet.Path {
		runPath = "/"
	} else {
		runPath = strings.Replace(projectGet.Path, rootPath, "", 1)
	}
	//获取root目录下的文件夹列表
	runPathDirList := []string{"/"}
	err = filepath.Walk(projectGet.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path != projectGet.Path && info.IsDir() {
			runPathDirList = append(runPathDirList, "/"+info.Name())
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	//判断user.ini是否存在
	userIni := false
	if util.PathExists(rootPath + "/user.ini") {
		userIni = true
	}
	projectInfo := &response.PHPProject{
		Project:        projectGet,
		PHPVersion:     phpVersion,
		TTWafStatus:    ttWafProjectConfigs[projectGet.Name].Status,
		SSL:            sslInfo,
		RunPath:        runPath,
		RunPathDirList: runPathDirList,
		UserIni:        userIni,
	}
	if err != nil {
		return nil, err
	}
	return projectInfo, nil
}

// Delete 删除项目
func (s *ProjectPHPService) Delete(projectData *model.Project, param *request.DeletePHPProjectR) (err error) {

	//删除nginx配置
	DeleteNginxConfig(projectData.Name)

	//删除ttwaf配置
	projectConfigs, err := GroupApp.TTWafServiceApp.GetProjectConfig()
	if err != nil {
		return
	}
	delete(projectConfigs, projectData.Name)

	//写入项目配置文件
	err = GroupApp.TTWafServiceApp.WriteProjectConfig(projectConfigs)
	if err != nil {
		return
	}

	//删除域名
	err = (&model.ProjectDomain{}).Delete(global.PanelDB, &model.ConditionsT{"project_id": projectData.ID})
	if err != nil {
		return
	}

	//删除ttwaf域名配置信息
	err = GroupApp.TTWafServiceApp.OperateDomainConfig(projectData.ID, projectData.Name, projectData.Path, constant.OperateTTWafDomainConfigByDelete)
	if err != nil {
		return
	}

	//根据条件删除其余项
	if param.ClearNginxLog {
		_ = os.Remove(fmt.Sprintf("%s/%s.log", global.Config.System.WwwLogPath, projectData.Name))
		_ = os.Remove(fmt.Sprintf("%s/%s.error.log", global.Config.System.WwwLogPath, projectData.Name))
	}
	if param.ClearPath {
		//检查敏感路径
		if !util.CheckSensitivePath(projectData.Path) {
			_ = os.RemoveAll(projectData.Path)
		}
	}

	//删除项目
	err = projectData.Delete(global.PanelDB, &model.ConditionsT{})
	if err != nil {
		return
	}
	return
}

// SetStatus 设置项目状态
func (s *ProjectPHPService) SetStatus(projectData *model.Project, action string) (err error) {
	if action == "start" {
		//检查到期时间
		if projectData.ExpireTime != 0 && projectData.ExpireTime < time.Now().Unix() {
			return errors.New("project has expired")
		}
		//设置项目运行路径为
		err = s.SetRunPath(projectData.Name, projectData.Path, "")
		if err != nil {
			return
		}
		projectData.Status = constant.ProjectStatusByRunning
	} else {
		//设置项目运行路径为
		err = s.SetRunPath(projectData.Name, global.Config.System.PanelPath+"/data/extensions/nginx/stop", "/")
		if err != nil {
			return
		}
		projectData.Status = constant.ProjectStatusByStop
	}
	return projectData.Update(global.PanelDB)
}

// SetExpireTime 设置项目到期时间
func (s *ProjectPHPService) SetExpireTime(projectData *model.Project, expireTime int64) (err error) {
	//判断到期时间是否小于当前时间
	if expireTime < time.Now().Unix() {
		return errors.New("expire time is less than current time")
	}
	projectData.ExpireTime = expireTime
	return projectData.Update(global.PanelDB)
}

// SetRunPath 设置运行目录
func (s *ProjectPHPService) SetRunPath(projectName string, projectPath string, runPath string) (err error) {
	if runPath == "/" {
		runPath = ""
	}
	//检查路径是否存在
	if !util.PathExists(projectPath + runPath) {
		return errors.New("path not exists")
	}
	return SetRootPath(projectName, projectPath+runPath)
}

// SetUserIni 设置user.ini
func (s *ProjectPHPService) SetUserIni(projectName string, status bool) (err error) {
	//获取运行目录
	rootPath, err := GetProjectRootPath(projectName)
	if err != nil {
		return err
	}
	if status {
		config := fmt.Sprintf("open_basedir=%s/:/tmp/", rootPath)
		return util.WriteFile(rootPath+"/user.ini", []byte(config), 0644)
	} else {
		_ = os.Remove(rootPath + "/user.ini")
		return
	}
}

// UsingPHPVersion	使用的PHP版本
func (s *ProjectPHPService) UsingPHPVersion(projectName string) (string, error) {
	confPath := ProjectMainConfFilePath(projectName)
	confBody, err := util.ReadFileStringBody(confPath)
	if err != nil {
		return "", err
	}

	//检查是否是自定义php配置
	customizePHPConfPath := ProjectConfDirPath(projectName) + "/enable-php-customize.conf"
	if strings.Contains(confBody, customizePHPConfPath) {
		//自定义配置
		body, err := util.ReadFileStringBody(customizePHPConfPath)
		if err != nil {
			return "", err
		}
		return body, nil
	}

	//匹配include enable-php-xxx.conf;
	rep := regexp.MustCompile(`\s+include\s+enable-php-(\w{2,5})\.conf;`)
	if rep.MatchString(confBody) {
		return rep.FindStringSubmatch(confBody)[1], nil
	}
	return "", errors.New("the PHP configuration item was not found, the project configuration file may have been manually modified, if not, please report this error")
}

// SwitchUsingPHPVersion 切换PHP版本
func (s *ProjectPHPService) SwitchUsingPHPVersion(projectName string, version, customize string) error {

	replaceValue := ""
	customizePath := ProjectConfDirPath(projectName) + "/enable-php-customize.conf"
	if version == "customize" {
		if !regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+:\d+|unix:[\w/\.-]+)$`).MatchString(customize) {
			return errors.New("the customized configuration format is incorrect")
		}
		customizeTmp := strings.Split(customize, ":")
		if customizeTmp[0] == "unix" {
			if !util.PathExists(customizeTmp[1]) {
				return errors.New("the customized configuration unix file does not exist")
			}
		} else {
			if _, err := util.CheckConnection("tcp", customizeTmp[0], customizeTmp[1], 3); err != nil {
				return errors.New("the customized configuration tcp port is not available")
			}
		}
		templateBody, err := util.ReadFileStringBody(GetExtensionsPath(constant.ExtensionNginxName) + "/template/enable_php/enable-php-customize.conf")
		if err != nil {
			return err
		}
		templateBody = strings.Replace(templateBody, "{{config}}", customize, -1)
		err = util.WriteFile(customizePath, []byte(templateBody), 0644)
		if err != nil {
			return err
		}
		replaceValue = customizePath
	} else {
		replaceValue = "enable-php-" + version + ".conf"
	}
	confPath := ProjectMainConfFilePath(projectName)
	confBody, err := util.ReadFileStringBody(confPath)
	if err != nil {
		return err
	}
	confBody = regexp.MustCompile(`\s+include\s+(.*enable-php-.*\.conf);`).ReplaceAllString(confBody, replaceValue)
	err = util.WriteFile(confPath, []byte(confBody), 0644)
	if err != nil {
		return err
	}
	return nil
}

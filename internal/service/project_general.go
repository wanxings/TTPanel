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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ProjectGeneralService struct {
}

func (s *ProjectGeneralService) GetProjectConfig(config string) (projectConfig *model.GeneralProjectConfig, err error) {
	projectConfig = &model.GeneralProjectConfig{}
	err = util.JsonStrToStruct(config, projectConfig)
	if err != nil {
		return
	}
	return
}

func (s *ProjectGeneralService) GetPidFilePath(projectName string) (filePath string) {
	filePath = fmt.Sprintf("%s/project/general/%s.pid", global.Config.Logger.RootPath, projectName)
	return
}
func (s *ProjectGeneralService) GetRunLogFilePath(projectName string) (filePath string) {
	filePath = fmt.Sprintf("%s/project/general/%s.log", global.Config.Logger.RootPath, projectName)
	return
}

func (s *ProjectGeneralService) Create(param *request.CreateGeneralProjectR) (err error) {
	if param.Name, err = util.IsValidProjectName(param.Name); err != nil {
		return
	}
	//检查项目名称是否存在
	get, err := (&model.Project{Name: param.Name}).Get(global.PanelDB)
	if err != nil {
		return
	}
	if get.ID > 0 {
		err = errors.New(helper.MessageWithMap("project.ProjectNameAlreadyExists", map[string]any{"Name": param.Name}))
		return
	}

	//检查项目路径是否存在
	param.Path = util.TrimPath(param.Path)
	if !util.PathExists(param.Path) {
		if param.AutoCreatePath {
			_ = os.MkdirAll(param.Path, 0755)
		} else {
			err = errors.New(fmt.Sprintf("path[%v] does not exist", param.Path))
			return
		}

	}

	//检查项目真实端口是否合法和是否被占用
	if param.Port < 10 || param.Port > 65535 {
		err = errors.New(helper.MessageWithMap("PortIsIllegalOrCommon", map[string]any{"Port": param.Port}))
		return
	}
	if util.CheckPortOccupied("tcp", param.Port) {
		err = errors.New(helper.MessageWithMap("PortOccupied", map[string]any{"Port": param.Port}))
		return
	}

	//准备project_config数据
	generalProjectConfig := model.GeneralProjectConfig{
		Command:   param.Command,
		Port:      param.Port,
		RunUser:   param.RunUser,
		IsPowerOn: param.IsPowerOn,
	}

	//准备project数据
	generalProject := model.Project{
		Name:        param.Name,
		Path:        param.Path,
		Status:      1,
		Ps:          param.Description,
		CategoryId:  0,
		ProjectType: constant.ProjectTypeByGeneral,
	}

	generalProject.ProjectConfig, err = util.StructToJsonStr(generalProjectConfig)
	if err != nil {
		return
	}

	//插入项目数据
	generalProjectCreate, err := generalProject.Create(global.PanelDB)
	if err != nil {
		return
	}

	//启动项目
	global.Log.Debugln("Start Project...")
	err = s.Start(generalProjectCreate.Name, generalProjectCreate.Path, generalProjectConfig.RunUser, generalProjectConfig.Command)
	if err != nil {
		return
	}
	return
}

// List 项目列表
func (s *ProjectGeneralService) List(param *request.GeneralProjectListR, offset, limit int) (projectList []*model.Project, total int64, err error) {
	condition := model.ConditionsT{
		"project_type = ?": constant.ProjectTypeByGeneral,
		"ORDER":            "create_time DESC",
	}
	if !util.StrIsEmpty(param.Query) {
		searchIDs, errs := (&model.Project{}).Search(global.PanelDB, &condition, param.Query)
		if errs != nil {
			err = errs
			return
		}
		if len(searchIDs) > 0 {
			condition["id IN ?"] = searchIDs
		}
	}
	projectList, total, err = (&model.Project{}).List(global.PanelDB, &condition, offset, limit)
	if err != nil {
		return
	}

	return
}

// Delete 删除项目
func (s *ProjectGeneralService) Delete(projectData *model.Project, param *request.DeleteGeneralProjectR) (err error) {
	if param.ClearRunLog {
		_ = os.Remove(s.GetRunLogFilePath(projectData.Name))
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
func (s *ProjectGeneralService) SetStatus(projectData *model.Project, action string) (err error) {
	//获取项目配置
	generalProjectConfig, err := s.GetProjectConfig(projectData.ProjectConfig)
	if err != nil {
		return
	}

	switch action {
	case constant.ProcessCommandByStart:
		if err = s.Start(projectData.Name, projectData.Path, generalProjectConfig.RunUser, generalProjectConfig.Command); err != nil {
			return
		}

	case constant.ProcessCommandByStop:
		if err = s.Stop(projectData.Name); err != nil {
			return
		}
	case constant.ProcessCommandByRestart:
		_ = s.Stop(projectData.Name)
		time.Sleep(time.Second * 4)
		if err = s.Start(projectData.Name, projectData.Path, generalProjectConfig.RunUser, generalProjectConfig.Command); err != nil {
			return
		}
	}
	return
}

// Stop 停止项目
func (s *ProjectGeneralService) Stop(projectName string) error {
	pidFilePath := s.GetPidFilePath(projectName)
	if !util.PathExists(pidFilePath) {
		//如果pid文件不存在表示项目未启动
		return nil
	}
	//读取pid文件
	pid, err := util.ReadFileStringBody(pidFilePath)
	if err != nil {
		return err
	}
	pid = strings.TrimSpace(pid)
	pidInt, _ := strconv.Atoi(pid)
	pidInt32 := int32(pidInt)

	_, err = (&TaskManagerService{}).KillProcess(pidInt32, false)
	if err != nil {
		return err
	}

	//删除pid文件
	_ = os.Remove(pidFilePath)
	return nil
}

// Start 启动项目
func (s *ProjectGeneralService) Start(projectName string, projectPath string, runUser string, command string) error {
	runLogFilePath := s.GetRunLogFilePath(projectName)
	pidFilePath := s.GetPidFilePath(projectName)
	err := os.MkdirAll(filepath.Dir(runLogFilePath), 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(pidFilePath), 0755)
	if err != nil {
		return err
	}
	if !util.PathExists(pidFilePath) {
		//如果pid文件不存在则创建
		_, _ = util.ExecShell("touch " + pidFilePath)
	} else {
		//如果pid文件存在设置权限
		_, _ = util.ExecShell(fmt.Sprintf("chown %s %s; chmod 755 %s", projectName, pidFilePath, pidFilePath))
	}
	if !util.PathExists(runLogFilePath) {
		//如果log文件不存在则创建
		_, _ = util.ExecShell("touch " + runLogFilePath)
	} else {
		//如果log文件存在设置权限
		_, _ = util.ExecShell(fmt.Sprintf("chown %s %s; chmod 755 %s",
			runUser,
			runLogFilePath,
			runLogFilePath,
		))
	}
	cmdStr := fmt.Sprintf("cd %s;nohup %s >> %s 2>&1 & pid=$!; echo $pid > %s",
		projectPath,
		command,
		runLogFilePath,
		pidFilePath,
	)

	result, err := util.ExecShellAsUser(cmdStr, runUser)
	if err != nil {
		global.Log.Debugf("Start->ExecShellAsUser Error:%v,%v\n", err, result)
		return err
	}
	////休眠两秒
	//time.Sleep(2 * time.Second)
	////检查是否启动成功
	//processList, _ := (&safe.TaskManagerService{}).ProcessList()
	//mainProcess, _ := s.GetProjectProcess(projectName, processList)
	//if mainProcess == nil {
	//	return errors.New(messages.ProjectStartFailed.MsgF([]interface{}{projectName}))
	//}
	return nil
}

// GetProjectProcess 获取项目进程列表
func (s *ProjectGeneralService) GetProjectProcess(projectName string, processList []response.ProcessInfoP) (main *response.ProcessInfoP, child []*response.ProcessInfoP) {
	pidFilePath := s.GetPidFilePath(projectName)
	if !util.PathExists(pidFilePath) {
		return
	}
	pid, err := util.ReadFileStringBody(pidFilePath)
	if err != nil {
		return
	}
	pid = strings.TrimSpace(pid)
	pidInt, _ := strconv.Atoi(pid)
	pidInt32 := int32(pidInt)
	if pidInt32 == 0 {
		return
	}

	for _, v := range processList {
		if v.Pid == pidInt32 {
			global.Log.Debugln("ProjectMainProcessP:", v)
			i := v
			main = &i
		} else if v.PPid == pidInt32 {
			i := v
			child = append(child, &i)
		}
	}
	return
}

// GetDetails 项目信息
func (s *ProjectGeneralService) GetDetails(projectBaseInfo *model.Project, processList []response.ProcessInfoP) (projectInfo *response.GeneralProject, err error) {

	projectInfo = &response.GeneralProject{
		Project:        projectBaseInfo,
		RunLogFilePath: s.GetRunLogFilePath(projectBaseInfo.Name),
	}
	projectInfo.MainProcess, projectInfo.ChildProcess = s.GetProjectProcess(projectBaseInfo.Name, processList)
	projectInfo.ProjectConfig, err = s.GetProjectConfig(projectBaseInfo.ProjectConfig)
	return
}

// SaveProjectConfig 保存项目配置
func (s *ProjectGeneralService) SaveProjectConfig(projectData *model.Project, param *request.SaveProjectConfigR) (err error) {
	//获取项目配置
	projectConfig, err := s.GetProjectConfig(projectData.ProjectConfig)

	//检查项目路径是否存在
	param.Path = util.TrimPath(param.Path)
	if !util.PathExists(param.Path) {
		if param.AutoCreatePath {
			_ = os.MkdirAll(param.Path, 0755)
		} else {
			return errors.New(fmt.Sprintf("path[%v] does not exist", param.Path))
		}

	}

	if projectConfig.Port != param.Port {
		//如果修改了端口检查端口是否合法和是否被占用
		if param.Port < 10 || param.Port > 65535 {
			return errors.New(helper.MessageWithMap("PortIsIllegalOrCommon", map[string]any{"Port": param.Port}))
		}
		if util.CheckPortOccupied("tcp", param.Port) {
			return errors.New(helper.MessageWithMap("PortOccupied", map[string]any{"Port": param.Port}))
		}
	}

	//覆盖项目配置
	projectConfig.Port = param.Port
	projectConfig.Command = param.Command
	projectConfig.RunUser = param.RunUser
	projectConfig.IsPowerOn = param.IsPowerOn
	projectData.ProjectConfig, err = util.StructToJsonStr(projectConfig)
	if err != nil {
		return
	}
	err = projectData.Update(global.PanelDB)
	if err != nil {
		return
	}

	//重启项目
	err = s.SetStatus(projectData, constant.ProcessCommandByRestart)
	if err != nil {
		return err
	}
	return nil
}

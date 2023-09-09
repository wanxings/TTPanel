package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type CronTaskService struct{}

func (s *CronTaskService) BatchCreate(taskInfoList []*request.CronTaskInfo) []error {
	var errs []error
	for _, taskInfo := range taskInfoList {
		err := s.Create(taskInfo)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return nil
}

// Create 创建计划任务
func (s *CronTaskService) Create(taskInfo *request.CronTaskInfo) error {
	hash := util.EncodeMD5(strconv.FormatInt(time.Now().Unix(), 10))
	err := WriteCronFile(hash, taskInfo)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			//删除cron文件中的计划任务
			RemoveForCron(hash)
			//删除taskShell文件
			RemoveForShell(hash)
		}
	}()

	taskInfo.Name = XssFilter(taskInfo.Name)
	configI, err := s.GetCronTaskConfig(taskInfo.Category, &taskInfo.CronTaskConfig)
	if err != nil {
		return err
	}
	configJsonStr := "{}"
	if configI != nil {
		configJsonStr, err = util.StructToJsonStr(configI)
		if err != nil {
			return err
		}
	}
	// 将计划任务插入数据库
	insertData := &model.CronTask{
		Name:          taskInfo.Name,
		Category:      taskInfo.Category,
		TimeType:      taskInfo.TimeType,
		TimeInc:       taskInfo.TimeInc,
		TimeHour:      taskInfo.TimeHour,
		TimeMinute:    taskInfo.TimeMinute,
		TimeCustomize: taskInfo.TimeCustomize,
		ShellBody:     taskInfo.ShellBody,
		Config:        configJsonStr,
		Hash:          hash,
		Status:        1,
	}
	global.Log.Debugf("CronTask-insertData:%+v", insertData)
	global.Log.Debugf("CronTask-taskInfo:%+v", taskInfo)
	_, err = (insertData).Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// GetCronTaskConfig 获取计划任务配置接口
func (s *CronTaskService) GetCronTaskConfig(category int, cronTaskConfig *request.CronTaskConfig) (config interface{}, err error) {
	switch category {
	case constant.CronTaskCategoryByShell: //执行shell
		return nil, nil
	case constant.CronTaskCategoryByBackupProject: //备份项目
		return struct {
			BackupProjectConfig model.BackupProjectConfig `json:"backup_project_config"`
		}{
			BackupProjectConfig: cronTaskConfig.BackupProjectConfig,
		}, err
	case constant.CronTaskCategoryByBackupDatabase: //备份数据库
		return struct {
			BackupDatabaseConfig model.BackupDatabaseConfig `json:"backup_database_config"`
		}{
			BackupDatabaseConfig: cronTaskConfig.BackupDatabaseConfig,
		}, err
	case constant.CronTaskCategoryByCutLog: //切割日志
		return struct {
			CutLogConfig model.CutLogConfig `json:"cut_log_config"`
		}{
			CutLogConfig: cronTaskConfig.CutLogConfig,
		}, err
	case constant.CronTaskCategoryByBackupDir: //备份目录
		return struct {
			BackupDirConfig model.BackupDirConfig `json:"backup_dir_config"`
		}{
			BackupDirConfig: cronTaskConfig.BackupDirConfig,
		}, err
	case constant.CronTaskCategoryByRequestUrl: //请求url
		return struct {
			RequestUrlConfig model.RequestUrlConfig `json:"request_url_config"`
		}{
			RequestUrlConfig: cronTaskConfig.RequestUrlConfig,
		}, err
	case constant.CronTaskCategoryByFreeMemory: //释放内存
		return nil, err
	default:
		return nil, errors.New("CronTaskCategory is invalid")
	}
}

// Edit 编辑任务
func (s *CronTaskService) Edit(taskEditData *request.EditR) error {
	//查询数据库中该计划任务的详细信息
	taskInfo, err := (&model.CronTask{ID: taskEditData.Id}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if taskInfo.ID == 0 {
		return errors.New("cronTask does not exist")
	}
	configI, err := s.GetCronTaskConfig(taskInfo.Category, &taskEditData.CronTaskConfig)
	if err != nil {
		return err
	}
	configJsonStr := "{}"
	if configI != nil {
		configJsonStr, err = util.StructToJsonStr(configI)
		if err != nil {
			return err
		}
	}
	taskInfo.Name = taskEditData.Name
	taskInfo.TimeType = taskEditData.TimeType
	taskInfo.TimeInc = taskEditData.TimeInc
	taskInfo.TimeHour = taskEditData.TimeHour
	taskInfo.TimeMinute = taskEditData.TimeMinute
	taskInfo.TimeCustomize = taskEditData.TimeCustomize
	taskInfo.ShellBody = taskEditData.ShellBody
	taskInfo.Config = configJsonStr
	//更新数据
	err = (taskInfo).Update(global.PanelDB)
	if err != nil {
		return err
	}
	//删除cron文件中旧的计划任务
	RemoveForCron(taskInfo.Hash)
	//将新任务写入cron文件中
	if taskInfo.Status != 0 {
		err = WriteCronFile(taskInfo.Hash, taskEditData.CronTaskInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// BatchSetStatus 批量设置任务状态
func (s *CronTaskService) BatchSetStatus(ids []int64, status int) ([]string, error) {
	taskNameList := make([]string, 0)
	for _, id := range ids {
		taskInfo, err := (&model.CronTask{ID: id}).Get(global.PanelDB)
		if err != nil {
			return nil, err
		}
		if taskInfo.ID == 0 {
			return nil, errors.New("cronTask does not exist")
		}
		switch status {
		case 0:
			//如果是关闭任务，判断任务原本的状态是否是关闭状态，如果是开启状态，那么就删除cron文件中的计划任务，否则不做任何操作
			if taskInfo.Status == 1 {
				RemoveForCron(taskInfo.Hash)
			}
		case 1:
			//如果是开启任务，判断任务原本的状态是否是关闭状态，如果是关闭状态，那么就重新写入cron文件，否则不做任何操作
			if taskInfo.Status == 0 {
				//如果原本是关闭状态，那么就重新写入cron文件
				configS := request.CronTaskConfig{}
				err = util.JsonStrToStruct(taskInfo.Config, &configS)
				if err != nil {
					return nil, err
				}
				err = WriteCronFile(taskInfo.Hash, &request.CronTaskInfo{
					Name:           taskInfo.Name,
					Category:       taskInfo.Category,
					TimeType:       taskInfo.TimeType,
					ShellBody:      taskInfo.ShellBody,
					TimeInc:        taskInfo.TimeInc,
					TimeHour:       taskInfo.TimeHour,
					TimeMinute:     taskInfo.TimeMinute,
					TimeCustomize:  taskInfo.TimeCustomize,
					CronTaskConfig: configS,
				})
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, errors.New("status is empty")
		}

		//更新数据库中的任务状态
		taskInfo.Status = status
		err = (taskInfo).Update(global.PanelDB)
		if err != nil {
			return nil, err
		}
		taskNameList = append(taskNameList, taskInfo.Name)
	}
	return taskNameList, nil
}

// BatchDelete 批量删除任务
func (s *CronTaskService) BatchDelete(ids []int64) ([]string, error) {
	taskNameList := make([]string, 0)
	for _, id := range ids {
		taskInfo, err := (&model.CronTask{ID: id}).Get(global.PanelDB)
		if err != nil {
			return nil, err
		}
		if taskInfo.ID == 0 {
			return nil, errors.New("cronTask does not exist")
		}
		//删除cron文件中的计划任务
		RemoveForCron(taskInfo.Hash)
		//删除taskShell文件
		RemoveForShell(taskInfo.Hash)
		//删除数据库中的计划任务
		err = (taskInfo).Delete(global.PanelDB, &model.ConditionsT{})
		if err != nil {
			return nil, err
		}
		taskNameList = append(taskNameList, taskInfo.Name)
	}
	return taskNameList, nil
}

// BatchRun 批量执行任务
func (s *CronTaskService) BatchRun(ids []int64) ([]string, error) {
	taskNameList := make([]string, 0)
	for _, id := range ids {
		taskInfo, err := (&model.CronTask{ID: id}).Get(global.PanelDB)
		if err != nil {
			return nil, err
		}
		if taskInfo.ID == 0 {
			return nil, errors.New("cronTask does not exist")
		}
		//执行任务
		err = RunTask(taskInfo.Hash)
		if err != nil {
			return nil, err
		}
		taskNameList = append(taskNameList, taskInfo.Name)
	}
	return taskNameList, nil
}

// GetExecutionLog 获取执行日志
func (s *CronTaskService) GetExecutionLog(id int64, line int) (string, error) {
	//查询任务详情
	taskInfo, err := (&model.CronTask{ID: id}).Get(global.PanelDB)
	if err != nil {
		return "", err
	}
	if taskInfo.ID == 0 {
		return "", errors.New("cronTask does not exist")
	}
	//获取执行日志
	shellPath := global.Config.System.ServerPath + "/cron/" + taskInfo.Hash + ".log"
	log, err := TailFile(shellPath, line)
	if err != nil {
		return "", err
	}
	return log, nil
}

// TailFile 取文件指定尾行数
func TailFile(path string, line int) (string, error) {
	body, err := util.ExecShell(fmt.Sprintf("tail -n %d %s", line, path))
	if err != nil {
		return "", err
	}
	return body, nil
}

// RunTask 执行任务
func RunTask(echo string) error {
	var err error
	//执行任务
	shellPath := global.Config.System.ServerPath + "/cron/" + echo
	_, err = util.ExecShell("chmod +x " + shellPath)
	if err != nil {
		return err
	}
	_, err = util.ExecShell("nohup " + shellPath + " >> " + shellPath + ".log 2>&1 &")
	if err != nil {
		return err
	}
	return nil
}

// List 获取任务列表
func (s *CronTaskService) List(param *request.TaskListR) ([]*model.CronTask, int64, error) {
	//构造查询条件
	conditions := &model.ConditionsT{"ORDER": "create_time DESC"}
	if !util.StrIsEmpty(param.QueryName) {
		(*conditions)["name"] = param.QueryName
	}
	if !util.StrIsEmpty(param.SType) {
		(*conditions)["s_type"] = param.SType
	}
	taskList, total, err := (&model.CronTask{}).List(global.PanelDB, conditions, 0, 0)
	if err != nil {
		return nil, 0, err
	}

	for _, task := range taskList {
		//取上次执行的时间
		//判断日志文件是否存在
		cronPath := fmt.Sprintf("%s/cron/%s.log", global.Config.System.ServerPath, task.Hash)
		if util.PathExists(cronPath) {
			// 获取文件的信息
			fileInfo, err := os.Stat(cronPath)
			if err != nil {
				global.Log.Debugf("CronTaskService.List.Stat CronPath:%s Error:%s", cronPath, err.Error())
			}
			// 获取修改时间
			task.LastRunTime = fileInfo.ModTime().Unix()
		} else {
			task.LastRunTime = task.CreateTime
		}
	}

	return taskList, total, nil
}

// Get 获取任务详情
func (s *CronTaskService) Get(id int64) (*model.CronTask, error) {
	taskInfo, err := (&model.CronTask{ID: id}).Get(global.PanelDB)
	if err != nil {
		return nil, err
	}
	if taskInfo.ID == 0 {
		return nil, errors.New("cronTask does not exist")
	}
	return taskInfo, nil
}

// WriteCronFile 写到cron文件
func WriteCronFile(cronHash string, taskInfo *request.CronTaskInfo) error {
	cronConfig, err := GetCronCycle(&request.CronCycle{
		Type:   taskInfo.TimeType,
		Minute: taskInfo.TimeMinute,
		Hour:   taskInfo.TimeHour,
		Inc:    taskInfo.TimeInc,
	})
	cronPath := global.Config.System.ServerPath + "/cron"
	err = WriteTaskShell(cronHash, taskInfo)
	if err != nil {
		return err
	}
	cronConfig += " " + cronPath + "/" + cronHash + " >> " + cronPath + "/" + cronHash + ".log 2>&1"
	err = AddShellToCron(cronConfig)
	if err != nil {
		return err
	}
	return nil
}

// RemoveForCron 删除cron
func RemoveForCron(hash string) bool {
	cronFilePath := GetCronRootPath()
	cronBody, err := util.ReadFileStringBody(cronFilePath)
	if err != nil {
		return false
	}
	reg := regexp.MustCompile(".+" + hash + ".+\n")
	cronBody = reg.ReplaceAllString(cronBody, "")
	err = util.WriteFile(cronFilePath, []byte(cronBody), 0644)
	if err != nil {
		return false
	}
	//修改文件夹权限
	_ = GroupApp.ExplorerServiceApp.BatchChangePermission([]string{cronFilePath}, "root", "600")
	ReloadCron()
	return true
}

// RemoveForShell 删除shell
func RemoveForShell(echo string) bool {
	shellPath := global.Config.System.ServerPath + "/cron/" + echo
	err := os.Remove(shellPath)
	if err != nil {
		return false
	}
	return true
}

// GetCronCycle 构造周期
func GetCronCycle(cycle *request.CronCycle) (cronTime string, err error) {
	var cronConfig string
	switch cycle.Type {
	case "day":
		cronConfig = GetDay(cycle)
	case "day_n":
		cronConfig = GetDayN(cycle)
	case "hour":
		cronConfig = GetHour(cycle)
	case "hour_n":
		cronConfig = GetHourN(cycle)
	case "minute_n":
		cronConfig = GetMinuteN(cycle)
	case "week":
		cronConfig = GetWeek(cycle)
	case "month":
		cronConfig = GetMonth(cycle)
	case "customize":
		cronConfig = cycle.Customize
	default:
		return "", errors.New("cycle.Type is empty")
	}
	return cronConfig, nil
}

// WriteTaskShell 编写任务脚本返回shell任务文件名
func WriteTaskShell(cronHash string, param *request.CronTaskInfo) error {
	var shell string
	shell = "#!/bin/bash\nPATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin\nexport PATH\n"
	//log := "-access_log"
	shell = shell + "\n" + "echo \"----------------------------------------------------------------------------\"\n" + "startDate=`date +\"%Y-%m-%d %H:%M:%S\"`\n" + "echo \"★[$startDate] Start\"\n" + "echo \"----------------------------------------------------------------------------\"\n"
	cronTaskShell := fmt.Sprintf("%s cd %s\n./TTPanel tools cronTask -i %s\n", shell, global.Config.System.PanelPath, cronHash)
	switch param.Category {
	case constant.CronTaskCategoryByShell: //执行shell
		shell += param.ShellBody
	case constant.CronTaskCategoryByBackupProject: //备份项目
		shell += cronTaskShell
	case constant.CronTaskCategoryByBackupDatabase: //备份数据库
		shell += cronTaskShell
	case constant.CronTaskCategoryByCutLog: //切割日志
		shell += cronTaskShell
	case constant.CronTaskCategoryByBackupDir: //备份目录
		shell += cronTaskShell
	case constant.CronTaskCategoryByRequestUrl: //请求url
		shell += cronTaskShell
	case constant.CronTaskCategoryByFreeMemory: //释放内存
		shell += fmt.Sprintf("%s cd %s/data/shell\nbash free_memory.sh\n", shell, global.Config.System.PanelPath)
	default:
		return errors.New("CronTask.Category is empty")
	}

	shell = shell + "\n" + "echo \"----------------------------------------------------------------------------\"\n" + "endDate=`date +\"%Y-%m-%d %H:%M:%S\"`\n" + "echo \"★[$endDate] Successful\"\n" + "echo \"----------------------------------------------------------------------------\"\n"
	cronPath := global.Config.System.ServerPath + "/cron"
	//判断cronPath是否存在
	if !util.PathExists(cronPath) {
		_, _ = util.ExecShell("mkdir -p " + cronPath)
	}
	file := cronPath + "/" + cronHash
	_ = util.WriteFile(file, []byte(CheckScript(shell)), 0700)
	_, _ = util.ExecShell("chmod 750 " + cronPath)
	return nil
}

// GetDay 获取任务构造Day
func GetDay(cycle *request.CronCycle) string {
	cronConfig := fmt.Sprintf("%d %d * * *", cycle.Minute, cycle.Hour)
	return cronConfig
}

// GetDayN 获取任务构造Day_N
func GetDayN(cycle *request.CronCycle) string {
	cronConfig := fmt.Sprintf("%d %d */%d * * ", cycle.Minute, cycle.Hour, cycle.Inc)
	return cronConfig
}

// GetHour 获取任务构造Hour
func GetHour(cycle *request.CronCycle) string {
	cronConfig := fmt.Sprintf("%d * * * * ", cycle.Minute)
	return cronConfig
}

// GetHourN 获取任务构造Hour_N
func GetHourN(cycle *request.CronCycle) string {
	cronConfig := fmt.Sprintf("%d */%d * * * ", cycle.Minute, cycle.Hour)
	return cronConfig
}

// GetMinuteN 获取任务构造Minute_N
func GetMinuteN(cycle *request.CronCycle) string {
	cronConfig := fmt.Sprintf("*/%d * * * * ", cycle.Minute)
	return cronConfig
}

// GetWeek 获取任务构造Week
func GetWeek(cycle *request.CronCycle) string {
	cronConfig := fmt.Sprintf("%d %d * * %d", cycle.Minute, cycle.Hour, cycle.Inc)
	return cronConfig
}

// GetMonth 获取任务构造Month
func GetMonth(cycle *request.CronCycle) string {
	cronConfig := fmt.Sprintf("%d %d %d * * ", cycle.Minute, cycle.Hour, cycle.Inc)
	return cronConfig
}

// CheckScript 检查脚本
func CheckScript(shell string) string {
	keys := []string{"shutdown", "init 0", "mkfs", "passwd", "chpasswd", "--stdin", "mkfs.ext", "mke2fs"}
	for _, key := range keys {
		strings.Replace(shell, key, "[***]", -1)
	}
	return shell
}

// GetCronRootPath 获取cron文件路径
func GetCronRootPath() string {
	cronFilePath := "/var/spool/cron/root"
	if util.PathExists("/usr/bin/apt-get") {
		cronFilePath = "/var/spool/cron/crontabs/root"
	}
	_ = os.MkdirAll(filepath.Dir(cronFilePath), 472)
	return cronFilePath
}

// AddShellToCron 将Shell脚本写到Cron文件
func AddShellToCron(shell string) error {
	cronFilePath := GetCronRootPath()
	//判断file是否存在
	if !util.PathExists(cronFilePath) {
		_ = util.WriteFile(cronFilePath, []byte(""), 0644)
		_, _ = util.ExecShell("chmod 600 '" + cronFilePath + "' && chown root.root " + cronFilePath)
	}
	cronBody, err := util.ReadFileStringBody(cronFilePath)
	if err != nil {
		return err
	}
	cronBody += shell + "\n"
	err = util.WriteFile(cronFilePath, []byte(cronBody), 0644)
	if err != nil {
		return err
	}
	_, _ = util.ExecShell("chmod 600 " + cronFilePath + " && chown root.root " + cronFilePath)
	//重载
	ReloadCron()
	return nil
}

// ReloadCron 重载计划任务配置
func ReloadCron() {
	if util.PathExists("/etc/init.d/crond") {
		_, _ = util.ExecShell("/etc/init.d/crond restart")
	}
	if util.PathExists("/etc/init.d/cron") {
		_, _ = util.ExecShell("/etc/init.d/cron restart")
	} else {
		_, _ = util.ExecShell("systemctl restart crond")
	}
	_, _ = util.ExecShell("systemctl restart crond")
}

// XssFilter xss过滤
func XssFilter(content string) string {
	return template.HTMLEscapeString(content) //函数用于将字符串中的特殊字符转换为转义字符，以便在 HTML 中安全显示。
}

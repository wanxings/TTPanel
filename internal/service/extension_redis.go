package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model/request"
	response2 "TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ExtensionRedisService struct{}

// Info 获取Redis详细信息
func (s *ExtensionRedisService) Info() (*response2.ExtensionsInfoResponse, error) {
	var redisInfo response2.ExtensionsInfoResponse
	err := ReadExtensionsInfo(constant.ExtensionRedisName, &redisInfo)
	if err != nil {
		return nil, err
	}
	//获取版本号和安装状态
	if version, ok := s.GetVersion(); ok {
		redisInfo.Description.Version = version
		redisInfo.Description.Install = true
		//获取运行状态
		redisInfo.Description.Status = s.IsRunning()
	}
	return &redisInfo, nil
}

// Install 安装Redis
func (s *ExtensionRedisService) Install(version string) error {
	//验证版本号是否正确
	if ok := util.IsVersion(version); !ok {
		return errors.New("version format error")
	}
	//检查是否在等待或者进行队列中
	taskName := "安装[Redis-" + version + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.InstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	//添加面板队列任务
	execStr := fmt.Sprintf(
		`cd %s && /bin/bash install_lib.sh && cd %s && /bin/bash install.sh install %s`,
		global.Config.System.PanelPath+"/data/shell", s.GetShellPath(), version)
	err := AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}

	return nil
}

// Uninstall 卸载Redis
func (s *ExtensionRedisService) Uninstall(version string) error {
	//检查是否在等待或者进行队列中
	taskName := "卸载[Redis-" + version + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.UninstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash install.sh uninstall %s`, s.GetShellPath(), version)
	err := AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}

	return nil
}

// SetStatus 设置Redis状态
func (s *ExtensionRedisService) SetStatus(action string) (err error) {
	redisInfo, err := s.Info()
	if err != nil {
		return err
	}
	cmdStr := ""
	switch action {
	case constant.ProcessCommandByStart:
		//启动
		cmdStr = redisInfo.Description.InitShell + " " + action
	case constant.ProcessCommandByStop:
		//关闭
		cmdStr = redisInfo.Description.InitShell + " " + action
	case constant.ProcessCommandByRestart:
		//重启
		cmdStr = redisInfo.Description.InitShell + " " + action
	case constant.ProcessCommandByReload:
		//重载配置
		cmdStr = redisInfo.Description.InitShell + " " + action
	default:
		return errors.New("? you can you up")
	}
	//执行命令
	output, err := util.ExecShell(cmdStr)
	if err != nil {
		return err
	}
	if !strings.Contains(output, "done") && !strings.Contains(output, "success") {
		return errors.New(output)
	}
	return nil
}

// PerformanceConfig 获取性能配置
func (s *ExtensionRedisService) PerformanceConfig() ([]*response2.RedisGeneralConfigP, error) {
	configBody, err := s.ReadConfig()
	if err != nil {
		return nil, err
	}

	configData := make([]*response2.RedisGeneralConfigP, len(response2.DefaultRedisGeneralConfigs))
	copy(configData, response2.DefaultRedisGeneralConfigs)

	for _, v := range configData {
		rep := regexp.MustCompile(fmt.Sprintf(`\n%s\s+(.+)`, v.Name))
		group := rep.FindStringSubmatch(configBody)
		if group != nil {
			switch v.Name {
			case "maxmemory":
				v.Value = strconv.Itoa(len(group[1]) / 1024 / 1024)
			default:
				v.Value = group[1]
			}
		}

	}
	return configData, nil

}

// SavePerformanceConfig 保存性能配置
func (s *ExtensionRedisService) SavePerformanceConfig(param request.RedisSavePerformanceConfigR) error {
	//校验ip地址
	if ok := util.CheckIP(param.Bind); !ok {
		return errors.New(helper.MessageWithMap("IPFormatError", map[string]any{"IP": param.Bind}))
	}
	//校验端口
	if ok := util.CheckPort(strconv.Itoa(param.Port)); !ok {
		return errors.New(helper.MessageWithMap("PortFormatError", map[string]any{"Port": param.Port}))
	}
	//校验密码
	if !util.StrIsEmpty(param.RequirePass) {
		prep := "[\\~\\`\\/\\=\\&]"
		match, _ := regexp.MatchString(prep, param.RequirePass)
		if match {
			return errors.New(helper.Message("redis.PasswordFormatError"))
		}
	}

	if ok := util.IsPublicIP(param.Bind); ok {
		if util.StrIsEmpty(param.RequirePass) {
			return errors.New(helper.Message("redis.MustSetPassword"))
		}
	}
	//转换内存
	param.MaxMemory = param.MaxMemory * 1024 * 1024

	//读取配置文件
	configBody, err := s.ReadConfig()
	if err != nil {
		return err
	}
	//备份配置文件
	err = s.BackupConfig(configBody)
	if err != nil {
		return err
	}

	// 遍历结构体字段
	t := reflect.TypeOf(param)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := reflect.ValueOf(param).Field(i)
		rep := regexp.MustCompile("\n" + field.Tag.Get("json") + "\\s+(.*)")

		if rep.MatchString(configBody) {
			newConf := fmt.Sprintf("\n%s %v", field.Tag.Get("json"), value)
			fmt.Println(newConf)
			configBody = rep.ReplaceAllString(configBody, newConf)
		} else {
			newConf := "\n# Redis configuration file example."
			newConf = newConf + fmt.Sprintf("\n%s %v", field.Tag.Get("json"), value)
			fmt.Println(newConf)
			configBody = configBody + newConf
		}
	}

	//写入配置文件
	err = s.WriteConfig(configBody)
	if err != nil {
		return err
	}

	return nil

}

// LoadStatus 获取Redis状态
func (s *ExtensionRedisService) LoadStatus() (map[string]string, error) {
	configBody, err := s.ReadConfig()
	if err != nil {
		return nil, err
	}

	//获取端口
	rep := regexp.MustCompile(`\nport\s+(.+)`)
	group := rep.FindStringSubmatch(configBody)
	if group == nil {
		return nil, errors.New("not found port")
	}
	port := group[1]
	//获取密码
	rep = regexp.MustCompile(`\nrequirepass\s+(.+)`)
	group = rep.FindStringSubmatch(configBody)
	password := ""
	if group != nil {
		password = group[1]
	}

	//获取bind
	rep = regexp.MustCompile(`\nbind\s+(.*?)\s+`)
	group = rep.FindStringSubmatch(configBody)
	if group == nil {
		return nil, errors.New("not found bind")
	}
	bind := group[1]

	cmdStr := fmt.Sprintf("%s -h %s -p %s", s.GetCliBinPath(), bind, port)
	if !util.StrIsEmpty(password) {
		cmdStr += fmt.Sprintf(" -a %s", password)
	}
	cmdStr += " info"
	shellResult, err := util.ExecShell(cmdStr)
	if err != nil {
		return nil, errors.New(err.Error() + ";output:" + shellResult)
	}
	data := strings.Split(shellResult, "\r\n")
	result := make(map[string]string)
	for _, v := range data {
		configs := strings.Split(v, ":")
		if len(configs) == 2 {
			result[configs[0]] = configs[1]
		}

	}
	return result, nil
}

// PersistentConfig 获取持久化配置
func (s *ExtensionRedisService) PersistentConfig() (map[string]interface{}, error) {
	returnData := make(map[string]interface{})
	configBody, err := s.ReadConfig()
	if err != nil {
		return nil, err
	}
	rep := regexp.MustCompile(`\nsave[\w\s\n]+`)
	group := rep.FindStringSubmatch(configBody)
	var RDBSave []map[string]string
	if group != nil {
		global.Log.Debugf("PersistentConfig->FindStringSubmatch group[0]:%v \n", group[0])
		saves := strings.Split(group[0], "\n")
		for _, save := range saves {
			saveL := strings.Split(save, " ")
			if len(saveL) > 1 {
				RDBSave = append(RDBSave, map[string]string{
					"time":    saveL[1],
					"changes": saveL[2],
				})
			}

		}
	}
	returnData["rdb"] = RDBSave

	//获取持久化保存目录
	rep = regexp.MustCompile(`\ndir\s+(.+)`)
	group = rep.FindStringSubmatch(configBody)
	if group == nil {
		returnData["dir"] = nil
	} else {
		returnData["dir"] = group[1]
	}

	//获取AOF配置
	//获取appendonly
	rep = regexp.MustCompile(`\nappendonly\s+(\w+)`)
	group = rep.FindStringSubmatch(configBody)
	if group == nil {
		return nil, errors.New("not found appendonly")
	}
	appendOnly := group[1]
	//获取appendfsync
	rep = regexp.MustCompile(`\nappendfsync\s+(\w+)`)
	group = rep.FindStringSubmatch(configBody)
	if group == nil {
		return nil, errors.New("not found appendfsync")
	}
	appendFsync := group[1]
	returnData["aof"] = map[string]string{
		"appendonly":  appendOnly,
		"appendfsync": appendFsync,
	}
	return returnData, err
}

// SavePersistentConfig 保存持久化配置
func (s *ExtensionRedisService) SavePersistentConfig(param *request.RedisSavePersistentConfigR) error {
	configBody, err := s.ReadConfig()
	if err != nil {
		return err
	}

	//备份配置文件
	err = s.BackupConfig(configBody)
	if err != nil {
		return err
	}

	//校验dir配置
	if !util.StrIsEmpty(param.Dir) {
		if util.PathExists(strings.TrimSuffix(param.Dir, "/")) && util.IsDir(strings.TrimSuffix(param.Dir, "/")) {
			_ = util.CreateDir(param.Dir + "redis_cache")
			//修改文件夹权限
			_ = util.ChangePermission(param.Dir+"redis_cache", "redis", "755")
			newConf := fmt.Sprintf("\ndir %s", param.Dir)
			rep := regexp.MustCompile(`\ndir\s+(.+)`)
			configBody = rep.ReplaceAllString(configBody, newConf)
		} else {
			return errors.New(fmt.Sprintf("Not a directory or directory does not exist  path:%v", strings.TrimSuffix(param.Dir, "/")))
		}
	}

	//校验rdb配置
	if param.Rdb != nil {
		newConf := "\n"
		for _, v := range param.Rdb {
			newConf += fmt.Sprintf("save %d %d \n", v.Time, v.Changes)
		}
		rep := regexp.MustCompile(`\nsave[\w\s\n]+`)
		configBody = rep.ReplaceAllString(configBody, newConf)
	}

	//校验aof配置
	if param.Aof != nil {
		newConf := fmt.Sprintf("\nappendonly %s", param.Aof.AppendOnly)
		rep := regexp.MustCompile(`\nappendonly\s+(\w+)`)
		configBody = rep.ReplaceAllString(configBody, newConf)
		newConf = fmt.Sprintf("\nappendfsync %s", param.Aof.AppendFsync)
		rep = regexp.MustCompile(`\nappendfsync\s+(\w+)`)
		configBody = rep.ReplaceAllString(configBody, newConf)
	}

	//写入配置文件
	err = s.WriteConfig(configBody)
	if err != nil {
		return err
	}
	return nil
}

// GetVersion  返回版本号和安装状态
func (s *ExtensionRedisService) GetVersion() (string, bool) {
	//检查是否已安装
	if util.PathExists(s.GetBinPath()) {
		result, _ := util.ExecShell("/www/server/redis/src/redis-server -v|awk '{print $3}'|cut -f2 -d'='")
		//去除换行和空格
		result = util.ClearStr(result)
		return result, true
	} else {
		return "", false
	}
}

func (s *ExtensionRedisService) GetBinPath() string {
	return global.Config.System.ServerPath + "/redis/src/redis-server"
}

func (s *ExtensionRedisService) GetCliBinPath() string {
	return global.Config.System.ServerPath + "/redis/src/redis-cli"
}

func (s *ExtensionRedisService) IsRunning() bool {
	result, _ := util.ExecShell("/etc/init.d/redis status")
	//去除换行和空格
	result = util.ClearStr(result)
	if strings.Contains(result, "running") {
		return true
	} else {
		return false
	}
}

func (s *ExtensionRedisService) GetShellPath() any {
	return global.Config.System.PanelPath + "/data/extensions/redis/install"
}

func (s *ExtensionRedisService) ReadConfig() (string, error) {
	configPath := fmt.Sprintf("%s/redis/redis.conf", global.Config.System.ServerPath)
	file, err := util.ReadFileStringBody(configPath)
	if err != nil {
		return "", err
	}
	return file, nil
}

// WriteConfig 写入配置文件
func (s *ExtensionRedisService) WriteConfig(body string) error {
	configPath := fmt.Sprintf("%s/redis/redis.conf", global.Config.System.ServerPath)
	err := util.WriteFile(configPath, []byte(body), 0600)
	if err != nil {
		return err
	}
	return nil
}

func (s *ExtensionRedisService) BackupConfig(body string) error {
	configPath := fmt.Sprintf("%s/redis/redis.conf.%d", global.Config.System.ServerPath, time.Now().Unix())
	err := util.WriteFile(configPath, []byte(body), 0600)
	if err != nil {
		return err
	}
	return nil
}

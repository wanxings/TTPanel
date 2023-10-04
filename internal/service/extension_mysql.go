package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model/request"
	response2 "TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ExtensionMysqlService struct{}

// Info 获取Mysql详细信息
func (s *ExtensionMysqlService) Info() (*response2.ExtensionsInfoResponse, error) {
	var mysqlInfo response2.ExtensionsInfoResponse
	err := ReadExtensionsInfo(constant.ExtensionMysqlName, &mysqlInfo)
	if err != nil {
		return nil, err
	}

	//判断是否安装
	if version, ok := s.IsInstalled(); ok {
		mysqlInfo.Description.Version = version
		mysqlInfo.Description.Install = true
		//获取运行状态
		mysqlInfo.Description.Status = s.IsRunning()
	} else {
		mysqlInfo.Description.Version = ""
		mysqlInfo.Description.Install = false
		mysqlInfo.Description.Status = false
	}

	return &mysqlInfo, nil
}

// Install 安装Mysql
func (s *ExtensionMysqlService) Install(version string) error {
	//验证版本号是否正确 Todo: 不同mysql版本此处校验有问题，后续优化
	//if ok := util.IsMysqlVersion(version); !ok {
	//	return errors.New("版本号错误")
	//}
	//检查是否在等待或者进行队列中
	taskName := "安装[Mysql-" + version + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.InstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash install_lib.sh && cd %s && /bin/bash install.sh install %s`,
		global.Config.System.PanelPath+"/data/shell", s.GetShellPath(), version)
	err := AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}

	return nil
}

// Uninstall 卸载Mysql
func (s *ExtensionMysqlService) Uninstall(version string) error {
	//验证版本号是否正确
	if ok := util.IsMysqlVersion(version); !ok {
		return errors.New("版本号错误")
	}

	//检查是否在等待或者进行队列中
	taskName := "卸载[Mysql-" + version + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.UninstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	//检查面板上是否有已存在的本地数据库 Todo: 卸载mysql前如果有存在的数据库阻止卸载，未完成,不确定需求是否合理，
	//_, i, err := (&model.Databases{}).List(global.PanelDB, &model.ConditionsT{"sid = ?": 0}, 1, 1)
	//if err != nil {
	//	return err
	//}
	//if i > 0 {
	//	return errors.New(helper.MessageWithMap("UninstallingOrWaiting", map[string]any{"Name": "Mysql"}))
	//}

	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash install.sh uninstall %s`, s.GetShellPath(), version)
	err := AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}
	return nil
}

// SetStatus 设置Mysql运行状态
func (s *ExtensionMysqlService) SetStatus(action string) error {
	cmdStr := ""
	switch action {
	case constant.ProcessCommandByStart:
		//启动
		cmdStr = "/etc/init.d/mysqld " + action
	case constant.ProcessCommandByStop:
		//关闭
		cmdStr = "/etc/init.d/mysqld " + action
	case constant.ProcessCommandByRestart:
		//重启
		cmdStr = "/etc/init.d/mysqld " + action
	case constant.ProcessCommandByReload:
		//重载配置
		cmdStr = "/etc/init.d/mysqld " + action
	default:
		return errors.New("无效动作")
	}
	//执行命令 Todo:这里需要优化,出现mysqld_safe A mysqld process already exists时缓冲区会阻塞住,
	err := util.ExecShellScriptS(cmdStr)
	if err != nil {
		return err
	}
	return nil
}

// PerformanceConfig 性能配置
func (s *ExtensionMysqlService) PerformanceConfig() ([]*response2.MysqlGeneralConfigP, error) {
	configData := make([]*response2.MysqlGeneralConfigP, len(response2.DefaultMysqlGeneralConfigs))
	copy(configData, response2.DefaultMysqlGeneralConfigs)
	dbService, err := GroupApp.DatabaseMysqlServiceApp.NewMysqlServiceBySid(0)
	if err != nil {
		return nil, errors.New("连接数据库失败，检查数据库是否正常")
	}
	rows, err := dbService.Raw("SHOW VARIABLES").Rows()
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		var key string
		var value string
		if err := rows.Scan(&key, &value); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s %s\n", key, value)
		// 处理变量名称和值
		for _, config := range configData {
			if config.Name == key {
				config.Value = value
				break
			}
		}
	}
	return configData, nil
}

// SavePerformanceConfig 设置性能配置
func (s *ExtensionMysqlService) SavePerformanceConfig(configs request.MysqlSetPerformanceConfigR) error {
	return nil
}

type MasterStatus struct {
	File     string
	Position uint
}

// LoadStatus 获取Mysql负载状态
func (s *ExtensionMysqlService) LoadStatus() (*response2.MysqlLoadStatusP, error) {
	dbService, err := GroupApp.DatabaseMysqlServiceApp.NewMysqlServiceBySid(0)
	if err != nil {
		return nil, errors.New("连接数据库失败，检查数据库是否正常")
	}
	rows, err := dbService.Raw("SHOW GLOBAL STATUS").Rows()
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	sourceData := make(map[string]string)
	for rows.Next() {
		var key string
		var value string
		if err := rows.Scan(&key, &value); err != nil {
			fmt.Println(err)
		}
		sourceData[key] = value
	}
	var tmp response2.MysqlLoadStatusP
	fmt.Println("Uptime", sourceData["Uptime"])
	// 遍历结构体字段
	t := reflect.ValueOf(&tmp).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		//value := reflect.ValueOf(tmp).Field(i)

		for k, v := range sourceData {
			if t.Type().Field(i).Tag.Get("json") == k {
				//value.SetString(v)
				field.SetString(v)
				break
			}
		}
	}

	if util.StrIsEmpty(tmp.Run) {
		uptime, _ := strconv.ParseInt(sourceData["Uptime"], 10, 64)
		tmp.Run = fmt.Sprintf("%d", time.Now().Unix()-uptime)
	}

	var result MasterStatus
	if err := dbService.Raw("SHOW MASTER STATUS").Scan(&result).Error; err != nil {
		return nil, err
	}
	tmp.File = result.File
	tmp.Position = fmt.Sprintf("%d", result.Position)

	return &tmp, nil
}

// ErrorLog 获取Mysql错误日志
func (s *ExtensionMysqlService) ErrorLog() (string, error) {
	cnfMap, err := s.GetMysqlCnf()
	if err != nil {
		return "", err
	}

	var errLogPath string
	dataDir := cnfMap["datadir"]
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			if strings.HasSuffix(file.Name(), ".err") {
				errLogPath = dataDir + "/" + file.Name()
			}
		}
	}
	if util.StrIsEmpty(errLogPath) {
		return "", errors.New("未找到错误日志")
	}

	return errLogPath, nil

}

// SlowLogs 获取Mysql慢查询日志
func (s *ExtensionMysqlService) SlowLogs() (string, error) {
	cnfMap, err := s.GetMysqlCnf()
	if err != nil {
		return "", err
	}
	slowLogsPath := cnfMap["datadir"] + "/mysql-slow.log"
	return slowLogsPath, nil
}

func (s *ExtensionMysqlService) GetShellPath() string {
	return global.Config.System.PanelPath + "/data/extensions/mysql/install"
}

func (s *ExtensionMysqlService) IsInstalled() (string, bool) {
	//检查是否已安装
	if util.PathExists(s.GetBinPath()) {
		result, _ := util.ExecShell("/www/server/mysql/bin/mysql -V | grep -oE '[0-9]+\\.[0-9]+\\.[0-9]+'")
		//去除换行和空格
		result = util.ClearStr(result)
		return result, true
	} else {
		return "", false
	}
}

// GetBinPath 取Mysql的bin路径
func (s *ExtensionMysqlService) GetBinPath() string {
	return global.Config.System.ServerPath + "/mysql/bin/mysql"
}

func (s *ExtensionMysqlService) IsRunning() bool {
	result, _ := util.ExecShell("/etc/init.d/mysqld status")
	//去除换行和空格
	result = util.ClearStr(result)
	if strings.Contains(result, "SUCCESS!MySQLrunning") || strings.Contains(result, "active(running)") {
		return true
	} else {
		return false
	}
}

// GetMysqlCnf 获取数据库配置信息
func (s *ExtensionMysqlService) GetMysqlCnf() (map[string]string, error) {
	data := make(map[string]string)
	cnfFile := "/etc/my.cnf"
	if !util.PathExists(cnfFile) {
		return nil, errors.New("my.cnf文件不存在")
	}
	cnfBody, err := util.ReadFileStringBody(cnfFile)
	if err != nil {
		return nil, err
	}
	rep := "datadir\\s*=\\s*(.+)\\n"
	match := regexp.MustCompile(rep).FindStringSubmatch(cnfBody)
	if len(match) > 1 {
		data["datadir"] = match[1]
	}
	rep = "port\\s*=\\s*([0-9]+)\\s*\\n"
	match = regexp.MustCompile(rep).FindStringSubmatch(cnfBody)
	if len(match) > 1 {
		data["port"] = match[1]
	}

	if len(data) == 0 {
		data["datadir"] = "/www/server/data"
		data["port"] = "3306"
	}
	return data, nil
}

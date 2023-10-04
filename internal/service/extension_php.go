package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	fcgiClient "github.com/tomasen/fcgi_client"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ExtensionPHPService struct{}

var PHPBinPath = "/usr/bin/php"

// Info 获取php扩展基本信息
func (s *ExtensionPHPService) Info() ([]*response.ExtensionsInfoResponse, error) {
	var phpItemInfo []*response.ExtensionsInfoResponse
	err := ReadExtensionsInfo(constant.ExtensionPHPName, &phpItemInfo)
	if err != nil {
		return nil, errors.New("php扩展信息文件损坏")
	}
	for _, phpInfo := range phpItemInfo {
		phpInfo.Description.Install = s.IsInstalled(phpInfo.VersionList[0].MVersion)
		if phpInfo.Description.Install {
			phpInfo.Description.Status = s.Status(phpInfo.VersionList[0].MVersion)
		} else {
			phpInfo.Description.Status = false
		}
	}
	return phpItemInfo, nil
}

// IsInstalled 指定版本的PHP是否已安装
func (s *ExtensionPHPService) IsInstalled(version string) bool {
	//检查是否已安装
	if util.PathExists(s.GetBinPath(version)) {
		return true
	} else {
		return false
	}
}

// Status 获取指定版本的PHP的状态
func (s *ExtensionPHPService) Status(version string) bool {
	cmd := "ps -ef | grep 'php/" + version + "' | grep -v grep | grep -v python | awk '{print $2}'"
	output, _ := util.ExecShell(cmd)
	pid := strings.TrimSpace(output)
	if util.StrIsEmpty(pid) {
		return false
	}
	return true
}

// Install 安装指定版本的PHP
func (s *ExtensionPHPService) Install(version string) error {
	//验证版本号是否正确
	if ok := util.IsPHPVersion(version); !ok {
		return errors.New("version Error")
	}
	//检查是否在等待或者进行队列中
	taskName := "安装[PHP-" + version + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.InstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash install_lib.sh && cd %s && /bin/bash install_php.sh install %s`,
		global.Config.System.PanelPath+"/data/shell", s.GetShellPath(), version)
	err := AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}
	return nil
}

// Uninstall 卸载指定版本的PHP
func (s *ExtensionPHPService) Uninstall(version string) error {
	//验证版本号是否正确
	if ok := util.IsPHPVersion(version); !ok {
		return errors.New("version Error")
	}

	//检查是否在等待或者进行队列中
	taskName := "卸载[PHP-" + version + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.UninstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash install_php.sh uninstall %s`, s.GetShellPath(), version)
	err := AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}

	return nil
}

// SetStatus 设置指定版本的PHP状态 action: start/stop/restart/reload
func (s *ExtensionPHPService) SetStatus(version string, action string) error {
	err := s.CheckInstall(version)
	if err != nil {
		return err
	}
	//验证版本号是否正确
	if ok := util.IsPHPVersion(version); !ok {
		return errors.New("version Error")
	}
	serverName := "php-fpm-" + version
	cmdStr := ""
	switch action {
	case constant.ProcessCommandByStart:
		//启动
		cmdStr = "/etc/init.d/" + serverName + " " + action
	case constant.ProcessCommandByStop:
		//关闭
		cmdStr = "/etc/init.d/" + serverName + " " + action
	case constant.ProcessCommandByRestart:
		//重启
		cmdStr = "/etc/init.d/" + serverName + " " + action
	case constant.ProcessCommandByReload:
		//重载配置
		cmdStr = "/etc/init.d/" + serverName + " " + action
	default:
		return errors.New("action is empty")
	}
	//执行命令
	output, err := util.ExecShell(cmdStr)
	if err != nil {
		return err
	}
	global.Log.Debugf("SetStatus->ExecShell output: %s \n", output)
	return nil
}

// ExtensionList php扩展列表
func (s *ExtensionPHPService) ExtensionList(version string) ([]*response.LibP, error) {
	err := s.CheckInstall(version)
	if err != nil {
		return nil, err
	}
	phpIni, err := s.readPHPIni(version)
	if err != nil {
		return nil, err
	}
	libJsons, err := s.readPHPLibJson()
	if err != nil {
		return nil, err
	}
	//获取正在运行的面板任务
	queueTasks, _, err := (&model.QueueTask{}).List(global.PanelDB, &model.ConditionsT{"FIXED": "status = " + fmt.Sprintf("%d", constant.QueueTaskStatusProcessing) + " OR status = " + fmt.Sprintf("%d", constant.QueueTaskStatusWait)}, 0, 0)
	if err != nil {
		return nil, err
	}
	for _, libJson := range libJsons {
		libJson.TaskStatus = constant.QueueTaskStatusSuccess
		for _, queueTask := range queueTasks {
			if strings.Contains(queueTask.ExecStr, libJson.Name) {
				libJson.TaskStatus = queueTask.Status
			}
		}
		if strings.Index(phpIni, libJson.Check) == -1 {
			libJson.Status = false
		} else {
			libJson.Status = true
		}
	}

	return libJsons, nil
}

// InstallLib 安装php扩展
func (s *ExtensionPHPService) InstallLib(version string, name string) error {
	err := s.CheckInstall(version)
	if err != nil {
		return err
	}

	//检查是否存在该扩展
	libByName, err := s.getLibByName(name)
	if err != nil {
		return err
	}
	if libByName == nil {
		return errors.New("not found lib")
	}

	//检查是否在等待或者进行队列中
	taskName := "安装[" + name + "-" + version + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.InstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash install_lib.sh && cd %s && /bin/bash install_lib.sh install %s %s`,
		global.Config.System.PanelPath+"/data/shell", s.GetShellPath(), name, version)
	err = AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}
	return nil
}

// UninstallLib 卸载php扩展
func (s *ExtensionPHPService) UninstallLib(version string, name string) error {
	err := s.CheckInstall(version)
	if err != nil {
		return err
	}

	//检查是否存在该扩展
	libByName, err := s.getLibByName(name)
	if err != nil {
		return err
	}
	if libByName == nil {
		return errors.New("not found lib")
	}

	//检查是否在等待或者进行队列中
	taskName := "卸载[" + name + "-" + version + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.UninstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash install_lib.sh uninstall %s %s`, s.GetShellPath(), name, version)
	err = AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}

	return nil
}

// GeneralConfig 获取php通用配置文件
func (s *ExtensionPHPService) GeneralConfig(version string) ([]response.PHPGeneralConfigP, error) {
	configData := make([]response.PHPGeneralConfigP, len(response.DefaultPHPGeneralConfigs))
	copy(configData, response.DefaultPHPGeneralConfigs)
	phpIniS, err := s.readPHPIni(version)
	if err != nil {
		return configData, err
	}
	for index, config := range configData {
		rep := config.Name + `\s*=\s*([0-9A-Za-z_&/ ~]+)(\s*;?|\r?\n)`
		re := regexp.MustCompile(rep)
		tmp := re.FindStringSubmatch(phpIniS)
		if len(tmp) == 0 {
			continue
		}
		configData[index].Value = tmp[1]
	}
	return configData, nil
}

// SaveGeneralConfig 保存php配置文件
func (s *ExtensionPHPService) SaveGeneralConfig(param request.SaveGeneralConfigR) error {
	//读取Ini配置文件
	phpIni, err := s.readPHPIni(param.Version)
	if err != nil {
		return err
	}

	// 遍历结构体字段
	t := reflect.TypeOf(param)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := reflect.ValueOf(param).Field(i)
		if field.Name == "Version" {
			continue
		}
		global.Log.Debugf("SaveGeneralConfig->for key:%s value%v\n", field.Name, value)
		rep := regexp.MustCompile(fmt.Sprintf("%s\\s*=\\s*(.+)\\r?\\n", field.Tag.Get("json")))
		val := fmt.Sprintf("%s = %s\n", field.Tag.Get("json"), value.Interface())
		phpIni = rep.ReplaceAllString(phpIni, val)
	}
	err = s.savePHPIni(param.Version, phpIni)
	if err != nil {
		return err
	}
	return nil
}

// DisableFunctionList 禁用函数列表
func (s *ExtensionPHPService) DisableFunctionList(version string) (string, error) {
	err := s.CheckInstall(version)
	if err != nil {
		return "", err
	}

	//读取Ini配置文件
	phpIni, err := s.readPHPIni(version)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`disable_functions\s*=\s{0,1}(.*)\n`)
	match := re.FindStringSubmatch(phpIni)
	if len(match) > 1 {
		disableFunctions := match[1]
		return disableFunctions, nil
	}

	return "", nil
}

// DeleteDisableFunction 删除禁用函数
func (s *ExtensionPHPService) DeleteDisableFunction(version string, name string) error {
	err := s.CheckInstall(version)
	if err != nil {
		return err
	}
	//获取禁用函数列表
	disableFunctions, err := s.DisableFunctionList(version)
	if err != nil {
		return err
	}
	//检查是否已存在,如果不存在则直接返回
	if !strings.Contains(disableFunctions, name) {
		return nil
	}
	disableFunctions = strings.Replace(disableFunctions, name, "", -1)
	disableFunctions = strings.Replace(disableFunctions, ",,", ",", -1)
	disableFunctions = strings.Trim(disableFunctions, ",")
	//读取Ini配置文件
	phpIni, err := s.readPHPIni(version)
	if err != nil {
		return err
	}
	//替换
	rep := regexp.MustCompile(`disable_functions\s*=\s{0,1}(.*)\n`)
	phpIni = rep.ReplaceAllString(phpIni, fmt.Sprintf("disable_functions = %s\n", disableFunctions))
	//保存
	err = s.savePHPIni(version, phpIni)
	if err != nil {
		return err
	}
	return nil
}

// AddDisableFunction 添加禁用函数
func (s *ExtensionPHPService) AddDisableFunction(version string, name string) error {
	err := s.CheckInstall(version)
	if err != nil {
		return err
	}
	//获取禁用函数列表
	disableFunctions, err := s.DisableFunctionList(version)
	if err != nil {
		return err
	}
	//检查是否已存在
	if strings.Contains(disableFunctions, name) {
		return nil
	}
	disableFunctions = disableFunctions + "," + name
	//读取Ini配置文件
	phpIni, err := s.readPHPIni(version)
	if err != nil {
		return err
	}
	//替换
	rep := regexp.MustCompile(`disable_functions\s*=\s{0,1}(.*)\n`)
	phpIni = rep.ReplaceAllString(phpIni, fmt.Sprintf("disable_functions = %s\n", disableFunctions))
	//保存
	err = s.savePHPIni(version, phpIni)
	if err != nil {
		return err
	}
	return nil
}

// PerformanceConfig 性能配置
func (s *ExtensionPHPService) PerformanceConfig(version string) (*response.PHPPerformanceConfigP, error) {
	err := s.CheckInstall(version)
	if err != nil {
		return nil, err
	}
	//读取php-fpm.conf配置文件
	fpmConf, err := s.readPHPFpmConf(version)
	if err != nil {
		return nil, err
	}

	var configData response.PHPPerformanceConfigP
	rep := regexp.MustCompile(`\s*pm.max_children\s*=\s*([0-9]+)\s*`)
	tmp := rep.FindStringSubmatch(fpmConf)
	tmp[1] = strings.TrimSpace(tmp[1])
	configData.MaxChildren, err = strconv.Atoi(tmp[1])
	if err != nil {
		configData.MaxChildren = 0
	}

	rep = regexp.MustCompile(`\s*pm.start_servers\s*=\s*([0-9]+)\s*`)
	tmp = rep.FindStringSubmatch(fpmConf)
	tmp[1] = strings.TrimSpace(tmp[1])
	configData.StartServers, err = strconv.Atoi(tmp[1])
	if err != nil {
		configData.StartServers = 0
	}

	rep = regexp.MustCompile(`\s*pm.min_spare_servers\s*=\s*([0-9]+)\s*`)
	tmp = rep.FindStringSubmatch(fpmConf)
	tmp[1] = strings.TrimSpace(tmp[1])
	configData.MinSpareServers, err = strconv.Atoi(tmp[1])
	if err != nil {
		configData.MinSpareServers = 0
	}

	rep = regexp.MustCompile(`\s*pm.max_spare_servers\s*=\s*([0-9]+)\s*`)
	tmp = rep.FindStringSubmatch(fpmConf)
	tmp[1] = strings.TrimSpace(tmp[1])
	configData.MaxSpareServers, err = strconv.Atoi(tmp[1])
	if err != nil {
		configData.MaxSpareServers = 0
	}

	rep = regexp.MustCompile(`\s*pm\s*=\s*(\w+)\s*`)
	tmp = rep.FindStringSubmatch(fpmConf)
	configData.Pm = tmp[1]

	rep = regexp.MustCompile(`\s*listen.allowed_clients\s*=\s*([\w\.,/]+)\s*`)
	tmp = rep.FindStringSubmatch(fpmConf)
	configData.Allowed = tmp[1]

	configData.Unix, configData.Bind, configData.Port, err = s.getFpmAddress(version, true)
	if err != nil {
		return nil, err
	}

	return &configData, nil
}

// SavePerformanceConfig 保存性能配置
func (s *ExtensionPHPService) SavePerformanceConfig(param *request.SavePerformanceConfigR) error {
	err := s.CheckInstall(param.Version)
	if err != nil {
		return err
	}
	//读取php-fpm.conf配置文件
	fpmConf, err := s.readPHPFpmConf(param.Version)
	if err != nil {
		return err
	}
	//替换
	rep := regexp.MustCompile(`\s*pm.max_children\s*=\s*([0-9]+)\s*`)
	fpmConf = rep.ReplaceAllString(fpmConf, fmt.Sprintf("pm.max_children = %d", param.MaxChildren)) //ok

	rep = regexp.MustCompile(`\s*pm.start_servers\s*=\s*([0-9]+)\s*`)
	fpmConf = rep.ReplaceAllString(fpmConf, fmt.Sprintf("pm.start_servers = %d", param.StartServers)) //ok

	rep = regexp.MustCompile(`\s*pm.min_spare_servers\s*=\s*([0-9]+)\s*`)
	fpmConf = rep.ReplaceAllString(fpmConf, fmt.Sprintf("pm.min_spare_servers = %d", param.MinSpareServers)) //ok

	rep = regexp.MustCompile(`\s*pm.max_spare_servers\s*=\s*([0-9]+)\s*`)
	fpmConf = rep.ReplaceAllString(fpmConf, fmt.Sprintf("pm.max_spare_servers = %d", param.MaxSpareServers)) //ok

	rep = regexp.MustCompile(`\s*pm\s*=\s*(\w+)\s*`)
	fpmConf = rep.ReplaceAllString(fpmConf, fmt.Sprintf("pm = %s", param.Pm)) //ok

	if param.Pm == "ondemand" {
		if match, _ := regexp.MatchString(`listen\.backlog\s*=\s*-1`, fpmConf); match {
			rep = regexp.MustCompile(`\s*listen\.backlog\s*=\s*([0-9-]+)\s*`)
			fpmConf = rep.ReplaceAllString(fpmConf, "\nlisten.backlog = 8192\n")
		}
	}
	var listen string
	if param.Unix == "unix" {
		listen = fmt.Sprintf("/tmp/php-cgi-%s.sock", param.Version)
	} else if param.Unix == "tcp" {
		defaultListen := fmt.Sprintf("127.0.0.1:10%s1", param.Version)
		if strings.Contains(param.Bind, "sock") {
			listen = defaultListen
		} else {
			listen = param.Bind
		}
	}

	rep = regexp.MustCompile(`\s*listen\s*=\s*.+\s*`)
	fpmConf = rep.ReplaceAllString(fpmConf, "\nlisten = "+listen+"\n") //ok

	rep = regexp.MustCompile(`\s*listen.allowed_clients\s*=\s*([\w\.,/]+)\s*`)
	fpmConf = rep.ReplaceAllString(fpmConf, fmt.Sprintf("\nlisten.allowed_clients = %s\n", param.Allowed)) //ok

	//写入php-fpm.conf配置文件
	err = s.writePHPFpmConf(param.Version, fpmConf)
	if err != nil {
		return err
	}
	return nil
}

// LoadStatus 获取PHP的负载状态
func (s *ExtensionPHPService) LoadStatus(version string) (*response.PHPLoadStatusP, error) {
	err := s.CheckInstall(version)
	if err != nil {
		return nil, err
	}
	uri := "/phpfpm_" + version + "_status?json"
	unix, bind, port, err := s.getFpmAddress(version, true)
	if err != nil {
		return nil, err
	}
	fpmAddress := ""
	if unix == "tcp" {
		fpmAddress = bind + ":" + strconv.Itoa(port)
	} else {
		fpmAddress = bind
	}

	fpm, err := s.requestPhpFpm("", uri, "GET", nil, fpmAddress)
	if err != nil {
		return nil, err
	}
	statusData := &response.PHPLoadStatusP{}
	err = json.Unmarshal(fpm, &statusData)
	if err != nil {
		return nil, err
	}
	return statusData, nil
}

// FpmLog FPM日志
func (s *ExtensionPHPService) FpmLog(version string) (string, error) {
	err := s.CheckInstall(version)
	if err != nil {
		return "", err
	}
	conf, err := s.readPHPFpmConf(version)
	if err != nil {
		return "", err
	}
	rep := regexp.MustCompile(`\s*error_log\s*=\s*(.+)\s*`)
	match := rep.FindStringSubmatch(conf)
	if len(match) < 2 {
		global.Log.Errorf("not found FpmLog error_log")
		return "", nil
	}
	logPath := match[1]
	return logPath, nil
}

// FpmSlowLog FPM慢日志
func (s *ExtensionPHPService) FpmSlowLog(version string) (string, error) {
	err := s.CheckInstall(version)
	if err != nil {
		return "", err
	}
	conf, err := s.readPHPFpmConf(version)
	if err != nil {
		return "", err
	}
	rep := regexp.MustCompile(`\s*slowlog\s*=\s*(.+)\s*`)
	match := rep.FindStringSubmatch(conf)
	if len(match) < 2 {
		global.Log.Errorf("not found FpmSlowLog,version:" + version)
		return "", nil
	}
	logPath := strings.TrimSpace(match[1])
	return logPath, nil
}

// PHPInfo PHP信息
func (s *ExtensionPHPService) PHPInfo(version string) (map[string]interface{}, error) {
	err := s.CheckInstall(version)
	if err != nil {
		return nil, err
	}

	serverPath := global.Config.System.ServerPath
	phpPath := serverPath + "/php/"
	phpBin := phpPath + version + "/bin/php"
	phpIni := phpPath + version + "/etc/php.ini"
	phpIniLit := serverPath + "/php/80/etc/php/80/litespeed/php.ini"
	if _, err := os.Stat(phpIniLit); err == nil {
		phpIni = phpIniLit
	}
	fmt.Println(phpBin + " -c " + phpIni + " " + global.Config.System.PanelPath + "/data/extensions/php/php_info.php")
	out, _ := util.ExecShell(phpBin + " -c " + phpIni + " " + global.Config.System.PanelPath + "/data/extensions/php/php_info.php")
	if strings.Contains(out, "Warning: JIT is incompatible") {
		out = strings.Split(strings.TrimSpace(out), "\n")[len(strings.Split(strings.TrimSpace(out), "\n"))-1]
	}
	fmt.Println(out)

	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(out), &result)
	if err != nil {
		return nil, err
	}
	phpInfo := make(map[string]interface{})
	phpInfo["php_version"] = result["php_version"]
	phpInfo["php_path"] = phpPath
	phpInfo["php_bin"] = phpBin
	phpInfo["php_ini"] = phpIni
	phpInfo["modules"] = result["modules"]
	phpInfo["ini"] = result["ini"]
	phpInfo["keys"] = map[string]string{
		"1cache":     "缓存器",
		"2crypt":     "加密解密库",
		"0db":        "数据库驱动",
		"4network":   "网络通信库",
		"5io_string": "文件和字符串处理库",
		"3photo":     "图片处理库",
		"6other":     "其它第三方库",
	}
	result["phpinfo"] = phpInfo
	delete(result, "php_version")
	delete(result, "modules")
	delete(result, "ini")
	return result, nil
}

// PHPInfoHtml PHP信息
func (s *ExtensionPHPService) PHPInfoHtml(version string) (string, error) {
	err := s.CheckInstall(version)
	if err != nil {
		return "", err
	}
	uri := "php_info_html.php"
	unix, bind, port, err := s.getFpmAddress(version, true)
	if err != nil {
		return "", err
	}

	fpmAddress := ""
	if unix == "tcp" {
		fpmAddress = bind + ":" + strconv.Itoa(port)
	} else {
		fpmAddress = bind
	}

	fpm, err := s.requestPhpFpm(global.Config.System.PanelPath+"/data/extensions/php", uri, "GET", nil, fpmAddress)
	if err != nil {
		return "", err
	}
	return string(fpm), nil

}

// CmdVersion 获取php命令行版本
func (s *ExtensionPHPService) CmdVersion() (version string, versionList []string, err error) {
	phpList, err := s.Info()
	if err != nil {
		return
	}
	for _, php := range phpList {
		if php.Description.Install {
			versionList = append(versionList, php.VersionList[0].MVersion)
		}
	}
	if !util.IsLink(PHPBinPath) {
		//不存在则尝试创建
		return
	}
	readlink, err := os.Readlink(PHPBinPath)
	if err != nil {
		return
	}
	for _, v := range versionList {
		if strings.Contains(readlink, "/"+v+"/") {
			version = v
			break
		}
	}
	return
}

// SetCmdVersion 设置php命令行版本
func (s *ExtensionPHPService) SetCmdVersion(version string) error {
	err := s.CheckInstall(version)
	if err != nil {
		return err
	}
	phpBin := "/usr/bin/php"
	phpBinSrc := fmt.Sprintf("/www/server/php/%s/bin/php", version)
	phpIze := "/usr/bin/phpize"
	phpIzeSrc := fmt.Sprintf("/www/server/php/%s/bin/phpize", version)
	phpFpm := "/usr/bin/php-fpm"
	phpFpmSrc := fmt.Sprintf("/www/server/php/%s/sbin/php-fpm", version)
	phpPecl := "/usr/bin/pecl"
	phpPeclSrc := fmt.Sprintf("/www/server/php/%s/bin/pecl", version)
	phpPear := "/usr/bin/pear"
	phpPearSrc := fmt.Sprintf("/www/server/php/%s/bin/pear", version)
	phpCliIni := "/etc/php-cli.ini"
	phpCliIniSrc := fmt.Sprintf("/www/server/php/%s/etc/php-cli.ini", version)
	IsChattr := false
	if result, _ := util.ExecShell("lsattr /usr|grep /usr/bin"); strings.Contains(result, "-i-") {
		IsChattr = true
		_, _ = util.ExecShell("chattr -i /usr/bin")
	}
	_, _ = util.ExecShell(fmt.Sprintf("rm -f %s %s %s %s %s %s ", phpBin, phpIze, phpFpm, phpPecl, phpPear, phpCliIni))
	_, _ = util.ExecShell(fmt.Sprintf("ln -sf %s %s", phpBinSrc, phpBin))
	_, _ = util.ExecShell(fmt.Sprintf("ln -sf %s %s", phpIzeSrc, phpIze))
	_, _ = util.ExecShell(fmt.Sprintf("ln -sf %s %s", phpFpmSrc, phpFpm))
	_, _ = util.ExecShell(fmt.Sprintf("ln -sf %s %s", phpPeclSrc, phpPecl))
	_, _ = util.ExecShell(fmt.Sprintf("ln -sf %s %s", phpPearSrc, phpPear))
	_, _ = util.ExecShell(fmt.Sprintf("ln -sf %s %s", phpCliIniSrc, phpCliIni))
	if IsChattr {
		_, _ = util.ExecShell("chattr +i /usr/bin")
	}
	return nil
}

// GetBinPath 取PHP的bin路径
func (s *ExtensionPHPService) GetBinPath(version string) string {
	return global.Config.System.ServerPath + "/php/" + version + "/bin/php"
}

// GetShellPath 取PHP安装shell脚本路径
func (s *ExtensionPHPService) GetShellPath() string {
	return global.Config.System.PanelPath + "/data/extensions/php/install"
}

// GetLibShellPath 取PHP安装LibShell脚本路径
func (s *ExtensionPHPService) GetLibShellPath() string {
	return global.Config.System.PanelPath + "/data/extensions/php/install"
}

// 读取存储PHPLib的json文件
func (s *ExtensionPHPService) readPHPLibJson() ([]*response.LibP, error) {
	// 读取JSON文件
	data, err := os.ReadFile(global.Config.System.PanelPath + "/data/extensions/php/Lib.json")
	if err != nil {
		return nil, err
	}

	// 解析json字符串
	var libItems []*response.LibP
	err = json.Unmarshal(data, &libItems)
	if err != nil {
		return nil, err
	}

	return libItems, nil
}

// getFpmAddress 获取FPM请求地址
//
// 参数：
//
//	version - php版本号
//	bind - 是否绑定到本地
//
// 返回值：
//
//	unix:unix套接字，bind:绑定地址，port:端口，
func (s *ExtensionPHPService) getFpmAddress(version string, bind bool) (string, string, int, error) {
	fpmAddress := "/tmp/php-cgi-" + version + ".sock"
	fpmConf, err := s.readPHPFpmConf(version)
	if err != nil {
		return "", "", 0, err
	}
	tmp := regexp.MustCompile(`listen\s*=\s*(.+)`).FindAllStringSubmatch(fpmConf, -1)
	if len(tmp) == 0 {
		return "unix", fpmAddress, 0, nil
	}
	listen := tmp[0][1]
	if strings.Contains(listen, "sock") {
		return "unix", fpmAddress, 0, nil
	}
	if strings.Contains(listen, ":") {
		listenTmp := strings.Split(listen, ":")
		port, _ := strconv.Atoi(listenTmp[1])
		if bind {
			return "tcp", listenTmp[0], port, nil
		} else {
			return "tcp", "127.0.0.1", port, nil
		}
	} else {
		port, _ := strconv.Atoi(listen)
		return "tcp", "127.0.0.1", port, nil
	}
}

// 读取php-fpm.conf配置文件
func (s *ExtensionPHPService) readPHPFpmConf(version string) (string, error) {
	content, err := util.ReadFileStringBody(global.Config.System.ServerPath + "/php/" + version + "/etc/php-fpm.conf")
	if err != nil {
		return "", err
	}
	return content, nil
}

// 写入php-fpm.conf配置文件
func (s *ExtensionPHPService) writePHPFpmConf(version, body string) error {
	err := util.WriteFile(global.Config.System.ServerPath+"/php/"+version+"/etc/php-fpm.conf", []byte(body), 644)
	if err != nil {
		return err
	}
	return nil
}

// readPHPIni 读取php.ini配置文件
func (s *ExtensionPHPService) readPHPIni(version string) (string, error) {
	content, err := util.ReadFileStringBody(global.Config.System.ServerPath + "/php/" + version + "/etc/php.ini")
	if err != nil {
		return "", err
	}
	return content, nil
}

// savePHPIni 保存php.ini配置文件
func (s *ExtensionPHPService) savePHPIni(version, body string) error {
	err := util.WriteFile(global.Config.System.ServerPath+"/php/"+version+"/etc/php.ini", []byte(body), 644)
	if err != nil {
		return err
	}
	return nil
}

// getLibByName 根据名称获取扩展信息
func (s *ExtensionPHPService) getLibByName(name string) (*response.LibP, error) {
	libJson, err := s.readPHPLibJson()
	if err != nil {
		return nil, err
	}
	var installLib *response.LibP
	for _, lib := range libJson {
		if lib.Name == name {
			installLib = lib
		}
	}
	return installLib, nil
}

func (s *ExtensionPHPService) requestPhpFpm(documentRoot string, uri string, method string, params map[string]string, addr string) ([]byte, error) {
	parts := strings.Split(uri, "?")
	scriptName, queryString := "", ""
	if len(parts) > 1 {
		scriptName = parts[0]
		queryString = parts[1]
	} else {
		scriptName = uri
	}

	if !strings.HasSuffix(documentRoot, "/") {
		documentRoot += "/"
	}
	fmt.Println(documentRoot, scriptName, queryString, method, params, addr)
	fmt.Println(fmt.Sprintf("%s%s", documentRoot, scriptName))
	env := map[string]string{
		"SCRIPT_FILENAME": fmt.Sprintf("%s%s", documentRoot, scriptName),
		"SCRIPT_NAME":     scriptName,
		"QUERY_STRING":    queryString,
		"REQUEST_METHOD":  method,
	}

	fcgi, err := fcgiClient.Dial("unix", addr)
	if err != nil {
		return nil, err
	}

	defer fcgi.Close()

	resp, err := fcgi.Get(env)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()

	}(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 0 {
		return nil, errors.New("php-fpm status code: " + strconv.Itoa(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// CheckInstall 检查是否已安装该版本的php
func (s *ExtensionPHPService) CheckInstall(version string) error {
	if !util.PathExists(s.GetBinPath(version)) {
		return errors.New(helper.Message("extension.NotInstalled"))
	}
	return nil
}

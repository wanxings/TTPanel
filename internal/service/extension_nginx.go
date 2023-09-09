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
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ExtensionNginxService struct{}

// Info 获取Nginx详细信息
func (s *ExtensionNginxService) Info() (*response2.ExtensionsInfoResponse, error) {
	var nginxInfo response2.ExtensionsInfoResponse
	err := ReadExtensionsInfo(constant.ExtensionNginxName, &nginxInfo)
	if err != nil {
		return nil, err
	}

	//判断是否安装
	if version, ok := s.IsInstalled(); ok {
		nginxInfo.Description.Version = version
		nginxInfo.Description.Install = true
		//获取运行状态
		nginxInfo.Description.Status = s.IsRunning()
	} else {
		nginxInfo.Description.Version = ""
		nginxInfo.Description.Install = false
		nginxInfo.Description.Status = false
	}

	return &nginxInfo, nil
}

// Install 安装Nginx
func (s *ExtensionNginxService) Install(version string) error {
	//检查是否在等待或者进行队列中
	taskName := "安装[Nginx-" + version + "]"
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

// Uninstall 卸载Nginx
func (s *ExtensionNginxService) Uninstall(version string) error {
	//检查是否在等待或者进行队列中
	taskName := "卸载[Nginx-" + version + "]"
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

// SetStatus 设置Nginx运行状态
func (s *ExtensionNginxService) SetStatus(action string) error {
	cmdStr := ""

	switch action {
	case constant.ProcessCommandByStart:
		//启动
		cmdStr = "/etc/init.d/nginx " + action
	case constant.ProcessCommandByStop:
		//关闭
		cmdStr = "/etc/init.d/nginx " + action
	case constant.ProcessCommandByRestart:
		//重启
		cmdStr = "/etc/init.d/nginx " + action
	case constant.ProcessCommandByReload:
		//重载配置
		cmdStr = "/etc/init.d/nginx " + action
	default:
		return errors.New("action is empty")
	}
	//执行命令
	//output, err := util.ExecShellWithResult("/bin/bash", "echo $("+cmdStr+")")
	output, err := util.ExecShell(cmdStr)
	if err != nil {
		return errors.New(err.Error() + output)
	}
	if !strings.Contains(output, "done") {
		return errors.New(output)
	}
	return nil
}

// PerformanceConfig 获取性能配置
func (s *ExtensionNginxService) PerformanceConfig() ([]*response2.NginxGeneralConfigP, error) {
	configBody, err := s.ReadConfig()
	if err != nil {
		return nil, err
	}
	proxyConfigBody, err := s.ReadProxyConfig()
	if err != nil {
		return nil, err
	}

	configData := make([]*response2.NginxGeneralConfigP, len(response2.DefaultNginxGeneralConfigs))
	copy(configData, response2.DefaultNginxGeneralConfigs)

	for _, v := range configData {
		rep := fmt.Sprintf("(%s)\\s+(\\w+)", v.Name)
		var matchList []string
		if v.From == "config" {
			matchList = regexp.MustCompile(rep).FindStringSubmatch(configBody)
		} else {
			matchList = regexp.MustCompile(rep).FindStringSubmatch(proxyConfigBody)
		}

		//如果没有匹配到第一个
		if len(matchList) < 2 {
			return nil, errors.New("key: " + v.Name + " is not found")
		}

		//如果没有匹配到第二个
		if len(matchList) < 3 {
			return nil, errors.New("value: " + v.Value + " is not found")
		}
		value := matchList[2]
		rep = "[kmgKMG]"
		//unit := regexp.MustCompile(rep).FindString(value)
		if matched, _ := regexp.MatchString(rep, value); matched {
			v.Unit = strings.ToUpper(string(value[len(value)-1]))
			v.Value = value[:len(value)-1]
			if len(v.Unit) == 1 {
				v.Ps = v.Unit + "B," + v.Ps
			} else {
				v.Ps = v.Unit + "," + v.Ps
			}
		} else {
			v.Unit = ""
			v.Value = value
		}
	}
	return configData, nil

}

// SavePerformanceConfig 保存性能配置
func (s *ExtensionNginxService) SavePerformanceConfig(param request.NginxSetPerformanceConfigR) error {
	configBody, err := s.ReadConfig()
	if err != nil {
		return err
	}
	backupConfigBody := configBody

	proxyConfigBody, err := s.ReadProxyConfig()
	if err != nil {
		return err
	}
	backupProxyConfigBody := proxyConfigBody

	// 遍历结构体字段
	t := reflect.TypeOf(param)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := reflect.ValueOf(param).Field(i)
		rep := regexp.MustCompile(fmt.Sprintf("%s\\s+[^kKmMgG\\;\\n]+", field.Tag.Get("json")))

		if rep.MatchString(configBody) {
			newConf := fmt.Sprintf("%s %s", field.Tag.Get("json"), value)
			configBody = rep.ReplaceAllString(configBody, newConf)
		} else if rep.MatchString(proxyConfigBody) {
			newConf := fmt.Sprintf("%s %s", field.Tag.Get("json"), value)
			proxyConfigBody = rep.ReplaceAllString(proxyConfigBody, newConf)
		}

	}

	//写入配置文件
	err = s.WriteConfig(configBody)
	if err != nil {
		return err
	}
	err = s.WriteProxyConfig(proxyConfigBody)
	if err != nil {
		return err
	}

	//检查Nginx配置文件是否有错误
	err = s.CheckConfig()
	if err != nil {
		//写入配置文件
		errs := s.WriteConfig(backupConfigBody)
		if errs != nil {
			return errs
		}
		errs = s.WriteProxyConfig(backupProxyConfigBody)
		if errs != nil {
			return errs
		}
		return err
	}
	err = s.SetStatus("reload")
	if err != nil {
		return err
	}
	return nil
}

// ReadConfig 读取配置文件
func (s *ExtensionNginxService) ReadConfig() (string, error) {
	configPath := fmt.Sprintf("%s/nginx/conf/nginx.conf", global.Config.System.ServerPath)
	file, err := util.ReadFileStringBody(configPath)
	if err != nil {
		return "", err
	}
	return file, nil
}

// WriteConfig 写入配置文件
func (s *ExtensionNginxService) WriteConfig(body string) error {
	configPath := fmt.Sprintf("%s/nginx/conf/nginx.conf", global.Config.System.ServerPath)
	err := util.WriteFile(configPath, []byte(body), 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadProxyConfig 读取代理配置文件
func (s *ExtensionNginxService) ReadProxyConfig() (string, error) {
	proxyConfigPath := fmt.Sprintf("%s/nginx/conf/proxy.conf", global.Config.System.ServerPath)
	file, err := util.ReadFileStringBody(proxyConfigPath)
	if err != nil {
		return "", err
	}
	return file, nil
}

// WriteProxyConfig 写入代理配置文件
func (s *ExtensionNginxService) WriteProxyConfig(body string) error {
	proxyConfigPath := fmt.Sprintf("%s/nginx/conf/proxy.conf", global.Config.System.ServerPath)
	err := util.WriteFile(proxyConfigPath, []byte(body), 0644)
	if err != nil {
		return err
	}
	return nil
}

// LoadStatus 获取Nginx的负载状态
func (s *ExtensionNginxService) LoadStatus() (*response2.NginxLoadStatusP, error) {
	//worker = int(public.ExecShell("ps aux|grep nginx|grep 'worker process'|wc -l")[0])-1
	//workermen = int(public.ExecShell("ps aux|grep nginx|grep 'worker process'|awk '{memsum+=$6};END {print memsum}'")[0]) / 1024
	statusData := &response2.NginxLoadStatusP{}
	var err error
	statusData.Worker, err = s.GetWorkerNum()
	if err != nil {
		return nil, err
	}
	statusData.WorkerMen, err = s.GetWorkerMenSize()
	if err != nil {
		return nil, err
	}
	statusData.WorkerCpu, err = s.GetWorkerCpu()
	if err != nil {
		return nil, err
	}

	err = s.CheckStatusConf()
	if err != nil {
		return nil, err
	}
	statusByHttp, err := s.GetNginxStatusByHttp()
	if err != nil {
		return nil, err
	}
	fmt.Println(statusByHttp)

	statusTmp := strings.Fields(statusByHttp)
	isRequestTime := false
	for _, s := range statusTmp {
		if s == "bar" {
			isRequestTime = true
			break
		}
	}
	if isRequestTime {
		statusData.Accepts = statusTmp[8]
		statusData.Handled = statusTmp[9]
		statusData.Requests = statusTmp[10]
		statusData.Reading = statusTmp[13]
		statusData.Writing = statusTmp[15]
		statusData.Waiting = statusTmp[17]
	} else {
		statusData.Accepts = statusTmp[9]
		statusData.Handled = statusTmp[7]
		statusData.Requests = statusTmp[8]
		statusData.Reading = statusTmp[11]
		statusData.Writing = statusTmp[13]
		statusData.Waiting = statusTmp[15]
	}
	statusData.Active = statusTmp[2]
	fmt.Println(statusTmp)
	return statusData, nil
}

// GetShellPath 获取shell安装文件路径
func (s *ExtensionNginxService) GetShellPath() string {
	return global.Config.System.PanelPath + "/data/extensions/nginx/install"
}

// 获取nginx安装路径

// GetNginxInfo 获取Nginx信息
func (s *ExtensionNginxService) GetNginxInfo(config *response2.ExtensionsInfoResponse) error {
	file, err := util.ReadFileStringBody(global.Config.System.PanelPath + "/data/extensions/nginx/info.json")
	if err != nil {
		return err
	}
	err = util.JsonStrToStruct(file, config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// IsInstalled
// Nginx是否已安装
// Return: 返回版本号和安装状态
func (s *ExtensionNginxService) IsInstalled() (string, bool) {
	//检查是否已安装
	if util.PathExists(s.GetBinPath()) {
		result, _ := util.ExecShell("/www/server/nginx/sbin/nginx -v 2>&1|grep version|awk '{print $3}'|cut -f2 -d'/'")
		//去除换行和空格
		result = util.ClearStr(result)
		return result, true
	} else {
		return "", false
	}
}

// GetBinPath 取Nginx的bin路径
func (s *ExtensionNginxService) GetBinPath() string {
	return global.Config.System.ServerPath + "/nginx/sbin/nginx"
}

func (s *ExtensionNginxService) IsRunning() bool {
	result, _ := util.ExecShell("/etc/init.d/nginx status")
	//去除换行和空格
	result = util.ClearStr(result)
	if strings.Contains(result, "running") {
		return true
	} else {
		return false
	}
}

func (s *ExtensionNginxService) CheckConfig() error {
	cmdStr := fmt.Sprintf("ulimit -n 8192 ;%s -t -c %s/nginx/conf/nginx.conf -v 2>&1", s.GetBinPath(), global.Config.System.ServerPath)
	result, err := util.ExecShell(cmdStr)
	if err != nil {
		return errors.New(err.Error() + " Error:" + result)
	}
	if strings.Contains(result, "successful") {
		return nil
	} else {
		return errors.New(result)
	}
}

// CheckStatusConf 检查状态配置文件
func (s *ExtensionNginxService) CheckStatusConf() error {
	filePath := global.Config.System.PanelPath + "/data/extensions/nginx/vhost/phpfpm_status.conf"
	if util.PathExists(filePath) {
		file, err := util.ReadFileStringBody(filePath)
		if err != nil {
			return err
		}
		if strings.Contains(file, "/nginx_status") {
			return nil
		} else {
			newConf := `server {
    listen 80;
    server_name 127.0.0.1;
    allow 127.0.0.1;
    location /nginx_status {
        stub_status on;
        access_log off;
    }
}`
			_ = util.WriteFile(filePath, []byte(newConf), 0644)
		}
		return nil
	}
	return nil
}

// GetWorkerNum 获取worker数量
func (s *ExtensionNginxService) GetWorkerNum() (int, error) {
	result, err := util.ExecShell("ps aux|grep nginx|grep 'worker process'|wc -l")
	if err != nil {
		return 0, err
	}
	result = util.ClearStr(result)
	newR, _ := strconv.Atoi(result)
	worker := int(newR) - 1
	return worker, nil
}

// GetWorkerMenSize 获取workerMen大小
func (s *ExtensionNginxService) GetWorkerMenSize() (int, error) {
	result, err := util.ExecShell("ps aux|grep nginx|grep 'worker process'|awk '{memsum+=$6};END {print memsum}'")
	if err != nil {
		return 0, err
	}
	result = util.ClearStr(result)
	newR, _ := strconv.Atoi(result)
	workerMen := newR / 1024
	return workerMen, nil
}

// GetWorkerCpu 获取workerCpu
func (s *ExtensionNginxService) GetWorkerCpu() (float64, error) {
	result, err := util.ExecShell("ps aux|grep nginx|grep 'worker process'|awk '{cpusum+=$3};END {print cpusum}'")
	if err != nil {
		return 0, err
	}
	result = util.ClearStr(result)
	newR, _ := strconv.Atoi(result)
	nginxCpu := float64(newR)
	return nginxCpu, nil
}

// GetNginxStatusByHttp http请求获取nginx状态内容
func (s *ExtensionNginxService) GetNginxStatusByHttp() (string, error) {
	url := "http://127.0.0.1/nginx_status"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

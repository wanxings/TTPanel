package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type ExtensionPhpmyadminService struct{}

// Info 获取Mysql详细信息
func (s *ExtensionPhpmyadminService) Info() (*response.ExtensionsInfoResponse, error) {
	var phpmyadminInfo response.ExtensionsInfoResponse
	err := ReadExtensionsInfo(constant.ExtensionPhpmyadminName, &phpmyadminInfo)
	if err != nil {
		return nil, err
	}

	//判断是否安装
	if version, ok := s.IsInstalled(); ok {
		phpmyadminInfo.Description.Version = version
		phpmyadminInfo.Description.Install = true
		config, _ := s.GetConfig()
		if status, ok := config["status"]; ok {
			phpmyadminInfo.Description.Status = status.(bool)
		}
	} else {
		phpmyadminInfo.Description.Version = ""
		phpmyadminInfo.Description.Install = false
		phpmyadminInfo.Description.Status = false
	}

	return &phpmyadminInfo, nil
}

// Install 安装phpmyadmin
func (s *ExtensionPhpmyadminService) Install(version string) error {
	//验证版本号是否正确
	if ok := util.IsMysqlVersion(version); !ok {
		return errors.New("version error")
	}
	//检查是否在等待或者进行队列中
	taskName := "安装[phpMyAdmin-" + version + "]"
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

// Uninstall 卸载phpmyadmin
func (s *ExtensionPhpmyadminService) Uninstall(version string) error {
	//验证版本号是否正确
	if ok := util.IsMysqlVersion(version); !ok {
		return errors.New("version error")
	}
	//检查是否在等待或者进行队列中
	taskName := "卸载[phpMyAdmin-" + version + "]"
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

func (s *ExtensionPhpmyadminService) IsInstalled() (string, bool) {
	//检查是否已安装
	if util.PathExists("/www/server/phpmyadmin/version.pl") {
		body, err := util.ReadFileStringBody("/www/server/phpmyadmin/version.pl")
		if err != nil {
			return "", false
		}
		return body, true
	} else {
		return "", false
	}
}

func (s *ExtensionPhpmyadminService) GetShellPath() string {
	return global.Config.System.PanelPath + "/data/extensions/phpmyadmin/install"
}

// GetPhpmyadminEntrance 获取phpmyadmin入口
func (s *ExtensionPhpmyadminService) GetPhpmyadminEntrance() (string, string) {
	phpmyadminPath := global.Config.System.ServerPath + "/phpmyadmin"
	if !util.PathExists(phpmyadminPath) {
		return "", ""
	}
	port := "888"
	entrance := "phpmyadmin"
	nginxConfigPath := global.Config.System.ServerPath + "/nginx/conf/nginx.conf"
	nginxConfigBody, err := util.ReadFileStringBody(nginxConfigPath)
	if err != nil {
		return "", ""
	}
	rep := regexp.MustCompile(`listen\s+([0-9]+)\s*;`)
	rtmp := rep.FindStringSubmatch(nginxConfigBody)
	if len(rtmp) > 1 {
		port = rtmp[1]
	}
	dirs, err := os.ReadDir(phpmyadminPath)
	if err != nil {
		return "", ""
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			if strings.HasPrefix(dir.Name(), "phpmyadmin") {
				entrance = dir.Name()
			}
		}
	}
	return entrance, port
}

// GetConfig 获取phpmyadmin配置
func (s *ExtensionPhpmyadminService) GetConfig() (map[string]interface{}, error) {
	nginxConfigPath := global.Config.System.ServerPath + "/nginx/conf/nginx.conf"
	phpConfigPath := global.Config.System.ServerPath + "/nginx/conf/enable-php.conf"
	auth := false
	status := false
	phpVersion := "54"
	if !util.PathExists(nginxConfigPath) {
		return nil, errors.New(fmt.Sprintf("not found file:%s", nginxConfigPath))
	}

	nginxConfigBody, err := util.ReadFileStringBody(nginxConfigPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("not found file:%s", nginxConfigPath))
	}

	if strings.Contains(nginxConfigBody, "AUTH_START") {
		auth = true
	}

	if !strings.Contains(nginxConfigBody, global.Config.System.ServerPath+"/stop") {
		status = true
	}

	if !util.PathExists(phpConfigPath) {
		err = util.CopyFile(global.Config.System.PanelPath+"/data/extensions/nginx/template/enable-php.conf", phpConfigPath)
		if err != nil {
			return nil, err
		}
	}
	phpConfigBody, err := util.ReadFileStringBody(phpConfigPath)
	if err != nil {
		return nil, err
	}
	rep := regexp.MustCompile(`php-cgi-([0-9]+)\.sock`)
	rtmp := rep.FindStringSubmatch(phpConfigBody)
	if len(rtmp) > 1 {
		phpVersion = rtmp[1]
	}

	url, port := s.GetPhpmyadminEntrance()

	return map[string]interface{}{
		"url":    fmt.Sprintf("http://%s:%s/%s", global.Config.System.PanelIP, port, url),
		"status": status,
		"auth":   auth,
		"port":   port,
		"php":    phpVersion,
	}, nil

}

// SetConfig 设置phpmyadmin配置
func (s *ExtensionPhpmyadminService) SetConfig(port, phpVersion string, auth, status bool) error {
	nginxConfigPath := global.Config.System.ServerPath + "/nginx/conf/nginx.conf"
	oldPort := "888"
	nginxConfigBody, err := util.ReadFileStringBody(nginxConfigPath)
	if err != nil {
		return errors.New(fmt.Sprintf("not found file:%s", nginxConfigPath))
	}
	global.Log.Debugf("ExtensionPhpmyadminService->SetConfig,port:%s,phpVersion:%s", port, phpVersion)
	if !util.StrIsEmpty(port) {
		//检查端口是否被占用
		if !util.CheckPort(port) {
			return errors.New(helper.MessageWithMap("PortOccupied", map[string]any{"Port": port}))
		}
		rep := regexp.MustCompile(`listen\s+([0-9]+)\s*;`)
		rtmp := rep.FindStringSubmatch(nginxConfigBody)
		if len(rtmp) > 1 {
			oldPort = rtmp[1]
		}
		if oldPort != port {
			nginxConfigBody = rep.ReplaceAllString(nginxConfigBody, "listen "+port+";")
			err = util.WriteFile(nginxConfigPath, []byte(nginxConfigBody), 0644)
			if err != nil {
				return err
			}
		}
	}
	if !util.StrIsEmpty(phpVersion) {
		phpConfigPath := global.Config.System.ServerPath + "/nginx/conf/enable-php.conf"
		if !util.PathExists(phpConfigPath) {
			err = util.CopyFile(global.Config.System.PanelPath+"/data/extensions/nginx/template/enable_php/enable-php.conf", phpConfigPath)
			if err != nil {
				return err
			}
		}
		phpConfigBody, err := util.ReadFileStringBody(phpConfigPath)
		if err != nil {
			return err
		}
		rep := regexp.MustCompile(`php-cgi-([0-9]+)\.sock`)
		phpConfigBody = rep.ReplaceAllString(phpConfigBody, "php-cgi-"+phpVersion+".sock")
		global.Log.Debugf("phpConfigBody:%s", phpConfigBody)
		err = util.WriteFile(phpConfigPath, []byte(phpConfigBody), 0644)
		if err != nil {
			return err
		}
	}
	err = GroupApp.ExtensionNginxServiceApp.SetStatus("reload")
	if err != nil {
		return err
	}
	return nil

}

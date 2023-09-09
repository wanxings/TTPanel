package service

import (
	"TTPanel/internal/global"
	authModel "TTPanel/internal/model"
	"TTPanel/pkg/util"
)

type SettingsService struct{}

// List 获取系统设置
func (s *SettingsService) List() (map[string]interface{}, error) {
	settingsData := make(map[string]interface{})
	settingsData["panel_name"] = global.Config.System.PanelName
	//面板端口
	settingsData["panel_port"] = global.Config.System.PanelPort
	//面板端口
	settingsData["panel_ip"] = global.Config.System.PanelIP
	//session超时时间
	settingsData["session_expire"] = global.Config.System.SessionExpire
	//安全入口
	settingsData["entrance"] = global.Config.System.Entrance
	//安全入口错误码
	settingsData["entrance_error_code"] = global.Config.System.EntranceErrorCode
	//BasicAuth认证
	settingsData["basic_auth"] = global.Config.System.BasicAuth
	//PanelApi
	settingsData["panel_api"] = global.Config.System.PanelApi
	//默认建站目录
	settingsData["default_website_directory"] = global.Config.System.DefaultProjectDirectory
	//默认备份目录
	settingsData["default_backup_directory"] = global.Config.System.DefaultBackupDirectory
	//语言
	settingsData["language"] = global.Config.System.Language
	//
	settingsData["server_date"], _ = util.ExecShell("date +\"%Y-%m-%d %H:%M:%S %Z %z\"")
	//面板用户
	user, err := (&authModel.User{ID: 1}).Get(global.PanelDB)
	if err != nil {
		return nil, err
	}
	settingsData["panel_user"] = user.Username
	return settingsData, nil
}

// SetBasicAuth 设置基础认证
func (s *SettingsService) SetBasicAuth(status bool, username, password string) error {
	global.Config.System.BasicAuth.Status = status
	global.Config.System.BasicAuth.Username = username
	global.Config.System.BasicAuth.Password = password
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}
	return nil
}

// SetPanelPort 设置面板端口
func (s *SettingsService) SetPanelPort(panelPort int) error {
	global.Config.System.PanelPort = panelPort
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}
	err = GroupApp.PanelServiceApp.OperatePanel("restart")
	if err != nil {
		return err
	}
	return nil
}

// SetEntrance 设置安全入口
func (s *SettingsService) SetEntrance(entrance string) error {
	global.Config.System.Entrance = entrance
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}
	err = GroupApp.PanelServiceApp.OperatePanel("restart")
	if err != nil {
		return err
	}
	return nil
}

// SetEntranceErrorCode 设置安全入口错误码
func (s *SettingsService) SetEntranceErrorCode(errorCode int) error {
	global.Config.System.EntranceErrorCode = errorCode
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}
	return nil
}

// SetUser 设置面板用户
func (s *SettingsService) SetUser(username, password string) error {
	user, err := (&authModel.User{ID: 1}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if !util.StrIsEmpty(username) {
		user.Username = username
	}
	if !util.StrIsEmpty(password) {
		user.Password, user.Salt = EncryptPasswordAndSalt(password)
	}
	err = user.Update(global.PanelDB)
	if err != nil {
		return err
	}
	adminToken := util.EncodeMD5(user.Username + user.Password)
	global.GoCache.Set("admin_token", adminToken, -1)
	return nil
}

// SetPanelName 设置面板名称
func (s *SettingsService) SetPanelName(panelName string) error {
	global.Config.System.PanelName = panelName
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}
	return nil
}

// SetPanelIP 设置面板IP
func (s *SettingsService) SetPanelIP(panelIP string) error {
	global.Config.System.PanelIP = panelIP
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}
	return nil
}

// SetDefaultWebsiteDirectory 设置默认建站目录
func (s *SettingsService) SetDefaultWebsiteDirectory(directory string) error {
	global.Config.System.DefaultProjectDirectory = directory
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}
	return nil
}

// SetDefaultBackupDirectory 设置默认备份目录
func (s *SettingsService) SetDefaultBackupDirectory(directory string) error {
	global.Config.System.DefaultBackupDirectory = directory
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}
	return nil
}

// SetPanelApi 设置面板API
func (s *SettingsService) SetPanelApi(status bool, key string, whitelist []string) error {
	global.Config.System.PanelApi.Status = status
	global.Config.System.PanelApi.Key = key
	global.Config.System.PanelApi.Whitelist = whitelist
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}
	return nil
}

// SetLanguage 设置面板语言
func (s *SettingsService) SetLanguage(language string) error {

	global.Config.System.Language = language
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}

	//切换语言
	// global.I18n = initialize.InitI18n()
	return nil
}

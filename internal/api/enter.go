package api

import "TTPanel/internal/service"

type Group struct {
	BackupApiApp              BackupApi
	CronTaskApiApp            CronTaskApi
	DatabaseMysqlApiApp       DatabaseMysqlApi
	ExplorerApiApp            ExplorerApi
	ExtensionDockerApiApp     ExtensionDockerApi
	ExtensionMysqlApiApp      ExtensionMysqlApi
	ExtensionNginxApiApp      ExtensionNginxApi
	ExtensionNodejsApiApp     ExtensionNodejsApi
	ExtensionPHPApiApp        ExtensionPHPApi
	ExtensionPhpmyadminApiApp ExtensionPhpmyadminApi
	ExtensionRedisApiApp      ExtensionRedisApi
	HostApiApp                HostApi
	LinuxToolsApiApp          LinuxToolsApi
	MonitorApiApp             MonitorApi
	SSHManageApiApp           SSHManageApi
	NotifyApiApp              NotifyApi
	PanelApiApp               PanelApi
	ProjectApiApp             ProjectApi
	ProjectGeneralApiApp      ProjectGeneralApi
	ProjectPHPApiApp          ProjectPHPApi
	QueueTaskApiApp           QueueTaskApi
	RecycleBinApiApp          RecycleBinApi
	SettingsApiApp            SettingsApi
	SSLApiApp                 SSLApi
	StorageApiApp             StorageApi
	SystemFirewallApiApp      SystemFirewallApi
	TaskManagerApiApp         TaskManagerApi
	TTWafApiApp               TTWafApi
	UserApiApp                UserApi
	WebSSHApiApp              WebSSHApi
	LogAuditApiApp            LogAuditApi
}

var GroupApp = new(Group)
var ServiceGroupApp = service.GroupApp

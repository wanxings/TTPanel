package service

type Group struct {
	UserServiceApp                UserService
	BackupServiceApp              BackupService
	CronTaskServiceApp            CronTaskService
	DatabaseMysqlServiceApp       DatabaseMysqlService
	ExplorerServiceApp            ExplorerService
	ExtensionDockerServiceApp     ExtensionDockerService
	ExtensionMysqlServiceApp      ExtensionMysqlService
	ExtensionNginxServiceApp      ExtensionNginxService
	ExtensionNodejsServiceApp     ExtensionNodejsService
	ExtensionPhpmyadminServiceApp ExtensionPhpmyadminService
	ExtensionPHPServiceApp        ExtensionPHPService
	ExtensionRedisServiceApp      ExtensionRedisService
	HostServiceApp                HostService
	MonitorServiceApp             MonitorService
	NotifyServiceApp              NotifyService
	SSHManageServiceApp           SSHManageService
	PanelServiceApp               PanelService
	ProjectServiceApp             ProjectService
	ProjectGeneralServiceApp      ProjectGeneralService
	ProjectPHPServiceApp          ProjectPHPService
	QueueTaskServiceApp           QueueTaskService
	RecycleBinServiceApp          RecycleBinService
	SettingsServiceApp            SettingsService
	SSLServiceApp                 SSLService
	StorageServiceApp             StorageService
	SystemFirewallServiceApp      SystemFirewallService
	TaskManagerServiceApp         TaskManagerService
	TTWafServiceApp               TTWafService
	WebSSHServiceApp              WebSSHService
	LinuxToolsServiceApp          LinuxToolsService
	LogAuditServiceApp            LogAuditService
}

var GroupApp = new(Group)

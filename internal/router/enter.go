package router

type Group struct {
	BackupRouterApp              BackupRouter
	CronTaskRouterApp            CronTaskRouter
	DatabaseMysqlRouterApp       DatabaseMysqlRouter
	ExplorerRouterApp            ExplorerRouter
	ExtensionDockerRouterApp     ExtensionDockerRouter
	ExtensionMysqlRouterApp      ExtensionMysqlRouter
	ExtensionNginxRouterApp      ExtensionNginxRouter
	ExtensionNodejsRouterApp     ExtensionNodejsRouter
	ExtensionPHPRouterApp        ExtensionPHPRouter
	ExtensionPhpmyadminRouterApp ExtensionPhpmyadminRouter
	ExtensionRedisRouterApp      ExtensionRedisRouter
	HostRouterApp                HostRouter
	MonitorRouterApp             MonitorRouter
	NotifyRouterApp              NotifyRouter
	PanelRouterApp               PanelRouter
	ProjectRouterApp             ProjectRouter
	ProjectGeneralRouterApp      ProjectGeneralRouter
	ProjectPHPRouterApp          ProjectPHPRouter
	QueueTaskRouterApp           QueueTaskRouter
	RecycleBinRouterApp          RecycleBinRouter
	SettingsRouterApp            SettingsRouter
	SSLRouterApp                 SSLRouter
	StorageRouterApp             StorageRouter
	SystemFirewallRouterApp      SystemFirewallRouter
	TaskManagerRouterApp         TaskManagerRouter
	TTWafRouterApp               TTWafRouter
	UserRouterApp                UserRouter
	UserBaseRouterApp            UserBaseRouter
	WebSSHRouterApp              WebSSHRouter
	StaticRouterApp              StaticRouter
	LinuxToolsRouterApp          LinuxToolsRouter
	SSHManageRouterApp           SSHManageRouter
	LogAuditRouterApp            LogAuditRouter
}

var GroupApp = new(Group)

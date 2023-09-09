package constant

const (
	ginModeByDebug   = "debug"
	ginModeByRelease = "release"
)
const (
	AuthModeByJwt            = "jwt"
	AuthModeByApiToken       = "apiToken"
	AuthModeBySession        = "session"
	AuthModeByTemporaryToken = "temporaryToken"
)
const (
	DataBaseTypeByMysql = "mysql"
)

const (
	ProjectRedirectTypeDomain = iota + 1
	ProjectRedirectTypePath
)
const (
	ProjectAccessRuleTypeBasicAuth = iota + 1
	ProjectAccessRuleTypeNoAccess
)
const (
	SystemFirewallStrategyAllow = iota //策略:允许
	SystemFirewallStrategyDeny  = iota //策略:拒绝
)
const (
	SystemFirewallProtocolTCP int = iota
	SystemFirewallProtocolUDP int = iota
	SystemFirewallProtocolTCPANDUDP
)
const (
	ProcessCommandByStart   = "start"
	ProcessCommandByStop    = "stop"
	ProcessCommandByRestart = "restart"
	ProcessCommandByReload  = "reload"
)
const (
	OperateTTWafDomainConfigByAdd    = "add"
	OperateTTWafDomainConfigByUpdate = "update"
	OperateTTWafDomainConfigByDelete = "delete"
)
const (
	ExtensionRedisName      = "redis"
	ExtensionNodejsName     = "nodejs"
	ExtensionNginxName      = "nginx"
	ExtensionPHPName        = "php"
	ExtensionMysqlName      = "mysql"
	ExtensionDockerName     = "docker"
	ExtensionPhpmyadminName = "phpmyadmin"
)
const (
	CompressTypeByZip   = "zip"
	CompressTypeByTar   = "tar"
	CompressTypeByTarGz = "tar.gz"
	CompressTypeByTarXz = "tar.xz"
	CompressTypeByXz    = "xz"
	CompressTypeByGz    = "gz"
	CompressTypeBy7z    = "7z"
)
const (
	OperationLogTypeByProjectManager    = iota + 1
	OperationLogTypeBySystemFirewall    = iota + 1
	OperationLogTypeBySystemTaskManager = iota + 1
	OperationLogTypeByLogAudit          = iota + 1
	OperationLogTypeByDatabase          = iota + 1
	OperationLogTypeByCronTask          = iota + 1
	OperationLogTypeByExplorer          = iota + 1
	OperationLogTypeByQueueTask         = iota + 1
	OperationLogTypeByUserLogin         = iota + 1
	OperationLogTypeByDockerManager     = iota + 1
	OperationLogTypeByHostManager       = iota + 1
	OperationLogTypeByBackup            = iota + 1
	OperationLogTypeByExtension         = iota + 1
	OperationLogTypeByLinuxTools        = iota + 1
	OperationLogTypeByMonitor           = iota + 1
	OperationLogTypeByNotify            = iota + 1
	OperationLogTypeByPanel             = iota + 1
	OperationLogTypeByRecycleBin        = iota + 1
	OperationLogTypeBySettings          = iota + 1
	OperationLogTypeBySystem            = iota + 1
	OperationLogTypeByStorage           = iota + 1
	OperationLogTypeByTaskManager       = iota + 1
	OperationLogTypeByTTWaf             = iota + 1
)
const (
	BackupCategoryByDatabase = iota + 1
	BackupCategoryByProject  = iota + 1
	BackupCategoryByPanel    = iota + 1
	BackupCategoryByDir      = iota + 1
)
const (
	ProjectTypeByPHP     = iota + 1
	ProjectTypeByGeneral = iota + 1
	ProjectTypeByNode    = iota + 1
	ProjectTypeByJava    = iota + 1
)
const (
	ProjectStatusByStop = iota
	ProjectStatusByRunning
	ProjectStatusByDel
)
const (
	QueueTaskStatusWait       int = iota
	QueueTaskStatusProcessing int = iota
	QueueTaskStatusSuccess    int = iota
	QueueTaskStatusError
)

const (
	ChattrCmd = "chattr"
	LsattrCmd = "lsattr"
)

const (
	SaveFileTypeFieldByNormal = iota
	SaveFileTypeFieldByNginxConf
)
const (
	NotifyCategoryByEmail      = iota + 1
	NotifyCategoryByDingTalk   = iota + 1
	NotifyCategoryByWeChatWork = iota + 1
	NotifyCategoryByAliSms     = iota + 1
	NotifyCategoryByTelegram   = iota + 1
)
const (
	NotifyLevelInfo    = "info"
	NotifyLevelWarning = "warning"
	NotifyLevelSuccess = "success"
	NotifyLevelDebug   = "debug"
)
const (
	SSLVerifyModeByDNSAccount = iota + 1
	SSLVerifyModeByManual     = iota + 1
	SslVerifyModeByFile
)
const (
	DnsTypeCloudflare   = "cloudflare"
	DnsTypeAliyun       = "alidns"
	DnsTypeTencentCloud = "tencentcloud"
)
const (
	StorageCategoryByS3         = iota + 1
	StorageCategoryByTencentCOS = iota + 1
	StorageCategoryByAliOSS     = iota + 1
	StorageCategoryByQiniuKodo  = iota + 1
	StorageCategoryByMinio      = iota + 1
)

// stable稳定版本 alpha开发版本 beta测试版本
const (
	PreReleaseVersionByStable = "stable"
	PreReleaseVersionByAlpha  = "alpha"
	PreReleaseVersionByBeta   = "beta"
)

const (
	CronTaskCategoryByShell          = iota + 1
	CronTaskCategoryByBackupProject  = iota + 1
	CronTaskCategoryByBackupDatabase = iota + 1
	CronTaskCategoryByCutLog         = iota + 1
	CronTaskCategoryByBackupDir      = iota + 1
	CronTaskCategoryByRequestUrl     = iota + 1
	CronTaskCategoryByFreeMemory     = iota + 1
)

const (
	AbnormalMonitoringCategoryByCertificateExpires = iota + 1
	AbnormalMonitoringCategoryBySiteProjectExpires = iota + 1
	AbnormalMonitoringCategoryByPanelLogin         = iota + 1
	AbnormalMonitoringCategoryBySSHLogin           = iota + 1
	AbnormalMonitoringCategoryByService            = iota + 1
)

const (
	SSHRootLoginTypeByYes                = "yes"
	SSHRootLoginTypeByNo                 = "no"
	SSHRootLoginTypeByWithoutPassword    = "without-password"
	SSHRootLoginTypeByForcedCommandsOnly = "forced-commands-only"
)

const (
	SSHKeyLoginTypeByEd25519 = "ed25519"
	SSHKeyLoginTypeByEcdsa   = "ecdsa"
	SSHKeyLoginTypeByRsa     = "rsa"
	SSHKeyLoginTypeByDsa     = "dsa"
)

const (
	MonitorEventCategoryByCPU                   = "cpu"
	MonitorEventCategoryByMemory                = "memory"
	MonitorEventCategoryByDiskSpace             = "disk_space"
	MonitorEventCategoryByDiskInode             = "disk_inode"
	MonitorEventCategoryByLoginPanel            = "login_panel"
	MonitorEventCategoryBySslExpirationTime     = "ssl_expiration_time"
	MonitorEventCategoryByProjectExpirationTime = "project_expiration_time"
	MonitorEventCategoryByServiceNginx          = "service_nginx"
	MonitorEventCategoryByServiceMysql          = "service_mysql"
	MonitorEventCategoryByServiceRedis          = "service_redis"
	MonitorEventCategoryByServiceDocker         = "service_docker"
)

const (
	SSHLoginStatusBySuccess = iota + 1
	SSHLoginStatusByFailed  = iota + 1
)

package conf

type Server struct {
	Sqlite  Sqlite  `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Logger  Logger  `mapstructure:"logger" json:"logger" yaml:"logger"`
	Monitor Monitor `mapstructure:"monitor" json:"monitor" yaml:"monitor"`
}
type Monitor struct {
	Status  bool `mapstructure:"status" json:"status" yaml:"status"`
	SaveDay int  `mapstructure:"save_day" json:"saveDay" yaml:"save_day"`
}
type Logger struct {
	RouterLog       bool   `mapstructure:"router_log" json:"router_log" yaml:"router_log"`
	RouterLogParams bool   `mapstructure:"router_log_params" json:"router_log_params" yaml:"router_log_params"`
	RootPath        string `mapstructure:"root_path" json:"root_path" yaml:"root_path"`
	FileExt         string `mapstructure:"file_ext" json:"file_ext" yaml:"file_ext"`
	LogLevel        string `mapstructure:"log_level" json:"logLevel" yaml:"log_level"`
}
type System struct {
	PanelIP                 string      `mapstructure:"panel_ip" json:"panel_ip" yaml:"panel_ip"`                                  // PanelIP
	PanelName               string      `mapstructure:"panel_name" json:"panel_name" yaml:"panel_name"`                            // 面板名称
	RunMode                 string      `mapstructure:"run_mode" json:"run_mode" yaml:"run_mode"`                                  // 环境值
	PanelPort               int         `mapstructure:"panel_port" json:"panel_port" yaml:"panel_port"`                            // 面板端口
	PreReleaseVersion       string      `mapstructure:"pre_release_version" json:"pre_release_version" yaml:"pre_release_version"` // 预发布版本
	SessionSecret           string      `mapstructure:"session_secret" json:"session_secret" yaml:"session_secret"`                // Session secret
	SessionExpire           int         `mapstructure:"session_expire" json:"session_expire" yaml:"session_expire"`                // jwt_expire 	// aes_key
	JwtSecret               string      `mapstructure:"jwt_secret" json:"jwt_secret" yaml:"jwt_secret"`                            // jwt_secret
	JwtExpire               int         `mapstructure:"jwt_expire" json:"jwt_expire" yaml:"jwt_expire"`                            // jwt_expire
	JwtIssuer               string      `mapstructure:"jwt_issuer" json:"jwt_issuer" yaml:"jwt_issuer"`                            // jwt_issuer
	Domain                  string      `mapstructure:"domain" json:"domain" yaml:"domain"`                                        // 域名
	Entrance                string      `mapstructure:"entrance" json:"entrance" yaml:"entrance"`                                  // 入口
	PanelPath               string      `mapstructure:"panel_path" json:"panel_path" yaml:"panel_path"`                            // 面板路径
	ServerPath              string      `mapstructure:"server_path" json:"server_path" yaml:"server_path"`                         // 环境服务路径
	WwwLogPath              string      `mapstructure:"www_log_path" json:"www_log_path" yaml:"www_log_path"`                      // 网站日志路径
	PluginPath              string      `mapstructure:"plugin_path" json:"plugin_path" yaml:"plugin_path"`                         // 插件路径 	// 分页大小
	Language                string      `mapstructure:"language" json:"language" yaml:"language"`                                  // 语言
	BasicAuth               BasicAuth   `mapstructure:"basic_auth" json:"basic_auth" yaml:"basic_auth"`                            // 基本认证
	RecycleBin              RecycleBin  `mapstructure:"recycle_bin" json:"recycle_bin" yaml:"recycle_bin"`
	FileHistory             FileHistory `mapstructure:"file_history" json:"file_history" yaml:"file_history"`
	MysqlRootPassword       string      `mapstructure:"mysql_root_password" json:"mysql_root_password" yaml:"mysql_root_password"`
	CloudNodes              []string    `mapstructure:"cloud_nodes" json:"cloud_nodes" yaml:"cloud_nodes"`                                           // 云节点
	DefaultProjectDirectory string      `mapstructure:"default_project_directory" json:"default_project_directory" yaml:"default_project_directory"` // 默认网站目录
	DefaultBackupDirectory  string      `mapstructure:"default_backup_directory" json:"default_backup_directory" yaml:"default_backup_directory"`    // 默认备份目录
	EntranceErrorCode       int         `mapstructure:"entrance_error_code" json:"entrance_error_code" yaml:"entrance_error_code"`                   // 入口错误码
	PanelApi                PanelApi    `mapstructure:"panel_api" json:"panel_api" yaml:"panel_api"`
	AutoCheckUpdate         bool        `mapstructure:"auto_check_update" json:"auto_check_update" yaml:"auto_check_update"`
}

type BasicAuth struct {
	Status   bool   `mapstructure:"status" json:"status" yaml:"status"`       // 状态
	Username string `mapstructure:"username" json:"username" yaml:"username"` // 用户名
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
}

type PanelApi struct {
	Status    bool     `mapstructure:"status" json:"status" yaml:"status"`          // 状态
	Key       string   `mapstructure:"key" json:"key" yaml:"key"`                   // 密钥
	Whitelist []string `mapstructure:"whitelist" json:"whitelist" yaml:"whitelist"` // 白名单
}

type RecycleBin struct {
	ExplorerStatus bool   `mapstructure:"explorer_status" json:"explorer_status" yaml:"explorer_status"` // 状态
	DatabaseStatus bool   `mapstructure:"database_status" json:"database_status" yaml:"database_status"` // 状态
	Directory      string `mapstructure:"directory" json:"directory" yaml:"directory"`                   // 回收站目录
}

type FileHistory struct {
	Status bool `mapstructure:"status" json:"status" yaml:"status"` // 状态
	Count  int  `mapstructure:"count" json:"count" yaml:"count"`    // 保存数量
}

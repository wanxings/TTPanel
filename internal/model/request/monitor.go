package request

type LogsR struct {
	Type      string `form:"type" json:"type" binding:"required"`
	StartTime int64  `form:"start_time" json:"start_time" binding:"required"`
	EndTime   int64  `form:"end_time" json:"end_time" binding:"required"`
}

type MonitorConfigR struct {
	Status  bool `form:"status" json:"status"`
	SaveDay int  `form:"saveDay" json:"saveDay" binding:"required"`
	MaxCpu  int  `form:"maxCpu" json:"maxCpu" binding:"required"`
	MaxMem  int  `form:"maxMem" json:"maxMem" binding:"required"`
	MaxNet  int  `form:"maxNet" json:"maxNet" binding:"required"`
}

type EventListR struct {
	Limit  int `form:"limit" json:"limit" binding:"required"`
	Status int `form:"status" json:"status"`
	Page   int `form:"page" json:"page" binding:"required"`
}

type BatchSetEventStatusR struct {
	Ids    []int64 `form:"ids" json:"ids" binding:"required"`
	Status int     `form:"status" json:"status"`
}
type MonitorEventConfig struct {
	Cpu struct {
		Status      bool    `json:"status"`
		NotifyCount int64   `json:"notify_count"`
		NotifyId    int64   `json:"notify_id"`
		Threshold   float64 `json:"threshold"`
	} `json:"cpu" binding:"required"`
	Mem struct {
		Status      bool    `json:"status"`
		NotifyCount int64   `json:"notify_count"`
		NotifyId    int64   `json:"notify_id"`
		Threshold   float64 `json:"threshold"`
	} `json:"mem" binding:"required"`
	DiskSpace struct {
		Status      bool     `json:"status"`
		NotifyCount int64    `json:"notify_count"`
		NotifyId    int64    `json:"notify_id"`
		Threshold   float64  `json:"threshold"`
		MountPoint  []string `json:"mount_point"`
	} `json:"disk_space" binding:"required"`
	DiskInode struct {
		Status      bool     `json:"status"`
		NotifyCount int64    `json:"notify_count"`
		NotifyId    int64    `json:"notify_id"`
		Threshold   float64  `json:"threshold"`
		MountPoint  []string `json:"mount_point"`
	} `json:"disk_inode" binding:"required"`
	LoginPanel struct {
		Status      bool  `json:"status"`
		NotifyCount int64 `json:"notify_count"`
		NotifyId    int64 `json:"notify_id"`
	} `json:"login_panel" binding:"required"`
	SslExpirationTime struct {
		Status      bool  `json:"status"`
		NotifyCount int64 `json:"notify_count"`
		NotifyId    int64 `json:"notify_id"`
		Day         int   `json:"day"`
	} `json:"ssl_expiration_time" binding:"required"`
	ProjectExpirationTime struct {
		Status      bool  `json:"status"`
		NotifyCount int64 `json:"notify_count"`
		NotifyId    int64 `json:"notify_id"`
		Day         int   `json:"day"`
	} `json:"project_expiration_time" binding:"required"`
	ServiceNginx struct {
		Status      bool  `json:"status"`
		NotifyCount int64 `json:"notify_count"`
		NotifyId    int64 `json:"notify_id"`
	} `json:"service_nginx" binding:"required"`
	ServiceMysql struct {
		Status      bool  `json:"status"`
		NotifyCount int64 `json:"notify_count"`
		NotifyId    int64 `json:"notify_id"`
	} `json:"service_mysql" binding:"required"`
	ServiceRedis struct {
		Status      bool  `json:"status"`
		NotifyCount int64 `json:"notify_count"`
		NotifyId    int64 `json:"notify_id"`
	} `json:"service_redis" binding:"required"`
	ServiceDocker struct {
		Status      bool  `json:"status"`
		NotifyCount int64 `json:"notify_count"`
		NotifyId    int64 `json:"notify_id"`
	} `json:"service_docker" binding:"required"`
}

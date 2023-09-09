package request

import (
	"TTPanel/internal/model"
)

type BatchCreateCronTaskR struct {
	List []*CronTaskInfo `json:"list"`
}

type EditR struct {
	Id int64 `json:"id" form:"id" binding:"required"`
	*CronTaskInfo
}
type CronTaskInfo struct {
	Name          string `json:"name"`
	Category      int    `json:"category"`
	TimeType      string `json:"time_type"`
	TimeInc       int    `json:"time_inc"`
	TimeHour      int    `json:"time_hour"`
	TimeMinute    int    `json:"time_minute"`
	TimeCustomize string `json:"time_customize"`
	ShellBody     string `json:"shell_body"`
	CronTaskConfig
}

type CronTaskConfig struct {
	BackupProjectConfig  model.BackupProjectConfig  `json:"backup_project_config"`
	BackupDatabaseConfig model.BackupDatabaseConfig `json:"backup_database_config"`
	CutLogConfig         model.CutLogConfig         `json:"cut_log_config"`
	BackupDirConfig      model.BackupDirConfig      `json:"backup_dir_config"`
	RequestUrlConfig     model.RequestUrlConfig     `json:"request_url_config"`
}
type CronCycle struct {
	Type      string `json:"type"`
	Minute    int    `json:"minute"`
	Hour      int    `json:"hour"`
	Inc       int    `json:"inc"`
	Customize string `json:"customize"`
}

type TaskListR struct {
	QueryName string `json:"query_name" form:"query_name"`
	SType     string `json:"sType" form:"sType"`
}
type TaskIdR struct {
	Id int64 `json:"id" form:"id" binding:"required"`
}

type BatchSetStatusR struct {
	IDs    []int64 `json:"ids" binding:"required"`
	Status int     `json:"status" form:"status"`
}
type GetLogR struct {
	Id   int64 `json:"id" form:"id" binding:"required"`
	Line int   `json:"line" form:"line"`
}

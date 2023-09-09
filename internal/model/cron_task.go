package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type CronTask struct {
	ID            int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Category      int    `gorm:"column:category;not null" json:"category" form:"category"`
	Hash          string `gorm:"column:hash;not null" json:"hash" form:"hash"`
	Name          string `gorm:"column:name;not null;unique" json:"name" form:"name"`
	TimeType      string `gorm:"column:time_type;not null" json:"time_type" form:"time_type"`
	TimeInc       int    `gorm:"column:time_inc;not null" json:"time_inc" form:"time_inc"`
	TimeHour      int    `gorm:"column:time_hour;not null" json:"time_hour" form:"time_hour"`
	TimeMinute    int    `gorm:"column:time_minute;not null" json:"time_minute" form:"time_minute"`
	TimeCustomize string `gorm:"column:time_customize;not null" json:"time_customize" form:"time_customize"`
	ShellBody     string `gorm:"column:shell_body;not null" json:"shell_body" form:"shell_body"`
	Config        string `gorm:"column:config;not null" json:"config" form:"config"`
	Status        int    `gorm:"column:status;not null" json:"status" form:"status"`
	LastRunTime   int64  `gorm:"-:all" json:"last_run_time" form:"last_run_time"`
	CreateTime    int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

type BackupProjectConfig struct {
	Id             int64    `json:"id"`
	StorageId      int64    `json:"storage_id"`
	KeepLocalFile  int      `json:"keep_local_file"`
	Number         int      `json:"number"`
	NotifyType     int      `json:"notify_type"`
	NotifyId       int64    `json:"notify_id"`
	ExclusionRules []string `json:"exclusion_rules"`
}

type BackupDatabaseConfig struct {
	Id            int64 `json:"id"`
	StorageId     int64 `json:"storage_id"`
	KeepLocalFile int   `json:"keep_local_file"`
	Number        int   `json:"number"`
	NotifyType    int   `json:"notify_type"`
	NotifyId      int64 `json:"notify_id"`
}

type CutLogConfig struct {
	Id     int64 `json:"id"`
	Number int   `json:"number"`
}

type BackupDirConfig struct {
	Path           string   `json:"path"`
	StorageId      int64    `json:"storage_id"`
	KeepLocalFile  int      `json:"keep_local_file"`
	Number         int      `json:"number"`
	NotifyType     int      `json:"notify_type"`
	NotifyId       int      `json:"notify_id"`
	ExclusionRules []string `json:"exclusion_rules"`
}

type RequestUrlConfig struct {
	Url        string `json:"url"`
	Keyword    string `json:"keyword"`
	Reg        string `json:"reg"`
	StatusCode int    `json:"status_code"`
	Timeout    int    `json:"timeout"`
	NotifyType []int  `json:"notify_type"`
	NotifyId   int64  `json:"notify_id"`
}

func (*CronTask) TableName() string {
	return "cron_task"
}
func (s *CronTask) Get(db *gorm.DB) (*CronTask, error) {
	var scheduledTask CronTask
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else if !util.StrIsEmpty(s.Hash) {
		db = db.Where("hash", s.Hash)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&scheduledTask).Error
	if err != nil {
		return &scheduledTask, err
	}

	return &scheduledTask, nil
}
func (s *CronTask) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*CronTask, int64, error) {
	var scheduledTask []*CronTask
	var err error
	var count int64

	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	if err = db.Model(s).Count(&count).Error; err != nil {
		return nil, count, err
	}

	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if err = db.Find(&scheduledTask).Error; err != nil {
		return nil, count, err
	}

	return scheduledTask, count, nil
}
func (s *CronTask) Create(db *gorm.DB) (*CronTask, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *CronTask) Update(db *gorm.DB) error {
	return db.Model(&CronTask{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *CronTask) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&CronTask{}, s.ID).Error
}
func (s *CronTask) Count(db *gorm.DB) (int64, error) {
	var count int64
	if s.Status > 0 {
		db = db.Where("status= ? ", s.Status)
	}
	db.Model(&CronTask{}).Count(&count)
	return count, nil
}

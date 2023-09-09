package model

import (
	"gorm.io/gorm"
)

type MonitorHighResourceProcesses struct {
	ID                int64   `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	ProcessName       string  `gorm:"column:process_name;not null" json:"process_name" form:"process_name"`
	Pid               int64   `gorm:"column:pid;not null" json:"pid" form:"pid"`
	CpuUsage          float64 `gorm:"column:cpu_usage;not null" json:"cpu_usage" form:"cpu_usage"`
	MemUsage          float64 `gorm:"column:mem_usage;not null" json:"mem_usage" form:"mem_usage"`
	Cmdline           string  `gorm:"column:cmdline;not null" json:"cmdline" form:"cmdline"`
	Username          string  `gorm:"column:username;not null" json:"username" form:"username"`
	ProcessCreateTime int64   `gorm:"column:process_create_time;not null" json:"process_create_time" form:"process_create_time"`
	NumThreads        int64   `gorm:"column:num_threads;not null" json:"num_threads" form:"num_threads"`
	CpuTimes          int64   `gorm:"column:cpu_times;not null" json:"cpu_times" form:"cpu_times"`
	IoCounters        int64   `gorm:"column:io_counters;not null" json:"io_counters" form:"io_counters"`
	OpenFiles         int64   `gorm:"column:open_files;not null" json:"open_files" form:"open_files"`
	CreateTime        int64   `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*MonitorHighResourceProcesses) TableName() string {
	return "monitor_high_resource_processes"
}

func (s *MonitorHighResourceProcesses) Get(db *gorm.DB) (*MonitorHighResourceProcesses, error) {
	var monitorHighResourceProcesses MonitorHighResourceProcesses
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&monitorHighResourceProcesses).Error
	if err != nil {
		return &monitorHighResourceProcesses, err
	}

	return &monitorHighResourceProcesses, nil
}
func (s *MonitorHighResourceProcesses) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*MonitorHighResourceProcesses, int64, error) {
	var monitorHighResourceProcesses []*MonitorHighResourceProcesses
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
	if err = db.Find(&monitorHighResourceProcesses).Error; err != nil {
		return nil, count, err
	}

	return monitorHighResourceProcesses, count, nil
}
func (s *MonitorHighResourceProcesses) Create(db *gorm.DB) (*MonitorHighResourceProcesses, error) {
	err := db.Create(&s).Error
	return s, err
}

func (s *MonitorHighResourceProcesses) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&MonitorHighResourceProcesses{}, s.ID).Error
}

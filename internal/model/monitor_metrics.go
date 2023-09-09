package model

import (
	"gorm.io/gorm"
)

type MonitorMetrics struct {
	ID                  int64   `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	CpuUsage            float64 `gorm:"column:cpu_usage;not null" json:"cpu_usage" form:"cpu_usage"`
	MemUsage            float64 `gorm:"column:mem_usage;not null" json:"mem_usage" form:"mem_usage"`
	LoadAvg1m           float64 `gorm:"column:load_avg_1m;not null" json:"load_avg_1m" form:"load_avg_1m"`
	LoadAvg5m           float64 `gorm:"column:load_avg_5m;not null" json:"load_avg_5m" form:"load_avg_5m"`
	LoadAvg15m          float64 `gorm:"column:load_avg_15m;not null" json:"load_avg_15m" form:"load_avg_15m"`
	ResourceUtilization float64 `gorm:"column:resource_utilization;not null" json:"resource_utilization" form:"resource_utilization"`
	CreateTime          int64   `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*MonitorMetrics) TableName() string {
	return "monitor_metrics"
}

func (s *MonitorMetrics) Get(db *gorm.DB) (*MonitorMetrics, error) {
	var monitorMetrics MonitorMetrics
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&monitorMetrics).Error
	if err != nil {
		return &monitorMetrics, err
	}

	return &monitorMetrics, nil
}
func (s *MonitorMetrics) List(db *gorm.DB, conditions *ConditionsT, startTime, endTime int64, offset, limit int) ([]*MonitorMetrics, error) {
	var monitorMetrics []*MonitorMetrics
	var err error

	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}

	if startTime > 0 && endTime > 0 {
		db = db.Where("create_time BETWEEN ? AND ?", startTime, endTime)
	}
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if err = db.Find(&monitorMetrics).Error; err != nil {
		return nil, err
	}
	return monitorMetrics, nil
}
func (s *MonitorMetrics) Create(db *gorm.DB) (*MonitorMetrics, error) {
	err := db.Create(&s).Error
	return s, err
}

func (s *MonitorMetrics) Delete(db *gorm.DB) error {
	return db.Delete(&MonitorMetrics{}, s.ID).Error
}
func (s *MonitorMetrics) BatchDelete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&MonitorMetrics{}).Error
}
func (s *MonitorMetrics) DeleteAll(db *gorm.DB) error {
	return db.Where("1 = 1").Delete(&MonitorMetrics{}).Error
}

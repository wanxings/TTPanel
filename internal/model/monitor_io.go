package model

import (
	"gorm.io/gorm"
)

type MonitorIo struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	ReadCount  uint64 `gorm:"column:read_count;not null" json:"read_count" form:"read_count"`
	WriteCount uint64 `gorm:"column:write_count;not null" json:"write_count" form:"write_count"`
	ReadBytes  uint64 `gorm:"column:read_bytes;not null" json:"read_bytes" form:"read_bytes"`
	WriteBytes uint64 `gorm:"column:write_bytes;not null" json:"write_bytes" form:"write_bytes"`
	ReadTime   uint64 `gorm:"column:read_time;not null" json:"read_time" form:"read_time"`
	WriteTime  uint64 `gorm:"column:write_time;not null" json:"write_time" form:"write_time"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*MonitorIo) TableName() string {
	return "monitor_io"
}
func (s *MonitorIo) Get(db *gorm.DB) (*MonitorIo, error) {
	var monitorIo MonitorIo
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&monitorIo).Error
	if err != nil {
		return &monitorIo, err
	}

	return &monitorIo, nil
}
func (s *MonitorIo) List(db *gorm.DB, conditions *ConditionsT, startTime, endTime int64, offset, limit int) ([]*MonitorIo, error) {
	var monitorIo []*MonitorIo
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
	if err = db.Find(&monitorIo).Error; err != nil {
		return nil, err
	}

	return monitorIo, nil
}
func (s *MonitorIo) Create(db *gorm.DB) (*MonitorIo, error) {
	err := db.Create(&s).Error
	return s, err
}

func (s *MonitorIo) Delete(db *gorm.DB) error {
	return db.Delete(&MonitorIo{}, s.ID).Error
}

func (s *MonitorIo) BatchDelete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&MonitorIo{}).Error
}

func (s *MonitorIo) DeleteAll(db *gorm.DB) error {
	return db.Where("1 = 1").Delete(&MonitorIo{}).Error
}

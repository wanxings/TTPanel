package model

import (
	"gorm.io/gorm"
)

type MonitorNetwork struct {
	ID          int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Up          uint64 `gorm:"column:up;not null" json:"up" form:"up"`
	Down        uint64 `gorm:"column:down;not null" json:"down" form:"down"`
	TotalUp     uint64 `gorm:"column:total_up;not null" json:"total_up" form:"total_up"`
	TotalDown   uint64 `gorm:"column:total_down;not null" json:"total_down" form:"total_down"`
	DownPackets uint64 `gorm:"column:down_packets;not null" json:"down_packets" form:"down_packets"`
	UpPackets   uint64 `gorm:"column:up_packets;not null" json:"up_packets" form:"up_packets"`
	CreateTime  int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*MonitorNetwork) TableName() string {
	return "monitor_network"
}

func (s *MonitorNetwork) Get(db *gorm.DB) (*MonitorNetwork, error) {
	var monitorNetwork MonitorNetwork
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&monitorNetwork).Error
	if err != nil {
		return &monitorNetwork, err
	}

	return &monitorNetwork, nil
}
func (s *MonitorNetwork) List(db *gorm.DB, conditions *ConditionsT, startTime, endTime int64, offset, limit int) ([]*MonitorNetwork, error) {
	var monitorNetwork []*MonitorNetwork
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
	if err = db.Find(&monitorNetwork).Error; err != nil {
		return nil, err
	}

	return monitorNetwork, nil
}
func (s *MonitorNetwork) Create(db *gorm.DB) (*MonitorNetwork, error) {
	err := db.Create(&s).Error
	return s, err
}

func (s *MonitorNetwork) Delete(db *gorm.DB) error {
	return db.Delete(&MonitorNetwork{}, s.ID).Error
}

func (s *MonitorNetwork) BatchDelete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&MonitorNetwork{}).Error
}

func (s *MonitorNetwork) DeleteAll(db *gorm.DB) error {
	return db.Where("1 = 1").Delete(&MonitorNetwork{}).Error
}

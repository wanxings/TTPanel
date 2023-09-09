package model

import (
	"gorm.io/gorm"
)

type MonitorEvent struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id"`
	Category   string `gorm:"column:category;not null" json:"category"`
	Log        string `gorm:"column:log" json:"log"`
	Status     int    `gorm:"column:status;not null" json:"status"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

func (*MonitorEvent) TableName() string {
	return "monitor_event"
}

func (s *MonitorEvent) Get(db *gorm.DB) (*MonitorEvent, error) {
	var project MonitorEvent
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&project).Error
	if err != nil {
		return &project, err
	}

	return &project, nil
}
func (s *MonitorEvent) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*MonitorEvent, int64, error) {
	var monitorEvent []*MonitorEvent
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
	if err = db.Find(&monitorEvent).Error; err != nil {
		return nil, 0, err
	}

	return monitorEvent, count, nil
}
func (s *MonitorEvent) Create(db *gorm.DB) (*MonitorEvent, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *MonitorEvent) UpdateOne(db *gorm.DB, filed string, value any, conditions *ConditionsT) error {
	for k, v := range *conditions {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	return db.Model(&MonitorEvent{}).Update(filed, value).Error
}
func (s *MonitorEvent) Update(db *gorm.DB) error {
	return db.Model(&MonitorEvent{}).Where("id = ?", s.ID).Save(s).Error
}
func (s *MonitorEvent) Count(db *gorm.DB, where *ConditionsT) (int64, error) {
	var count int64
	for k, v := range *where {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	db.Model(&MonitorEvent{}).Count(&count)
	return count, nil
}

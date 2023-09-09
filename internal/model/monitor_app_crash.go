package model

import (
	"gorm.io/gorm"
)

type MonitorAppCrash struct {
	ID           int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	AppName      string `gorm:"column:app_name;not null" json:"app_name" form:"app_name"`
	ErrorType    string `gorm:"column:error_type;not null" json:"error_type" form:"error_type"`
	ErrorMessage string `gorm:"column:error_message;not null" json:"error_message" form:"error_message"`
	CreateTime   int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*MonitorAppCrash) TableName() string {
	return "monitor_app_crash"
}

func (s *MonitorAppCrash) Get(db *gorm.DB) (*MonitorAppCrash, error) {
	var monitorAppCrash MonitorAppCrash
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&monitorAppCrash).Error
	if err != nil {
		return &monitorAppCrash, err
	}

	return &monitorAppCrash, nil
}
func (s *MonitorAppCrash) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*MonitorAppCrash, int64, error) {
	var monitorAppCrash []*MonitorAppCrash
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
	if err = db.Find(&monitorAppCrash).Error; err != nil {
		return nil, count, err
	}

	return monitorAppCrash, count, nil
}
func (s *MonitorAppCrash) Create(db *gorm.DB) (*MonitorAppCrash, error) {
	err := db.Create(&s).Error
	return s, err
}

func (s *MonitorAppCrash) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&MonitorAppCrash{}, s.ID).Error
}

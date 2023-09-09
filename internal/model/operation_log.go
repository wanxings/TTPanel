package model

import (
	"gorm.io/gorm"
)

type OperationLog struct {
	ID            int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Type          int    `gorm:"column:type;not null" json:"type" form:"type"`
	Uid           int64  `gorm:"column:uid;not null" json:"uid" form:"uid"`
	Username      string `gorm:"column:username;not null" json:"username" form:"username"`
	Log           string `gorm:"column:log;not null" json:"log" form:"log"`
	IP            string `gorm:"column:ip;not null" json:"ip" form:"ip"`
	IPAttribution string `gorm:"column:ip_attribution;not null" json:"ip_attribution" form:"ip_attribution"`
	CreateTime    int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*OperationLog) TableName() string {
	return "operation_log"
}

func (s *OperationLog) Get(db *gorm.DB) (*OperationLog, error) {
	var operationLog OperationLog
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&operationLog).Error
	if err != nil {
		return &operationLog, err
	}

	return &operationLog, nil
}
func (s *OperationLog) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*OperationLog, int64, error) {
	var operationLog []*OperationLog
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
	if err = db.Find(&operationLog).Error; err != nil {
		return nil, count, err
	}

	return operationLog, count, nil
}
func (s *OperationLog) Create(db *gorm.DB) (*OperationLog, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *OperationLog) Update(db *gorm.DB) error {
	return db.Model(&OperationLog{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *OperationLog) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&OperationLog{}, s.ID).Error
}

func (s *OperationLog) DeleteAll(db *gorm.DB) error {
	db = db.Where("1=1")
	return db.Delete(&OperationLog{}).Error
}

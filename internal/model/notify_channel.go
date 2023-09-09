package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type NotifyChannel struct {
	ID          int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Category    int    `gorm:"column:category;not null" json:"category" form:"category"`
	Name        string `gorm:"column:name;not null;unique" json:"name" form:"name"`
	Config      string `gorm:"column:config;not null" json:"config" form:"config"`
	Description string `gorm:"column:description" json:"description" form:"description"`
	CreateTime  int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*NotifyChannel) TableName() string {
	return "notify_channel"
}

func (t *NotifyChannel) Get(db *gorm.DB) (*NotifyChannel, error) {
	var notifyChannel NotifyChannel
	if t.ID > 0 {
		db = db.Where("id= ? ", t.ID)
	} else if !util.StrIsEmpty(t.Name) {
		db = db.Where("name = ?", t.Name)
	} else {
		return &notifyChannel, nil
	}

	if err := db.Limit(1).Find(&notifyChannel).Error; err != nil {
		return &notifyChannel, err
	}
	return &notifyChannel, nil
}
func (t *NotifyChannel) List(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT, offset, limit int) ([]*NotifyChannel, int64, error) {
	var notifyChannel []*NotifyChannel
	var err error
	var count int64

	for k, v := range *where {
		if k == "ORDER" {
			db = db.Order(v)
		} else if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	for k, v := range *whereOr {
		if k == "FIXED" {
			db = db.Or(v)
		} else {
			db = db.Or(k, v)
		}

	}
	if err = db.Model(t).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if err = db.Find(&notifyChannel).Error; err != nil {
		return nil, 0, err
	}

	return notifyChannel, count, nil
}
func (t *NotifyChannel) Create(db *gorm.DB) (*NotifyChannel, error) {
	err := db.Create(&t).Error
	return t, err
}

func (t *NotifyChannel) Updates(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	return db.Model(&NotifyChannel{}).Updates(&t).Error
}
func (t *NotifyChannel) Update(db *gorm.DB) error {
	return db.Model(&NotifyChannel{}).Where("id = ?", t.ID).Save(t).Error
}

func (t *NotifyChannel) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&NotifyChannel{}, t.ID).Error
}

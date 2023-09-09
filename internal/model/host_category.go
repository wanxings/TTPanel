package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type HostCategory struct {
	ID     int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Name   string `gorm:"column:name;not null;unique" json:"name" form:"name"`
	Remark string `gorm:"column:remark" json:"remark" form:"remark"`
}

func (*HostCategory) TableName() string {
	return "host_category"
}

func (t *HostCategory) Get(db *gorm.DB) (*HostCategory, error) {
	var hostCategory HostCategory
	if t.ID > 0 {
		db = db.Where("id= ? ", t.ID)
	} else if !util.StrIsEmpty(t.Name) {
		db = db.Where("name = ?", t.Name)
	} else {
		return &hostCategory, nil
	}

	if err := db.Limit(1).Find(&hostCategory).Error; err != nil {
		return &hostCategory, err
	}
	return &hostCategory, nil
}
func (t *HostCategory) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*HostCategory, int64, error) {
	var hostCategory []*HostCategory
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
	if err = db.Model(t).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if err = db.Find(&hostCategory).Error; err != nil {
		return nil, 0, err
	}

	return hostCategory, count, nil
}
func (t *HostCategory) Create(db *gorm.DB) (*HostCategory, error) {
	err := db.Create(&t).Error
	return t, err
}

func (t *HostCategory) Updates(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	return db.Model(&HostCategory{}).Updates(&t).Error
}
func (t *HostCategory) Update(db *gorm.DB) error {
	return db.Model(&HostCategory{}).Where("id = ?", t.ID).Save(t).Error
}

func (t *HostCategory) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&HostCategory{}, t.ID).Error
}

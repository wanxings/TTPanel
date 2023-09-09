package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type Storage struct {
	ID          int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Category    int    `gorm:"column:category;not null" json:"category" form:"category"`
	Name        string `gorm:"column:name;not null;unique" json:"name" form:"name"`
	Config      string `gorm:"column:config;not null" json:"config" form:"config"`
	Description string `gorm:"column:description" json:"description" form:"description"`
	CreateTime  int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*Storage) TableName() string {
	return "storage"
}

func (t *Storage) Get(db *gorm.DB) (*Storage, error) {
	var storage Storage
	if t.ID > 0 {
		db = db.Where("id= ? ", t.ID)
	} else if !util.StrIsEmpty(t.Name) {
		db = db.Where("name = ?", t.Name)
	} else {
		return &storage, nil
	}

	if err := db.Limit(1).Find(&storage).Error; err != nil {
		return &storage, err
	}
	return &storage, nil
}
func (t *Storage) List(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT, offset, limit int) ([]*Storage, int64, error) {
	var storage []*Storage
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
	if err = db.Find(&storage).Error; err != nil {
		return nil, 0, err
	}

	return storage, count, nil
}
func (t *Storage) Create(db *gorm.DB) (*Storage, error) {
	err := db.Create(&t).Error
	return t, err
}

func (t *Storage) Updates(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	return db.Model(&Storage{}).Updates(&t).Error
}
func (t *Storage) Update(db *gorm.DB) error {
	return db.Model(&Storage{}).Where("id = ?", t.ID).Save(t).Error
}

func (t *Storage) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&Storage{}, t.ID).Error
}

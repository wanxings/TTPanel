package model

import (
	"gorm.io/gorm"
)

type ProjectBinding struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id"`
	Pid        int64  `gorm:"column:pid;not null" json:"pid"`
	Domain     string `gorm:"column:domain;not null" json:"domain"`
	Path       string `gorm:"column:path" json:"path"`
	Port       int64  `gorm:"column:port;not null" json:"port"`
	CreateTime string `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

func (*ProjectBinding) TableName() string {
	return "project_binding"
}
func (s *ProjectBinding) Get(db *gorm.DB) (*ProjectBinding, error) {
	var projectBinding ProjectBinding
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&projectBinding).Error
	if err != nil {
		return &projectBinding, err
	}

	return &projectBinding, nil
}
func (s *ProjectBinding) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*ProjectBinding, int64, error) {
	var projectBinding []*ProjectBinding
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
	if err = db.Find(&projectBinding).Error; err != nil {
		return nil, count, err
	}

	return projectBinding, count, nil
}
func (s *ProjectBinding) Create(db *gorm.DB) (*ProjectBinding, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *ProjectBinding) Update(db *gorm.DB) error {
	return db.Model(&ProjectBinding{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *ProjectBinding) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&ProjectBinding{}, s.ID).Error
}

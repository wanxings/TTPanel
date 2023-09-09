package model

import (
	"errors"
	"gorm.io/gorm"
)

type ProjectDomain struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id"`
	ProjectId  int64  `gorm:"column:project_id;not null" json:"project_id"`
	Domain     string `gorm:"column:domain;not null" json:"domain"`
	Port       int    `gorm:"column:port;not null" json:"port"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

func (*ProjectDomain) TableName() string {
	return "project_domain"
}
func (s *ProjectDomain) Get(db *gorm.DB) (*ProjectDomain, error) {
	var projectDomain ProjectDomain
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&projectDomain).Error
	if err != nil {
		return &projectDomain, err
	}

	return &projectDomain, nil
}
func (s *ProjectDomain) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*ProjectDomain, int64, error) {
	var projectDomain []*ProjectDomain
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
	if err = db.Find(&projectDomain).Error; err != nil {
		return nil, count, err
	}

	return projectDomain, count, nil
}
func (s *ProjectDomain) Create(db *gorm.DB) (*ProjectDomain, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *ProjectDomain) Count(db *gorm.DB, conditions *ConditionsT) (count int64) {
	db = db.Model(s)
	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	db.Count(&count)
	return
}
func (s *ProjectDomain) Update(db *gorm.DB) error {
	return db.Model(&ProjectDomain{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *ProjectDomain) Delete(db *gorm.DB, conditions *ConditionsT) error {
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else if len(*conditions) == 0 {
		return errors.New("not delete condition")
	}
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&ProjectDomain{}).Error
}

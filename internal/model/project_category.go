package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type ProjectCategory struct {
	ID   int64  `gorm:"column:id;autoIncrement" db:"id" json:"id"`
	Name string `gorm:"column:name;not null;unique" db:"name" json:"name"`
	Ps   string `gorm:"column:ps" db:"ps" json:"ps"`
}

func (*ProjectCategory) TableName() string {
	return "project_category"
}
func (s *ProjectCategory) Get(db *gorm.DB) (*ProjectCategory, error) {
	var projectCategory ProjectCategory
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&projectCategory).Error
	if err != nil {
		return &projectCategory, err
	}

	return &projectCategory, nil
}
func (s *ProjectCategory) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*ProjectCategory, int64, error) {
	var projectCategory []*ProjectCategory
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
	if err = db.Find(&projectCategory).Error; err != nil {
		return nil, count, err
	}

	return projectCategory, count, nil
}
func (s *ProjectCategory) Create(db *gorm.DB) (*ProjectCategory, error) {
	err := db.Create(&s).Error
	return s, err
}

func (s *ProjectCategory) GetByName(db *gorm.DB) (*ProjectCategory, error) {
	var projectCategory ProjectCategory
	if !util.StrIsEmpty(s.Name) {
		db = db.Where("name", s.Name)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Find(&projectCategory).Error
	if err != nil {
		return &projectCategory, err
	}

	return &projectCategory, nil
}

func (s *ProjectCategory) Update(db *gorm.DB) error {
	return db.Model(&ProjectCategory{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *ProjectCategory) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&ProjectCategory{}, s.ID).Error
}

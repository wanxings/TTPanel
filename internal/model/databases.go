package model

import (
	"TTPanel/pkg/util"
	"errors"
	"gorm.io/gorm"
)

type Databases struct {
	ID          int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Pid         int64  `gorm:"column:pid" json:"pid" form:"pid"`
	DbType      int    `gorm:"column:db_type;not null" json:"db_type" form:"db_type"`
	Sid         int64  `gorm:"column:sid;not null" json:"sid" form:"sid"`
	Type        string `gorm:"column:type;not null" json:"type" form:"type"`
	Name        string `gorm:"column:name;not null;unique" json:"name" form:"name"`
	Username    string `gorm:"column:username;not null" json:"username" form:"username"`
	Password    string `gorm:"column:password;not null" json:"password" form:"password"`
	Accept      string `gorm:"column:accept;not null" json:"accept" form:"accept"`
	Ps          string `gorm:"column:ps" json:"ps" form:"ps"`
	CreateTime  int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
	BackupCount int64  `gorm:"-" json:"backup_count" form:"backup_count"`
}

func (*Databases) TableName() string {
	return "databases"
}
func (s *Databases) Get(db *gorm.DB) (*Databases, error) {
	var databases Databases
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else if !util.StrIsEmpty(s.Name) {
		db = db.Where("name", s.Name)
	} else if !util.StrIsEmpty(s.Username) {
		db = db.Where("username", s.Username)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&databases).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &databases, nil
	}

	if err != nil {
		return &databases, err
	}

	return &databases, nil
}
func (s *Databases) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*Databases, int64, error) {
	var databases []*Databases
	var err error
	var count int64

	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else if k == "FIXED" {
			db = db.Where(v)
		} else if k == "OR" {
			db = db.Or(v)
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
	if err = db.Find(&databases).Error; err != nil {
		return nil, count, err
	}
	return databases, count, nil
}
func (s *Databases) Create(db *gorm.DB) (*Databases, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *Databases) Update(db *gorm.DB) error {
	return db.Model(&Databases{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *Databases) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&Databases{}, s.ID).Error
}
func (s *Databases) Count(db *gorm.DB) (int64, error) {
	var count int64
	db.Model(&Databases{}).Count(&count)
	return count, nil
}

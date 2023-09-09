package model

import (
	"gorm.io/gorm"
)

type DockerRepository struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Url        string `gorm:"column:url;not null" json:"url" form:"url"`
	Username   string `gorm:"column:username" json:"username" form:"username"`
	Password   string `gorm:"column:password" json:"password" form:"password"`
	Name       string `gorm:"column:name;not null;unique" json:"name" form:"name"`
	Namespace  string `gorm:"column:namespace" json:"namespace" form:"namespace"`
	Remark     string `gorm:"column:remark" json:"remark" form:"remark"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*DockerRepository) TableName() string {
	return "repository"
}

func (s *DockerRepository) Get(db *gorm.DB) (*DockerRepository, error) {
	var repository DockerRepository
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else if s.Name != "" {
		db = db.Where("name", s.Name)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&repository).Error
	if err != nil {
		return &repository, err
	}

	return &repository, nil
}
func (s *DockerRepository) List(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT, offset, limit int) ([]*DockerRepository, int64, error) {
	var repository []*DockerRepository
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
	if err = db.Model(s).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if err = db.Find(&repository).Error; err != nil {
		return nil, count, err
	}

	return repository, count, nil
}
func (s *DockerRepository) Create(db *gorm.DB) (*DockerRepository, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *DockerRepository) Update(db *gorm.DB) error {
	return db.Model(&DockerRepository{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *DockerRepository) Delete(db *gorm.DB) error {
	return db.Delete(&DockerRepository{}, s.ID).Error
}

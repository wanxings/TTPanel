package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type DockerComposeTemplate struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Path       string `gorm:"column:path;not null" json:"path" form:"path"`
	AddInPath  int    `gorm:"column:add_in_path" json:"add_in_path" form:"add_in_path"`
	Name       string `gorm:"column:name;not null;unique" json:"name" form:"name"`
	Remark     string `gorm:"column:remark" json:"remark" form:"remark"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*DockerComposeTemplate) TableName() string {
	return "compose_template"
}

func (s *DockerComposeTemplate) Get(db *gorm.DB) (*DockerComposeTemplate, error) {
	var composeTemplate DockerComposeTemplate
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else if !util.StrIsEmpty(s.Name) {
		db = db.Where("name", s.Name)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&composeTemplate).Error
	if err != nil {
		return &composeTemplate, err
	}

	return &composeTemplate, nil
}
func (s *DockerComposeTemplate) List(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT, offset, limit int) ([]*DockerComposeTemplate, int64, error) {
	var composeTemplate []*DockerComposeTemplate
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
	if err = db.Find(&composeTemplate).Error; err != nil {
		return nil, count, err
	}

	return composeTemplate, count, nil
}
func (s *DockerComposeTemplate) Create(db *gorm.DB) (*DockerComposeTemplate, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *DockerComposeTemplate) Update(db *gorm.DB) error {
	return db.Model(&DockerComposeTemplate{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *DockerComposeTemplate) Delete(db *gorm.DB) error {
	return db.Delete(&DockerComposeTemplate{}, s.ID).Error
}

package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type DockerCompose struct {
	ID              int64              `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Name            string             `gorm:"column:name;not null;unique" json:"name" form:"name"`
	ServerName      string             `gorm:"column:server_name" json:"server_name" form:"server_name"`
	TemplateID      int64              `gorm:"column:template_id" json:"template_id" form:"template_id"`
	Path            string             `gorm:"column:path;not null" json:"path" form:"path"`
	ContainerNumber int                `gorm:"-" json:"containerNumber" form:"containerNumber"`
	Containers      []ComposeContainer `gorm:"-" json:"containers" form:"containers"`
	Remark          string             `gorm:"column:remark" json:"remark" form:"remark"`
	CreateTime      int64              `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}
type ComposeContainer struct {
	ContainerID string `json:"containerID"`
	Name        string `json:"name"`
	CreateTime  string `json:"createTime"`
	State       string `json:"state"`
}

func (*DockerCompose) TableName() string {
	return "compose"
}

func (s *DockerCompose) Get(db *gorm.DB) (*DockerCompose, error) {
	var compose DockerCompose
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else if !util.StrIsEmpty(s.Name) {
		db = db.Where("name", s.Name)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&compose).Error
	if err != nil {
		return &compose, err
	}

	return &compose, nil
}
func (s *DockerCompose) List(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT, offset, limit int) ([]*DockerCompose, int64, error) {
	var compose []*DockerCompose
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
	if err = db.Find(&compose).Error; err != nil {
		return nil, count, err
	}

	return compose, count, nil
}
func (s *DockerCompose) Create(db *gorm.DB) (*DockerCompose, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *DockerCompose) Count(db *gorm.DB, where *ConditionsT) (int64, error) {
	var count int64
	db = db.Model(s)
	for k, v := range *where {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	err := db.Count(&count).Error
	return count, err
}
func (s *DockerCompose) Update(db *gorm.DB) error {
	return db.Model(&DockerCompose{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *DockerCompose) Delete(db *gorm.DB) error {
	return db.Delete(&DockerCompose{}, s.ID).Error
}

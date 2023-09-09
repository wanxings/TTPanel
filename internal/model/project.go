package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type Project struct {
	ID            int64  `gorm:"column:id;autoIncrement" json:"id"`
	Name          string `gorm:"column:name;not null;unique" json:"name"`
	Path          string `gorm:"column:path;not null" json:"path"`
	CategoryId    int64  `gorm:"column:category_id;not null" json:"category_id"`
	ProjectType   int    `gorm:"column:project_type;not null" json:"project_type" form:"project_type"`
	ProjectConfig string `gorm:"column:project_config;not null" json:"project_config" form:"project_config"`
	ExpireTime    int64  `gorm:"column:expire_time" json:"expire_time"`
	Status        int    `gorm:"column:status;not null" json:"status"`
	Ps            string `gorm:"column:ps" json:"ps"`
	CreateTime    int64  `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

func (*Project) TableName() string {
	return "project"
}

type GeneralProjectConfig struct {
	Command      string `json:"command"`
	IsPowerOn    bool   `json:"is_power_on"`
	RunUser      string `json:"run_user"`
	Port         int    `json:"port"`
	ExternalBind bool   `json:"external_bind"`
}

func (s *Project) Get(db *gorm.DB) (*Project, error) {
	var project Project
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else if !util.StrIsEmpty(s.Name) {
		db = db.Where("name", s.Name)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	if s.ProjectType > 0 {
		db = db.Where("project_type", s.ProjectType)
	}
	err := db.Limit(1).Find(&project).Error
	if err != nil {
		return &project, err
	}

	return &project, nil
}

// List 结构体需设置ProjectType，否则默认是php项目
func (s *Project) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*Project, int64, error) {
	var project []*Project
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
	if err = db.Model(&s).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if err = db.Find(&project).Error; err != nil {
		return nil, count, err
	}

	return project, count, nil
}

// Search 结构体需设置ProjectType，否则默认是php项目
func (s *Project) Search(db *gorm.DB, conditions *ConditionsT, query string) ([]int64, error) {
	var pIds []int64
	var err error
	query = "%" + query + "%"
	tnProject := s.TableName()
	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	tnDomainGenre := db.NamingStrategy.TableName("project_domain")
	db = db.Joins("LEFT JOIN " + tnDomainGenre + " ON " + tnDomainGenre + ".project_id = " + tnProject + ".id ")
	db = db.Where("name LIKE ? OR ps LIKE ? OR domain LIKE ?", query, query, query)
	if err = db.Model(&s).Distinct().Pluck(tnProject+".id", &pIds).Error; err != nil {
		return nil, err
	}
	return pIds, nil
}
func (s *Project) Create(db *gorm.DB) (*Project, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *Project) Update(db *gorm.DB) error {
	return db.Model(&Project{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *Project) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&Project{}, s.ID).Error
}
func (s *Project) Count(db *gorm.DB) (int64, error) {
	var count int64
	if s.Status > 0 {
		db = db.Where("status= ? ", s.Status)
	}
	db.Model(&Project{}).Count(&count)
	return count, nil
}

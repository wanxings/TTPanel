package model

import (
	"gorm.io/gorm"
)

type DatabaseServer struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	DbType     string `gorm:"column:db_type;not null" json:"db_type" form:"db_type"`
	Charset    string `gorm:"column:charset" json:"charset" form:"charset"`
	Host       string `gorm:"column:host;not null" json:"host" form:"host"`
	Port       int64  `gorm:"column:port;not null" json:"port" form:"port"`
	User       string `gorm:"column:user;not null" json:"user" form:"user"`
	Password   string `gorm:"column:password;not null" json:"password" form:"password"`
	Ps         string `gorm:"column:ps" json:"ps" form:"ps"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*DatabaseServer) TableName() string {
	return "database_server"
}

func (s *DatabaseServer) Get(db *gorm.DB) (*DatabaseServer, error) {
	var databaseServer DatabaseServer
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&databaseServer).Error
	if err != nil {
		return &databaseServer, err
	}

	return &databaseServer, nil
}
func (s *DatabaseServer) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*DatabaseServer, int64, error) {
	var databaseServer []*DatabaseServer
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
	if err = db.Find(&databaseServer).Error; err != nil {
		return nil, count, err
	}

	return databaseServer, count, nil
}
func (s *DatabaseServer) Create(db *gorm.DB) (*DatabaseServer, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *DatabaseServer) Update(db *gorm.DB) error {
	return db.Model(&DatabaseServer{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *DatabaseServer) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&DatabaseServer{}, s.ID).Error
}

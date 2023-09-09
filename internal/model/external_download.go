package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
	"time"
)

type ExternalDownload struct {
	ID          int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Token       string `gorm:"column:token;not null;unique" json:"token" form:"token"`
	FilePath    string `gorm:"column:file_path;not null" json:"file_path" form:"file_path"`
	Description string `gorm:"column:description" json:"description" form:"description"`
	ExpireTime  int64  `gorm:"column:expire_time" json:"expire_time" form:"expire_time"`
	CreateTime  int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*ExternalDownload) TableName() string {
	return "external_download"
}

func (t *ExternalDownload) Get(db *gorm.DB) (*ExternalDownload, error) {
	var externalDownload ExternalDownload
	if t.ID > 0 {
		db = db.Where("id= ? ", t.ID)
	} else if !util.StrIsEmpty(t.Token) {
		db = db.Where("token = ?", t.Token)
	} else {
		return &externalDownload, nil
	}
	db = db.Where("expire_time> ? ", time.Now().Unix())
	if err := db.Limit(1).Find(&externalDownload).Error; err != nil {
		return &externalDownload, err
	}
	return &externalDownload, nil
}
func (t *ExternalDownload) Create(db *gorm.DB) (*ExternalDownload, error) {
	err := db.Create(&t).Error
	return t, err
}
func (t *ExternalDownload) List(db *gorm.DB) ([]*ExternalDownload, error) {
	var externalDownload []*ExternalDownload
	var err error
	db = db.Order("create_time DESC")
	if err = db.Find(&externalDownload).Error; err != nil {
		return nil, err
	}
	return externalDownload, nil
}

func (t *ExternalDownload) CleanExpired(db *gorm.DB) error {
	db = db.Where("expire_time < ? ", time.Now().Unix())
	return db.Delete(&ExternalDownload{}, t.ID).Error
}
func (t *ExternalDownload) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&ExternalDownload{}, t.ID).Error
}

package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type Host struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	CID        int64  `gorm:"column:cid;not null" json:"cid" form:"cid"`
	Name       string `gorm:"column:name;not null;unique" json:"name" form:"name"`
	Address    string `gorm:"column:address" json:"address" form:"address"`
	Port       int    `gorm:"column:port" json:"port" form:"port"`
	User       string `gorm:"column:user" json:"user" form:"user"`
	Password   string `gorm:"column:password" json:"password" form:"password"`
	PrivateKey string `gorm:"column:private_key" json:"private_key" form:"private_key"`
	Remark     string `gorm:"column:remark" json:"remark" form:"remark"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*Host) TableName() string {
	return "host"
}

func (t *Host) Get(db *gorm.DB) (*Host, error) {
	var host Host
	if t.ID > 0 {
		db = db.Where("id= ? ", t.ID)
	} else if !util.StrIsEmpty(t.Name) {
		db = db.Where("name = ?", t.Name)
	} else if t.CID > 0 {
		db = db.Where("cid = ?", t.CID)
	} else {
		return &host, nil
	}

	if err := db.Limit(1).Find(&host).Error; err != nil {
		return &host, err
	}
	return &host, nil
}
func (t *Host) List(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT, offset, limit int) ([]*Host, int64, error) {
	var host []*Host
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
	if err = db.Model(t).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if err = db.Find(&host).Error; err != nil {
		return nil, 0, err
	}

	return host, count, nil
}
func (t *Host) Create(db *gorm.DB) (*Host, error) {
	err := db.Create(&t).Error
	return t, err
}

func (t *Host) Updates(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	return db.Model(&Host{}).Updates(&t).Error
}
func (t *Host) Update(db *gorm.DB) error {
	return db.Model(&Host{}).Where("id = ?", t.ID).Save(t).Error
}

func (t *Host) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&Host{}, t.ID).Error
}

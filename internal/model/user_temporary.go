package model

import (
	"gorm.io/gorm"
)

type TemporaryUser struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Token      string `gorm:"column:token;not null" json:"token" form:"token"`
	Remark     string `gorm:"column:remark" json:"remark" form:"remark"`
	ExpireTime int    `gorm:"column:expire_time;not null" json:"expire_time" form:"expire_time"`
	LoginIp    string `gorm:"column:login_ip;not null" json:"login_ip" form:"login_ip"`
	LoginTime  int64  `gorm:"column:login_time;not null" json:"login_time" form:"login_time"`
}

func (*TemporaryUser) TableName() string {
	return "temporary_user"
}
func (u *TemporaryUser) Get(db *gorm.DB) (*TemporaryUser, error) {
	var user TemporaryUser
	if u.ID > 0 {
		db = db.Where("id= ? ", u.ID)
	} else {
		return &user, nil
	}

	err := db.Limit(1).Find(&user).Error
	if err != nil {
		return &user, err
	}

	return &user, nil
}
func (u *TemporaryUser) List(db *gorm.DB, where *ConditionsT, offset, limit int) ([]*TemporaryUser, int64, error) {
	var temporaryUser []*TemporaryUser
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
	if err = db.Model(u).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if err = db.Find(&temporaryUser).Error; err != nil {
		return nil, 0, err
	}

	return temporaryUser, count, nil
}
func (u *TemporaryUser) Create(db *gorm.DB) (*TemporaryUser, error) {
	err := db.Create(&u).Error
	return u, err
}

func (u *TemporaryUser) Update(db *gorm.DB) error {
	return db.Model(&User{}).Where("id = ?", u.ID).Save(u).Error
}

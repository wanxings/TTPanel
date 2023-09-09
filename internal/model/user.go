package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type User struct {
	ID        int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Username  string `gorm:"column:username;not null" json:"username" form:"username"`
	Password  string `gorm:"column:password;not null" json:"password" form:"password"`
	LoginIp   string `gorm:"column:login_ip;not null" json:"login_ip" form:"login_ip"`
	LoginTime int64  `gorm:"column:login_time;not null" json:"login_time" form:"login_time"`
	Email     string `gorm:"column:email;not null" json:"email" form:"email"`
	Salt      string `gorm:"column:salt;not null" json:"salt" form:"salt"`
}

func (*User) TableName() string {
	return "user"
}
func (u *User) Get(db *gorm.DB) (*User, error) {
	var user User
	if u.ID > 0 {
		db = db.Where("id= ? ", u.ID)
	} else if !util.StrIsEmpty(u.Username) {
		db = db.Where("username = ?", u.Username)
	} else if !util.StrIsEmpty(u.Email) {
		db = db.Where("email = ?", u.Email, 0, 2)
	} else {
		return &user, nil
	}

	err := db.Limit(1).Find(&user).Error
	if err != nil {
		return &user, err
	}

	return &user, nil
}
func (u *User) Create(db *gorm.DB) (*User, error) {
	err := db.Create(&u).Error
	return u, err
}

func (u *User) Update(db *gorm.DB) error {
	return db.Model(&User{}).Where("id = ?", u.ID).Save(u).Error
}

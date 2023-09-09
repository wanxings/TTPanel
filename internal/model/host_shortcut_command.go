package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type HostShortcutCommand struct {
	ID          int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Name        string `gorm:"column:name;not null" json:"name" form:"name"`
	Cmd         string `gorm:"column:cmd;not null" json:"cmd" form:"cmd"`
	Description string `gorm:"column:description" json:"description" form:"description"`
	CreateTime  int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*HostShortcutCommand) TableName() string {
	return "host_shortcut_command"
}

func (t *HostShortcutCommand) Get(db *gorm.DB) (*HostShortcutCommand, error) {
	var hostShortcutCommand HostShortcutCommand
	if t.ID > 0 {
		db = db.Where("id= ? ", t.ID)
	} else if !util.StrIsEmpty(t.Name) {
		db = db.Where("name = ?", t.Name)
	} else {
		return &hostShortcutCommand, nil
	}

	if err := db.Limit(1).Find(&hostShortcutCommand).Error; err != nil {
		return &hostShortcutCommand, err
	}
	return &hostShortcutCommand, nil
}
func (t *HostShortcutCommand) List(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT, offset, limit int) ([]*HostShortcutCommand, int64, error) {
	var hostShortcutCommand []*HostShortcutCommand
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
	if err = db.Find(&hostShortcutCommand).Error; err != nil {
		return nil, 0, err
	}

	return hostShortcutCommand, count, nil
}
func (t *HostShortcutCommand) Create(db *gorm.DB) (*HostShortcutCommand, error) {
	err := db.Create(&t).Error
	return t, err
}

func (t *HostShortcutCommand) Updates(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	return db.Model(&HostShortcutCommand{}).Updates(&t).Error
}
func (t *HostShortcutCommand) Update(db *gorm.DB) error {
	return db.Model(&HostShortcutCommand{}).Where("id = ?", t.ID).Save(t).Error
}

func (t *HostShortcutCommand) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&HostShortcutCommand{}, t.ID).Error
}

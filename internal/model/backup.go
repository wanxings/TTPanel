package model

import (
	"gorm.io/gorm"
)

type Backup struct {
	ID          int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Category    int    `gorm:"column:category;not null" json:"category" form:"category"`
	Pid         int64  `gorm:"column:pid" json:"pid" form:"pid"`
	CronTaskID  int64  `gorm:"column:cron_task_id" json:"cron_task_id" form:"cron_task_id"`
	StorageId   int64  `gorm:"column:storage_id" json:"storage_id" form:"storage_id"`
	FileName    string `gorm:"column:file_name;not null" json:"file_name" form:"file_name"`
	FilePath    string `gorm:"column:file_path;not null" json:"file_path" form:"file_path"`
	Size        int64  `gorm:"column:size;not null" json:"size" form:"size"`
	Description string `gorm:"column:description" json:"description" form:"description"`
	CreateTime  int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*Backup) TableName() string {
	return "backup"
}
func (t *Backup) Get(db *gorm.DB) (*Backup, error) {
	var backup Backup
	if t.ID > 0 {
		db = db.Where("id= ? ", t.ID)
	} else {
		return &backup, nil
	}

	if err := db.Limit(1).Find(&backup).Error; err != nil {
		return &backup, err
	}
	return &backup, nil
}
func (t *Backup) List(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT, offset, limit int) ([]*Backup, int64, error) {
	var backup []*Backup
	var err error
	var count int64
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}

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
	if err = db.Find(&backup).Error; err != nil {
		return nil, 0, err
	}
	if err = db.Model(t).Count(&count).Error; err != nil {
		return nil, count, err
	}
	return backup, count, nil
}
func (t *Backup) Count(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT) (int64, error) {
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
	if err = db.Model(&Backup{}).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}
func (t *Backup) Create(db *gorm.DB) (*Backup, error) {
	err := db.Create(&t).Error
	return t, err
}

func (t *Backup) Updates(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	return db.Model(&Backup{}).Updates(&t).Error
}
func (t *Backup) Update(db *gorm.DB) error {
	return db.Model(&Backup{}).Where("id = ?", t.ID).Save(t).Error
}

func (t *Backup) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&Backup{}, t.ID).Error
}

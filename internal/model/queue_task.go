package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type QueueTask struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Name       string `gorm:"column:name;not null" json:"name" form:"name"`
	Type       int    `gorm:"column:type;not null" json:"type" form:"type"`
	Status     int    `gorm:"column:status;not null" json:"status" form:"status"` //1.进行中2.已完成3.异常
	StartTime  int64  `gorm:"column:start_time" json:"start_time" form:"start_time"`
	EndTime    int64  `gorm:"column:end_time" json:"end_time" form:"end_time"`
	ExecStr    string `gorm:"column:exec_str;not null" json:"exec_str" form:"exec_str"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" gorm:"autoCreateTime" json:"create_time" form:"create_time"`
}

func (*QueueTask) TableName() string {
	return "queue_task"
}

func (t *QueueTask) Get(db *gorm.DB) (*QueueTask, error) {
	var taskQueue QueueTask
	if t.ID > 0 {
		db = db.Where("id= ? ", t.ID)
	} else if !util.StrIsEmpty(t.Name) {
		db = db.Where("name = ?", t.Name)
	} else if t.Type > 0 {
		db = db.Where("type = ?", t.Type)
	} else if t.Status > 0 {
		db = db.Where("type = ?", t.Type)
	} else {
		return &taskQueue, nil
	}

	if err := db.Limit(1).Find(&taskQueue).Error; err != nil {
		return &taskQueue, err
	}
	return &taskQueue, nil
}
func (t *QueueTask) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*QueueTask, int64, error) {
	var taskQueue []*QueueTask
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
	if err = db.Model(t).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	if err = db.Find(&taskQueue).Error; err != nil {
		return nil, 0, err
	}

	return taskQueue, count, nil
}
func (t *QueueTask) Create(db *gorm.DB) (*QueueTask, error) {
	err := db.Create(&t).Error
	return t, err
}

func (t *QueueTask) UpdateOne(db *gorm.DB, filed string, value any, conditions *ConditionsT) error {
	for k, v := range *conditions {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	return db.Model(&QueueTask{}).Update(filed, value).Error
}
func (t *QueueTask) Update(db *gorm.DB) error {
	return db.Model(&QueueTask{}).Where("id = ?", t.ID).Save(t).Error
}
func (t *QueueTask) Delete(db *gorm.DB) error {
	return db.Model(&QueueTask{}).Where("id = ?", t.ID).Delete(&QueueTask{}).Error
}
func (t *QueueTask) Count(db *gorm.DB, where ConditionsT) (int64, error) {
	var count int64
	for k, v := range where {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	db.Model(&QueueTask{}).Count(&count)
	return count, nil
}

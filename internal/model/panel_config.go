package model

import (
	"TTPanel/pkg/util"
	"gorm.io/gorm"
)

type PanelConfig struct {
	Key   string `gorm:"column:key;not null;unique" json:"key" form:"key"`
	Value string `gorm:"column:value" json:"value" form:"value"`
}

func (*PanelConfig) TableName() string {
	return "panel_config"
}

func (t *PanelConfig) Get(db *gorm.DB) (*PanelConfig, error) {
	var panelConfig PanelConfig
	if !util.StrIsEmpty(t.Key) {
		db = db.Where("key = ?", t.Key)
	} else {
		return &panelConfig, nil
	}

	if err := db.Limit(1).Find(&panelConfig).Error; err != nil {
		return &panelConfig, err
	}
	return &panelConfig, nil
}
func (t *PanelConfig) Create(db *gorm.DB) (*PanelConfig, error) {
	err := db.Create(&t).Error
	return t, err
}

func (t *PanelConfig) Updates(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	return db.Model(&PanelConfig{}).Updates(&t).Error
}
func (t *PanelConfig) Update(db *gorm.DB) error {
	return db.Model(&PanelConfig{}).Where("key = ?", t.Key).Save(t).Error
}

func (t *PanelConfig) Delete(db *gorm.DB) error {
	return db.Delete(&PanelConfig{}, t.Key).Error
}

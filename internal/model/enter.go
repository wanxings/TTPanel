package model

import (
	"gorm.io/gorm"
)

// Model 公共Model
type Model struct {
	ID int64 `gorm:"primary_key" json:"id"`
}
type ConditionsT map[string]interface{}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	return
}

func (m *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	// if !tx.Statement.Changed("modified_on") {
	// 	tx.Statement.SetColumn("modified_on", time.Now().Unix())
	// }

	return
}

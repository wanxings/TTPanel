package model

import (
	"gorm.io/gorm"
)

type FirewallRuleIp struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Ip         string `gorm:"column:ip;not null" json:"ip" form:"ip"`
	Strategy   int    `gorm:"column:strategy;not null" json:"strategy" form:"strategy"`
	Ps         string `gorm:"column:ps" json:"ps" form:"ps"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*FirewallRuleIp) TableName() string {
	return "firewall_rule_ip"
}

func (s *FirewallRuleIp) Get(db *gorm.DB) (*FirewallRuleIp, error) {
	var firewallRuleIp FirewallRuleIp
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&firewallRuleIp).Error
	if err != nil {
		return &firewallRuleIp, err
	}

	return &firewallRuleIp, nil
}
func (s *FirewallRuleIp) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*FirewallRuleIp, int64, error) {
	var firewallRuleIp []*FirewallRuleIp
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
	if err = db.Find(&firewallRuleIp).Error; err != nil {
		return nil, count, err
	}

	return firewallRuleIp, count, nil
}
func (s *FirewallRuleIp) Create(db *gorm.DB) (*FirewallRuleIp, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *FirewallRuleIp) Update(db *gorm.DB) error {
	return db.Model(&FirewallRuleIp{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *FirewallRuleIp) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&FirewallRuleIp{}, s.ID).Error
}

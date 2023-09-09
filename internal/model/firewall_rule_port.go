package model

import (
	"gorm.io/gorm"
)

type FirewallRulePort struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Port       string `gorm:"column:port;not null" json:"port" form:"port"`
	Strategy   int    `gorm:"column:strategy;not null" json:"strategy" form:"strategy"`
	SourceIp   string `gorm:"column:source_ip" json:"source_ip" form:"source_ip"`
	Protocol   string `gorm:"column:protocol;not null" json:"protocol" form:"protocol"`
	Ps         string `gorm:"column:ps" json:"ps" form:"ps"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*FirewallRulePort) TableName() string {
	return "firewall_rule_port"
}
func (s *FirewallRulePort) Get(db *gorm.DB) (*FirewallRulePort, error) {
	var firewallRulePort FirewallRulePort
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&firewallRulePort).Error
	if err != nil {
		return &firewallRulePort, err
	}

	return &firewallRulePort, nil
}
func (s *FirewallRulePort) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*FirewallRulePort, int64, error) {
	var firewallRulePort []*FirewallRulePort
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
	if err = db.Find(&firewallRulePort).Error; err != nil {
		return nil, count, err
	}

	return firewallRulePort, count, nil
}
func (s *FirewallRulePort) Create(db *gorm.DB) (*FirewallRulePort, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *FirewallRulePort) Update(db *gorm.DB) error {
	return db.Model(&FirewallRulePort{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *FirewallRulePort) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&FirewallRulePort{}, s.ID).Error
}

package model

import (
	"gorm.io/gorm"
)

type FirewallRuleForward struct {
	ID         int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	SourcePort int64  `gorm:"column:source_port;not null" json:"source_port" form:"source_port"`
	TargetIp   string `gorm:"column:target_ip;not null" json:"target_ip" form:"target_ip"`
	TargetPort int64  `gorm:"column:target_port;not null" json:"target_port" form:"target_port"`
	Protocol   string `gorm:"column:protocol;not null" json:"protocol" form:"protocol"`
	Ps         string `gorm:"column:ps" json:"ps" form:"ps"`
	CreateTime int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*FirewallRuleForward) TableName() string {
	return "firewall_rule_forward"
}

func (s *FirewallRuleForward) Get(db *gorm.DB) (*FirewallRuleForward, error) {
	var firewallRuleForward FirewallRuleForward
	if s.ID > 0 {
		db = db.Where("id", s.ID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.Limit(1).Find(&firewallRuleForward).Error
	if err != nil {
		return &firewallRuleForward, err
	}

	return &firewallRuleForward, nil
}
func (s *FirewallRuleForward) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*FirewallRuleForward, int64, error) {
	var firewallRuleForward []*FirewallRuleForward
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
	if err = db.Find(&firewallRuleForward).Error; err != nil {
		return nil, count, err
	}

	return firewallRuleForward, count, nil
}
func (s *FirewallRuleForward) Create(db *gorm.DB) (*FirewallRuleForward, error) {
	err := db.Create(&s).Error
	return s, err
}
func (s *FirewallRuleForward) Update(db *gorm.DB) error {
	return db.Model(&FirewallRuleForward{}).Where("id = ? ", s.ID).Save(s).Error
}
func (s *FirewallRuleForward) Delete(db *gorm.DB, conditions *ConditionsT) error {
	for k, v := range *conditions {
		db = db.Where(k, v)
	}
	return db.Delete(&FirewallRuleForward{}, s.ID).Error
}

package model

import (
	"gorm.io/gorm"
)

type TTWafBanIpLog struct {
	ID             int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	ServerName     string `gorm:"column:server_name;not null" json:"server_name" form:"server_name"`
	Reason         string `gorm:"column:reason;not null" json:"reason" form:"reason"`
	RuleDetailed   string `gorm:"column:rule_detailed;not null" json:"rule_detailed" form:"rule_detailed"`
	Ip             string `gorm:"column:ip;not null" json:"ip" form:"ip"`
	IpCity         string `gorm:"column:ip_city;not null" json:"ip_city" form:"ip_city"`
	IpCountry      string `gorm:"column:ip_country;not null" json:"ip_country" form:"ip_country"`
	IpSubdivisions string `gorm:"column:ip_subdivisions;not null" json:"ip_subdivisions" form:"ip_subdivisions"`
	IpContinent    string `gorm:"column:ip_continent;not null" json:"ip_continent" form:"ip_continent"`
	IpLongitude    string `gorm:"column:ip_longitude;not null" json:"ip_longitude" form:"ip_longitude"`
	IpLatitude     string `gorm:"column:ip_latitude;not null" json:"ip_latitude" form:"ip_latitude"`
	UserAgent      string `gorm:"column:user_agent;not null" json:"user_agent" form:"user_agent"`
	BlockTotal     int64  `gorm:"column:block_total;not null" json:"block_total" form:"block_total"`
	BlockTime      int64  `gorm:"column:block_time;not null" json:"block_time" form:"block_time"`
	Status         int64  `gorm:"column:status;not null" json:"status" form:"status"`
	CreateTime     int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*TTWafBanIpLog) TableName() string {
	return "ban_ip_log"
}
func (t *TTWafBanIpLog) List(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT, offset, limit int) ([]*TTWafBanIpLog, int64, error) {
	var ttWafBanIpLog []*TTWafBanIpLog
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
	if err = db.Find(&ttWafBanIpLog).Error; err != nil {
		return nil, 0, err
	}

	return ttWafBanIpLog, count, nil
}
func (t *TTWafBanIpLog) BatchUpdates(db *gorm.DB, updateFiled string, where *ConditionsT) error {
	for k, v := range *where {
		if k == "ORDER" {
			db = db.Order(v)
		} else if k == "FIXED" {
			db = db.Where(v)
		} else {
			db = db.Where(k, v)
		}
	}
	return db.Model(t).Select(updateFiled).Updates(t).Error
}
func (t *TTWafBanIpLog) Count(db *gorm.DB, where *ConditionsT) (int64, error) {
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
	if err := db.Model(t).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

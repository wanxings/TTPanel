package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type TTWafBlockIpLog struct {
	ID              int64  `gorm:"column:id;autoIncrement" json:"id" form:"id"`
	Hash            string `gorm:"column:hash;not null" json:"hash" form:"hash"`
	ServerName      string `gorm:"column:server_name;not null" json:"server_name" form:"server_name"`
	Ip              string `gorm:"column:ip;not null" json:"ip" form:"ip"`
	IpCity          string `gorm:"column:ip_city;not null" json:"ip_city" form:"ip_city"`
	IpCountry       string `gorm:"column:ip_country;not null" json:"ip_country" form:"ip_country"`
	IpSubdivisions  string `gorm:"column:ip_subdivisions;not null" json:"ip_subdivisions" form:"ip_subdivisions"`
	IpContinent     string `gorm:"column:ip_continent;not null" json:"ip_continent" form:"ip_continent"`
	IpLongitude     string `gorm:"column:ip_longitude;not null" json:"ip_longitude" form:"ip_longitude"`
	IpLatitude      string `gorm:"column:ip_latitude;not null" json:"ip_latitude" form:"ip_latitude"`
	RuleType        string `gorm:"column:rule_type;not null" json:"rule_type" form:"rule_type"`
	Uri             string `gorm:"column:uri;not null" json:"uri" form:"uri"`
	UserAgent       string `gorm:"column:user_agent;not null" json:"user_agent" form:"user_agent"`
	FilterRule      string `gorm:"column:filter_rule;not null" json:"filter_rule" form:"filter_rule"`
	RuleDescription string `gorm:"column:rule_description;not null" json:"rule_description" form:"rule_description"`
	SourceValue     string `gorm:"column:source_value;not null" json:"source_value" form:"source_value"`
	RuleReg         string `gorm:"column:rule_reg;not null" json:"rule_reg" form:"rule_reg"`
	RiskValue       string `gorm:"column:risk_value;not null" json:"risk_value" form:"risk_value"`
	CreateTime      int64  `gorm:"column:create_time;autoCreateTime" json:"create_time" form:"create_time"`
}

func (*TTWafBlockIpLog) TableName() string {
	return "block_ip_log"
}

type CountByDayData struct {
	Date  string
	Count int
}

type TopServerNamesData struct {
	ServerName string
	Count      int
}

func (t *TTWafBlockIpLog) List(db *gorm.DB, where *ConditionsT, whereOr *ConditionsT, offset, limit int) ([]*TTWafBlockIpLog, int64, error) {
	var ttWafBlockIpLog []*TTWafBlockIpLog
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
	if err = db.Find(&ttWafBlockIpLog).Error; err != nil {
		return nil, 0, err
	}

	return ttWafBlockIpLog, count, nil
}

func (t *TTWafBlockIpLog) Count(db *gorm.DB, where *ConditionsT) (int64, error) {
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

func (t *TTWafBlockIpLog) CountByDay(db *gorm.DB, days int) ([]CountByDayData, error) {
	var result []CountByDayData
	err := db.Model(t).
		Select("strftime('%Y-%m-%d',create_time, 'unixepoch')  as date, COUNT(id) as count").
		Where("create_time > ?", time.Now().AddDate(0, 0, -days).Unix()).
		Group("date").
		Order("date").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (t *TTWafBlockIpLog) TopServerNames(db *gorm.DB, countField string, limit int) ([]TopServerNamesData, error) {
	var result []TopServerNamesData
	err := db.Model(t).
		Select(fmt.Sprintf("%v, COUNT(*) as count", countField)).
		Where("create_time >= ? AND create_time < ?", time.Now().Unix()-int64(time.Now().Hour()*3600+time.Now().Minute()*60+time.Now().Second()), time.Now().Unix()-int64(time.Now().Hour()*3600+time.Now().Minute()*60+time.Now().Second()-86400)).
		Group(countField).
		Order("count DESC").
		Limit(limit).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

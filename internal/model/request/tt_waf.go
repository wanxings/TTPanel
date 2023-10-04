package request

import (
	"TTPanel/internal/model/response"
)

type TTWafGlobalSetR struct {
	Cc              *response.TTWafCCSetting              `json:"cc"`
	AttackTolerance *response.TTWafAttackToleranceSetting `json:"attack_tolerance"`
	Analytics       *response.TTWafAnalyticsSetting       `json:"analytics"`
}

type TTWafSaveProjectConfig struct {
	ProjectID int64 `json:"project_id" binding:"required"`
	response.TTWafProjectConfig
}

type TTWafProjectConfigR struct {
	ProjectID int64 `json:"project_id" binding:"required"`
}

type TTWafBlockListR struct {
	Ip        string `json:"ip"`
	Url       string `json:"url"`
	Domain    string `json:"domain"`
	UserAgent string `json:"user_agent"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	Limit     int    `json:"limit" binding:"required"`
	Page      int    `json:"page" binding:"required"`
}

type TTWafBanListR struct {
	Ip        string `json:"ip"`
	Url       string `json:"url"`
	Domain    string `json:"domain"`
	UserAgent string `json:"user_agent"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	Limit     int    `json:"limit" binding:"required"`
	Page      int    `json:"page" binding:"required"`
}

type TTWafGetRegRuleR struct {
	RuleName string `json:"rule_name" binding:"required"`
}

type TTWafSaveRegRuleR struct {
	RuleName string                  `json:"rule_name" binding:"required"`
	RuleList []response.TTWafRegRule `json:"rule_list" binding:"required"`
}

type TTWafIPR struct {
	Ips []string `json:"ips" binding:"required"`
}

type TTWafIPBlackWhiteR struct {
	IPV4 [][]int64 `json:"ipv4"`
	IPV6 []string  `json:"ipv6"`
}

type AnalyticsOverviewR struct {
	ProjectId int64 `json:"project_id" binding:"required"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

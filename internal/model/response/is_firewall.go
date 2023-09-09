package response

import (
	"TTPanel/internal/model"
)

type FirewallInfo struct {
	Ping   bool   `json:"ping"`
	Status bool   `json:"status"`
	Name   string `json:"name"`
}

type FirewallRulePortResp struct {
	Status bool `json:"status"`
	*model.FirewallRulePort
}

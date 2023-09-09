package request

type BatchCreatePortRuleR struct {
	Rules []*CreatePortRuleR `json:"rules" binding:"required"`
}
type CreatePortRuleR struct {
	Port     int    `json:"port" binding:"required"`
	Strategy int    `json:"strategy" binding:"required"`
	SourceIp string `json:"source_ip"`
	Ps       string `json:"ps"`
	Protocol string `json:"protocol" binding:"required"`
}

type FirewallListR struct {
	Query string `json:"query"`
	Limit int    `json:"limit" binding:"required"`
	Page  int    `json:"page" binding:"required"`
}
type UpdatePortRuleR struct {
	ID       int64  `json:"id" form:"id" binding:"required"`
	Port     string `json:"port" form:"port" binding:"required"`
	Strategy int    `json:"strategy" form:"strategy" binding:"required"`
	SourceIp string `json:"source_ip" form:"source_ip"`
	Protocol string `json:"protocol" form:"protocol" binding:"required"`
	Ps       string `json:"ps" form:"ps"`
}

type BatchCreateIPRuleR struct {
	Rules []*CreateIPRuleR `json:"rules" binding:"required"`
}

type CreateIPRuleR struct {
	SourceIp string `json:"source_ip" binding:"required"`
	Strategy int    `json:"strategy" binding:"required"`
	Ps       string `json:"ps"`
}

type BatchDeleteIPRuleR struct {
	Ids []int64 `json:"ids"  binding:"required"`
}

type UpdateIPRuleR struct {
	ID       int64  `json:"id" form:"id" binding:"required"`
	SourceIp string `json:"source_ip" form:"source_ip" binding:"required"`
	Strategy int    `json:"strategy" form:"strategy" binding:"required"`
	Ps       string `json:"ps" form:"ps"`
}

type BatchCreateForwardRuleR struct {
	Rules []*CreateForwardRuleR `json:"rules" binding:"required"`
}

type CreateForwardRuleR struct {
	SourcePort int64  `json:"source_port" form:"source_port"`
	TargetIp   string `json:"target_ip" form:"target_ip"`
	TargetPort int64  `json:"target_port" form:"target_port"`
	Protocol   string `json:"protocol" form:"protocol"`
	Ps         string `json:"ps" form:"ps"`
}

type BatchDeleteForwardRuleR struct {
	Ids []int64 `json:"ids"  binding:"required"`
}

type UpdateForwardRuleR struct {
	ID         int64  `json:"id" form:"id" binding:"required"`
	SourcePort int64  `json:"source_port" form:"source_port"`
	TargetIp   string `json:"target_ip" form:"target_ip"`
	TargetPort int64  `json:"target_port" form:"target_port"`
	Protocol   string `json:"protocol" form:"protocol"`
	Ps         string `json:"ps" form:"ps"`
}

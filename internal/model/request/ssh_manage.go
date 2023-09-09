package request

type SetSSHStatusR struct {
	Action string `json:"action" form:"action" binding:"required,oneof=stop start"`
}

type OperateSSHKeyLoginR struct {
	Action  string `json:"action" form:"action" binding:"required,oneof=on off"`
	KeyType string `json:"key_type" form:"key_type"`
}

type OperatePasswordLoginR struct {
	Action string `json:"action" form:"action" binding:"required,oneof=on off"`
}

type GetSSHLoginStatisticsR struct {
	Refresh bool `json:"refresh" form:"refresh"`
}

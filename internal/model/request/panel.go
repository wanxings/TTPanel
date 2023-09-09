package request

type OperateR struct {
	Action string `json:"action" form:"action" binding:"required,oneof=stop restart"`
}

type UpdateR struct {
	Version string `json:"version" form:"version" binding:"required"`
}

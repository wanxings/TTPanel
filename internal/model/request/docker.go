package request

type DockerSetStatusR struct {
	Action string `json:"action" form:"action" binding:"required"`
}

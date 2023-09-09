package request

type QueueTaskListR struct {
	Status int `json:"status" form:"status"`
	Limit  int `json:"limit" form:"limit" binding:"required"`
	Page   int `json:"page" form:"page" binding:"required"`
}

type QueueTaskDelR struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}

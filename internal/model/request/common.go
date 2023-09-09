package request

type BatchIDR struct {
	IDs []int64 `json:"ids" binding:"required"`
}

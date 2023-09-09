package request

type OperationLogListR struct {
	Query      string `json:"query"`
	QueryField string `json:"query_field"`
	Type       int    `json:"type"`
	Limit      int    `json:"limit" binding:"required"`
	Page       int    `json:"page" binding:"required"`
}

type SSHLoginLogListR struct {
	Query  string `json:"query"`
	Status int    `json:"status"`
	Limit  int    `json:"limit"  binding:"required"`
	Page   int    `json:"page"  binding:"required"`
}

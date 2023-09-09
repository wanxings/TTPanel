package request

type SetRecycleBinStatusR struct {
	ExplorerStatus bool `json:"explorer_status" form:"explorer_status"`
	DatabaseStatus bool `json:"database_status" form:"database_status"`
}

type RecoveryFileR struct {
	Hash  string `json:"hash" form:"hash" binding:"required"`
	Cover bool   `json:"cover" form:"cover"`
}

type DeleteRecoveryFileR struct {
	Hash string `json:"hash" form:"hash" binding:"required"`
}

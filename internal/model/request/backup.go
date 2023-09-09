package request

type BackupDatabaseR struct {
	Id            int64 `json:"id" binding:"required"`
	StorageId     int64 `json:"storage_id"`
	KeepLocalFile int   `json:"keep_local_file"`
}

type BackupProjectR struct {
	Id             int64    `json:"id" binding:"required"`
	StorageId      int64    `json:"storage_id"`
	KeepLocalFile  int      `json:"keep_local_file"`
	ExclusionRules []string `json:"exclusion_rules"`
}

type BackupDirR struct {
	Path           string   `json:"path" binding:"required"`
	StorageId      int64    `json:"storage_id"`
	KeepLocalFile  int      `json:"keep_local_file"`
	ExclusionRules []string `json:"exclusion_rules"`
	Description    string   `json:"description"`
}

type BackupPanelR struct {
	StorageId     int64 `json:"storage_id"`
	KeepLocalFile int   `json:"keep_local_file"`
}

type BackupListR struct {
	Category int   `json:"category" binding:"required"`
	Pid      int64 `json:"pid"`
	Page     int   `json:"page"  binding:"required"`
	Limit    int   `json:"limit"  binding:"required"`
}

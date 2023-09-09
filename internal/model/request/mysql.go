package request

type CreateMysqlR struct {
	DatabaseName     string `form:"database_name" json:"database_name" binding:"required"`
	UserName         string `form:"user_name" json:"user_name" binding:"required"`
	Password         string `form:"password" json:"password" binding:"required"`
	Coding           string `form:"coding" json:"coding" binding:"required"`
	AccessPermission string `form:"access_permission" json:"access_permission"`
	Ps               string `form:"ps" json:"ps"`
	Sid              int64  `form:"sid" json:"sid"`
	Pid              int64  `form:"pid" json:"pid"`
}

type ListMysqlR struct {
	Query string `form:"query" json:"query"`
	Limit int    `form:"limit" json:"limit" binding:"required"`
	Page  int    `form:"page" json:"page" binding:"required"`
	Sid   int    `form:"sid" json:"sid"`
}

type ListMysqlServerR struct {
	DBType string `form:"db_type" json:"db_type"`
}

type SetRootPwdR struct {
	Password string `form:"password" json:"password" binding:"required"`
	Sid      int64  `form:"sid" json:"sid"`
}

type SyncGetDBR struct {
	Sid int64 `form:"sid" json:"sid"`
}

type SyncToDBR struct {
	Ids []int64 `form:"ids" json:"ids"`
}

type SetAccessPermissionR struct {
	AccessPermission string `form:"access_permission" json:"access_permission" binding:"required"`
	UserName         string `form:"user_name" json:"user_name" binding:"required"`
}

type GetAccessPermissionR struct {
	UserName string `form:"user_name" json:"user_name" binding:"required"`
}

type SetPwdR struct {
	UserName string `form:"user_name" json:"user_name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type CheckDeleteDatabaseR struct {
	IDs []int64 `form:"ids" json:"ids" binding:"required"`
}

type DeleteDatabaseR struct {
	ID int64 `form:"id" json:"id" binding:"required"`
}

type ImportDatabaseR struct {
	FilePath string `json:"file_path"`
	IsAsync  bool   `json:"is_async"`
	BackupId int64  `json:"backup_id"`
	Id       int64  `json:"id"  binding:"required"`
}

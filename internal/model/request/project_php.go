package request

type CreatePHPProjectR struct {
	MoreDomains []*DomainItem `json:"more_domains" from:"more_domains"`
	Domain      *DomainItem   `json:"domain" from:"domain" binding:"required"`
	DataBase    *DataBase     `json:"database" from:"database"`
	Name        string        `json:"name" from:"name" binding:"required"`
	Ps          string        `json:"ps" from:"ps"`
	Path        string        `json:"path" from:"path" binding:"required"`
	TypeId      int           `json:"type_id" from:"type_id"`
	PHPVersion  string        `json:"php_version" from:"php_version" binding:"required"`
}
type DataBase struct {
	Sql      string `json:"sql" from:"sql" `
	DBName   string `json:"db_name" from:"db_name"`
	User     string `json:"user" from:"user"`
	Password string `json:"password" from:"password"`
	Coding   string `json:"coding" from:"coding"`
}

type PHPProjectListR struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
	Type  int    `json:"type"`
	Page  int    `json:"page"`
}

type SwitchUsingPHPVersionR struct {
	ProjectId int64  `json:"project_id"`
	Version   string `json:"version"`
	Customize string `json:"customize"`
}

type DeletePHPProjectR struct {
	ProjectId     int64 `json:"project_id" binding:"required"`
	ClearPath     bool  `json:"clear_path"`
	ClearNginxLog bool  `json:"clear_nginx_log"`
}

type SetPHPProjectStatusR struct {
	ProjectId int64  `json:"project_id" binding:"required"`
	Action    string `json:"action" binding:"required,oneof=stop start"`
}

type SetPHPProjectRunPathR struct {
	ProjectId int64  `json:"project_id" binding:"required"`
	Path      string `json:"path" binding:"required"`
}

type SetPHPProjectPathR struct {
	ProjectId int64  `json:"project_id" binding:"required"`
	Path      string `json:"path" binding:"required"`
}

type SetPHPProjectUserIniR struct {
	ProjectId int64 `json:"project_id" binding:"required"`
	Status    bool  `json:"status"`
}

type SetPHPProjectExpireTimeR struct {
	ProjectId  int64 `json:"project_id" binding:"required"`
	ExpireTime int64 `json:"expire_time" binding:"required"`
}

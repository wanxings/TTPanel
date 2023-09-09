package request

type CreateGeneralProjectR struct {
	Path           string `json:"path" binding:"required"`
	AutoCreatePath bool   `json:"auto_create_path"`
	Name           string `json:"name" binding:"required"`
	Port           int    `json:"port"`
	Command        string `json:"command" binding:"required"`
	RunUser        string `json:"run_user"`
	IsPowerOn      bool   `json:"is_power_on"`
	Description    string `json:"description"`
}

type GeneralProjectListR struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
	Type  int    `json:"type"`
	Page  int    `json:"page"`
}

type DeleteGeneralProjectR struct {
	ProjectId   int64 `json:"project_id" binding:"required"`
	ClearPath   bool  `json:"clear_path"`
	ClearRunLog bool  `json:"clear_run_log"`
}

type Domain struct {
	Name string `json:"name" from:"name" binding:"required"`
	Port int    `json:"port" from:"port" binding:"required"`
}

type DomainListR struct {
	ProjectId int64 `json:"project_id" from:"project_id" binding:"required"`
}

type ProjectInfoR struct {
	ProjectId int64 `json:"project_id" from:"project_id" binding:"required"`
}

type SaveProjectConfigR struct {
	ProjectId      int64  `json:"project_id" binding:"required"`
	Path           string `json:"path" binding:"required"`
	AutoCreatePath bool   `json:"auto_create_path"`
	Port           int    `json:"port" binding:"required"`
	Command        string `json:"command" binding:"required"`
	RunUser        string `json:"run_user" binding:"required"`
	IsPowerOn      bool   `json:"is_power_on"`
	Description    string `json:"description"`
}

type AddDomainR struct {
	ProjectId int64    `json:"project_id" from:"project_id" binding:"required"`
	Domains   []Domain `json:"domains" from:"domains" binding:"required"`
}

type SetStatusR struct {
	ProjectId int64  `json:"project_id" from:"project_id" binding:"required"`
	Action    string `json:"action" from:"action" binding:"required,oneof=start stop restart"`
}

type SetBindExtranetR struct {
	ProjectId int64 `json:"project_id" from:"project_id" binding:"required"`
	Status    bool  `json:"status" from:"status"`
}

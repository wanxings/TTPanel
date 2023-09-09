package request

type AddDomainsR struct {
	Domains   []*DomainItem `json:"domains" from:"domains"`
	ProjectID int64         `json:"project_id" from:"project_id" binding:"required"`
}

type SaveNginxConfR struct {
	ProjectId int64  `json:"project_id" binding:"required"`
	ConfBody  string `json:"conf_body" binding:"required"`
}

type SaveDefaultIndexR struct {
	ProjectId int64  `json:"project_id" from:"project_id" binding:"required"`
	Index     string `json:"index" from:"index" binding:"required"`
}

type BatchDeleteDomainR struct {
	ProjectId int64   `json:"project_id" from:"project_id" binding:"required"`
	IDS       []int64 `json:"ids" from:"ids" binding:"required"`
}

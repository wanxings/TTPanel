package request

type CreateDnsAccountR struct {
	Type          string            `json:"type" binding:"required"`
	Name          string            `json:"name" binding:"required"`
	Authorization map[string]string `json:"authorization" binding:"required"`
}

type CreateAcmeAccountR struct {
	Email string `json:"email" binding:"required"`
}

type ApplyCertificateR struct {
	Domains     []string `json:"domains" binding:"required"`
	AcmeAccount string   `json:"acme_account" binding:"required"`
	VerifyMode  int      `json:"verify_mode" binding:"required"`
	ProjectId   int64    `json:"project_id"`
	DnsAccount  string   `json:"dns_account"`
	SetUp       bool     `json:"set_up"`
}

type GetResolveR struct {
	Domains     []string `json:"domains" binding:"required"`
	AcmeAccount string   `json:"acme_account" binding:"required"`
}

type CheckDnsRecordsR struct {
	List []string `json:"list" binding:"required"`
}

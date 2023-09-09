package request

type CreateCategoryR struct {
	Name string `json:"name" from:"name" binding:"required"`
	Ps   string `json:"ps" from:"ps"`
}

type EditCategoryR struct {
	ID   int64  `json:"id" from:"id" binding:"required"`
	Name string `json:"name" from:"name" binding:"required"`
	Ps   string `json:"ps" from:"ps"`
}

type SetPsR struct {
	IDR
	Ps string `json:"ps" from:"ps"`
}

type DomainItem struct {
	Name string `json:"name" from:"name" binding:"required"`
	Port int    `json:"port" from:"port" binding:"required"`
}

type IDR struct {
	ProjectID int64 `json:"project_id" from:"project_id" binding:"required"`
}

type SetSslR struct {
	ProjectId int64  `json:"project_id" binding:"required"`
	SSLKey    string `json:"ssl_key"`
	Key       string `json:"key"`
	Csr       string `json:"csr"`
}

type AlwaysUseHttpsR struct {
	ProjectId int64 `json:"project_id" binding:"required"`
	Action    bool  `json:"action" binding:"required"`
}

type CreateRedirectR struct {
	ProjectId       int64    `json:"project_id"  binding:"required"`
	Status          bool     `json:"status"`
	PreserveUriArgs bool     `json:"preserve_uri_args"`
	Type            int      `json:"type"  binding:"required"`
	Code            int      `json:"code"  binding:"required,oneof=301 302"`
	Path            string   `json:"path"`
	Domains         []string `json:"domains"`
	TargetUrl       string   `json:"target_url"  binding:"required"`
}

type BatchEditRedirectR struct {
	ProjectID int64                      `json:"project_id" from:"project_id" binding:"required"`
	List      map[string]CreateRedirectR `json:"list"  binding:"required"`
}

type BatchDeleteRedirectR struct {
	ProjectId int64    `json:"project_id" binding:"required"`
	Keys      []string `json:"keys" binding:"required"`
}

type CreateAntiLeechConfigR struct {
	ProjectID          int64    `json:"project_id" binding:"required"`
	SuffixList         []string `json:"suffix_list" binding:"required"`
	PassDomains        []string `json:"pass_domains" binding:"required"`
	ResponseStatusCode int      `json:"response_status_code" binding:"required"`
	RefererNone        bool     `json:"referer_none"`
}

type CreateReverseProxyConfigR struct {
	ProjectId  int64  `json:"project_id" binding:"required"`
	Status     bool   `json:"status"`
	Name       string `json:"name" binding:"required"`
	ExactMatch bool   `json:"exact_match"`
	CacheTime  int    `json:"cache_time"`
	ProxyDir   string `json:"proxy_dir" binding:"required"`
	TargetUrl  string `json:"target_url"`
	Host       string `json:"host"`

	ProxyConnectTimeout int    `json:"proxy_connect_timeout"`
	ProxyReadTimeout    int    `json:"proxy_read_timeout"`
	ProxySendTimeout    int    `json:"proxy_send_timeout"`
	ProxyHttpVersion    string `json:"proxy_http_version"`

	SubFilters []struct {
		Old string `json:"old"`
		New string `json:"new"`
	} `json:"sub_filters"`
	SubFilterTypes []string `json:"sub_filter_types"`
	SubFilterOnce  string   `json:"sub_filter_once"`
}

type BatchEditReverseProxyConfigR struct {
	ProjectID int64                                `json:"project_id" from:"project_id" binding:"required"`
	List      map[string]CreateReverseProxyConfigR `json:"list"  binding:"required"`
}

type BatchDeleteReverseProxyConfigR struct {
	ProjectId int64    `json:"project_id" binding:"required"`
	Keys      []string `json:"keys" binding:"required"`
}

type CreateAccessRuleConfigR struct {
	ProjectId       int64  `json:"project_id" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Dir             string `json:"dir" binding:"required"`
	RuleType        int    `json:"rule_type" binding:"required"`
	BasicAuthConfig struct {
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"basic_auth_config"`
	NoAccessConfig struct {
		SuffixList []string `json:"suffix_list"`
	} `json:"no_access_config"`
}

type EditAccessRuleConfigR struct {
	ProjectId int64                   `json:"project_id" binding:"required"`
	Key       string                  `json:"key" binding:"required"`
	Config    CreateAccessRuleConfigR `json:"config"`
}
type BatchDeleteAccessRuleConfigR struct {
	ProjectId int64    `json:"project_id" binding:"required"`
	Keys      []string `json:"keys" binding:"required"`
}

package response

type TTWafConfig struct {
	//默认配置
	Cc               TTWafCCSetting               `json:"cc" binding:"required"`
	Analytics        TTWafAnalyticsSetting        `json:"analytics" binding:"required"`
	SemanticAnalysis TTWafSemanticAnalysisSetting `json:"semantic_analysis"`
	Get              TTWafGetSetting              `json:"get" binding:"required"`
	Post             TTWafPostSetting             `json:"post" binding:"required"`
	UserAgent        TTWafUserAgentSetting        `json:"user_agent" binding:"required"`
	Cookie           TTWafCookieSetting           `json:"cookie" binding:"required"`
	BlockCountry     TTWafBlockCountrySetting     `json:"block_country" binding:"required"`
	MethodType       []TTWafMethodTypeSetting     `json:"method_type" binding:"required"`
	HeaderLen        []TTWafHeaderLenSetting      `json:"header_len" binding:"required"`
	//防火墙配置
	Status      bool `json:"status"`
	IpWhiteList struct {
		IPV4 [][]int64 `json:"ipv4"`
		IPV6 []string  `json:"ipv6"`
	} `json:"ip_white_list"`
	IpBlackList struct {
		IPV4 [][]int64 `json:"ipv4"`
		IPV6 []string  `json:"ipv6"`
	} `json:"ip_black_list"`
	UrlWhiteList    []string                    `json:"url_white_list"`
	UrlBlackList    []string                    `json:"url_black_list"`
	AttackTolerance TTWafAttackToleranceSetting `json:"attack_tolerance" binding:"required"`
	PostArgsLimit   int                         `json:"post_args_limit"`
	CdnHeaders      []string                    `json:"cdn_headers"`
	LogsPath        string                      `json:"logs_path"`
	LogSave         int                         `json:"log_save"`
	SavePostData    bool                        `json:"save_post_data"`
}

type TTWafProjectConfig struct {
	Status           bool                         `json:"status"`
	DisablePhpPath   []string                     `json:"disable_php_path"`
	DisablePath      []string                     `json:"disable_path"`
	DisableExt       []string                     `json:"disable_ext"`
	DisableUploadExt []string                     `json:"disable_upload_ext"`
	DisableRule      TTWafDisableRuleSetting      `json:"disable_rule"  binding:"required"`
	Cc               TTWafCCSetting               `json:"cc"  binding:"required"`
	Analytics        TTWafAnalyticsSetting        `json:"analytics" binding:"required"`
	SemanticAnalysis TTWafSemanticAnalysisSetting `json:"semantic_analysis"  binding:"required"`
	BlockCountry     TTWafBlockCountrySetting     `json:"block_country"  binding:"required"`
	Get              bool                         `json:"get"`
	Post             bool                         `json:"post"`
	Cookie           bool                         `json:"cookie"`
	UserAgent        bool                         `json:"user_agent"`
	CDN              bool                         `json:"cdn"`
	AttackTolerance  TTWafAttackToleranceSetting  `json:"attack_tolerance"  binding:"required"`
}

type TTWafCCSetting struct {
	RespCode    int      `json:"resp_code"`
	Type        int      `json:"type"`
	Description string   `json:"description"`
	Increase    bool     `json:"increase"`
	Limit       int      `json:"limit"`
	BanTime     int      `json:"ban_time"`
	Status      bool     `json:"status"`
	RespFile    string   `json:"resp_file"`
	Cycle       int      `json:"cycle"`
	Country     []string `json:"country"`
}
type TTWafDomainSetting struct {
	Path       string   `json:"path"`
	DomainList []string `json:"domain_list"`
}

type TTWafAnalyticsSetting struct {
	Status        bool          `json:"status"`
	LogSave       int           `json:"log_save"`
	ExcludeSuffix []string      `json:"exclude_suffix"`
	ExcludeCode   []int         `json:"exclude_code"`
	ExcludePath   []interface{} `json:"exclude_path"`
	ExcludeIp     struct {
		Ipv4 [][]int64     `json:"ipv4"`
		Ipv6 []interface{} `json:"ipv6"`
	} `json:"exclude_ip"`
}
type TTWafSemanticAnalysisSetting struct {
	GetSql  bool `json:"get_sql"`
	GetXss  bool `json:"get_xss"`
	PostSql bool `json:"post_sql"`
	PostXss bool `json:"post_xss"`
}

type TTWafGetSetting struct {
	RespCode    int    `json:"resp_code"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	RespFile    string `json:"resp_file"`
}

type TTWafPostSetting struct {
	RespCode    int    `json:"resp_code"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	RespFile    string `json:"resp_file"`
}

type TTWafUserAgentSetting struct {
	RespCode    int    `json:"resp_code"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	RespFile    string `json:"resp_file"`
}

type TTWafCookieSetting struct {
	RespCode    int    `json:"resp_code"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	RespFile    string `json:"resp_file"`
}

type TTWafBlockCountrySetting struct {
	RespCode    int      `json:"resp_code"`
	Description string   `json:"description"`
	Reverse     bool     `json:"reverse"`
	Status      bool     `json:"status"`
	List        []string `json:"list"`
	RespFile    string   `json:"resp_file"`
}

type TTWafMethodTypeSetting struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

type TTWafHeaderLenSetting struct {
	Name string `json:"name"`
	Len  int    `json:"len"`
}

type TTWafAttackToleranceSetting struct {
	TimeWindow  int `json:"time_window"`
	AttackLimit int `json:"attack_limit"`
	BanTime     int `json:"ban_time"`
}

type TTWafDisableRuleSetting struct {
	Uri       []int `json:"uri"`
	PostArgs  []int `json:"post_args"`
	GetArgs   []int `json:"get_args"`
	Cookie    []int `json:"cookie"`
	UserAgent []int `json:"user_agent"`
}

type TTWafRegRule struct {
	Status      bool   `json:"status"`
	Reg         string `json:"reg"`
	Description string `json:"description"`
}

// 总ip数，总请求数，总流量大小
type TTWafAnalyticsOverview struct {
	TotalIP      int `json:"total_ip"`
	TotalReq     int `json:"total_req"`
	TotalTraffic int `json:"total_traffic"`
}

type TTWafAccessLog struct {
	StatusCode     int    `json:"status_code"`
	BodyLength     int    `json:"body_length"`
	ServerProtocol string `json:"server_protocol"`
	Method         string `json:"method"`
	XForwardedFor  string `json:"x_forwarded_for"`
	Uri            string `json:"uri"`
	RequestTime    int    `json:"request_time"`
	HeadersBody    string `json:"headers_body"`
	CreateTime     int    `json:"create_time"`
	UserAgent      string `json:"user_agent"`
	Domain         string `json:"domain"`
	ServerName     string `json:"server_name"`
	Ip             string `json:"ip"`
}

type TTWafAccessTotal struct {
	TotalRequests    int            `json:"total_requests"`
	TotalBodyLength  int            `json:"total_body_length"`
	UriStatistics    map[string]int `json:"uri_statistics"`
	SpiderStatistics map[string]int `json:"spider_statistics"`
	DeviceStatistics map[string]int `json:"device_statistics"`
	IPList           map[string]int `json:"ip_list"`
}

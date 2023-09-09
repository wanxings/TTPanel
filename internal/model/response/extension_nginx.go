package response

type NginxLoadStatusP struct {
	Accepts   string  `json:"accepts"`
	Handled   string  `json:"handled"`
	Requests  string  `json:"requests"`
	Reading   string  `json:"Reading"`
	Writing   string  `json:"Writing"`
	Waiting   string  `json:"Waiting"`
	Active    string  `json:"active"`
	Worker    int     `json:"worker"`
	WorkerCpu float64 `json:"worker_cpu"`
	WorkerMen int     `json:"worker_men"`
}

type NginxGeneralConfigP struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
	Ps    string `json:"ps"`
	From  string `json:"from"`
}

var DefaultNginxGeneralConfigs = []*NginxGeneralConfigP{
	{
		Name:  "worker_processes",
		Value: "",
		Unit:  "",
		Ps:    "处理进程,auto表示自动,数字表示进程数",
		From:  "config",
	},
	{
		Name:  "worker_connections",
		Value: "",
		Unit:  "",
		Ps:    "最大并发链接数",
		From:  "config",
	},
	{
		Name:  "keepalive_timeout",
		Value: "",
		Unit:  "",
		Ps:    "连接超时时间",
		From:  "config",
	},
	{
		Name:  "gzip",
		Value: "",
		Unit:  "",
		Ps:    "是否开启压缩传输",
		From:  "config",
	},
	{
		Name:  "gzip_min_length",
		Value: "",
		Unit:  "",
		Ps:    "最小压缩文件",
		From:  "config",
	},
	{
		Name:  "gzip_comp_level",
		Value: "",
		Unit:  "",
		Ps:    "压缩率",
		From:  "config",
	},
	{
		Name:  "client_max_body_size",
		Value: "",
		Unit:  "",
		Ps:    "最大上传文件",
		From:  "config",
	},
	{
		Name:  "server_names_hash_bucket_size",
		Value: "",
		Unit:  "",
		Ps:    "服务器名字的hash表大小",
		From:  "config",
	},
	{
		Name:  "client_header_buffer_size",
		Value: "",
		Unit:  "",
		Ps:    "客户端请求头buffer大小",
		From:  "config",
	},
	{
		Name:  "client_body_buffer_size",
		Value: "",
		Unit:  "",
		Ps:    "请求主体缓冲区",
		From:  "proxyConfig",
	},
}

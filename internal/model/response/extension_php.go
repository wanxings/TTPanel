package response

type LibP struct {
	Name       string   `json:"name"`
	Versions   []string `json:"versions"`
	Type       string   `json:"type"`
	Msg        string   `json:"msg"`
	Shell      string   `json:"shell"`
	Check      string   `json:"check"`
	TaskStatus int      `json:"task_status"`
	Status     bool     `json:"status"`
}

type PHPLoadStatusP struct {
	Pool               string `json:"pool"`
	ProcessManager     string `json:"process manager"`
	StartTime          int    `json:"start time"`
	StartSince         int    `json:"start since"`
	AcceptedConn       int    `json:"accepted conn"`
	ListenQueue        int    `json:"listen queue"`
	MaxListenQueue     int    `json:"max listen queue"`
	ListenQueueLen     int    `json:"listen queue len"`
	IdleProcesses      int    `json:"idle processes"`
	ActiveProcesses    int    `json:"active processes"`
	TotalProcesses     int    `json:"total processes"`
	MaxActiveProcesses int    `json:"max active processes"`
	MaxChildrenReached int    `json:"max children reached"`
	SlowRequests       int    `json:"slow requests"`
}

type PHPPerformanceConfigP struct {
	MaxChildren     int    `json:"max_children"`
	StartServers    int    `json:"start_servers"`
	MinSpareServers int    `json:"min_spare_servers"`
	MaxSpareServers int    `json:"max_spare_servers"`
	Pm              string `json:"pm"`
	Allowed         string `json:"allowed"`
	Unix            string `json:"unix"`
	Port            int    `json:"port"`
	Bind            string `json:"bind"`
}

type PHPGeneralConfigP struct {
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Value string `json:"value"`
	Ps    string `json:"ps"`
}

var DefaultPHPGeneralConfigs = []PHPGeneralConfigP{
	{
		Name:  "short_open_tag",
		Type:  1,
		Value: "",
		Ps:    "短标签支持",
	},
	{
		Name:  "asp_tags",
		Type:  1,
		Value: "",
		Ps:    "ASP标签",
	},
	{
		Name:  "display_errors",
		Type:  1,
		Value: "",
		Ps:    "是否输出详细错误信息",
	},
	{
		Name:  "cgi.fix_pathinfo",
		Type:  0,
		Value: "",
		Ps:    "是否开启pathinfo",
	},
	{
		Name:  "date.timezone",
		Type:  3,
		Value: "",
		Ps:    "时区",
	},
	{
		Name:  "file_uploads",
		Type:  1,
		Value: "",
		Ps:    "文件上传",
	},
	{
		Name:  "memory_limit",
		Type:  2,
		Value: "",
		Ps:    "内存限制",
	},
	{
		Name:  "post_max_size",
		Type:  2,
		Value: "",
		Ps:    "POST提交最大限制",
	},
	{
		Name:  "upload_max_filesize",
		Type:  2,
		Value: "",
		Ps:    "上传文件最大限制",
	},
	{
		Name:  "max_execution_time",
		Type:  2,
		Value: "",
		Ps:    "最大脚本运行时间",
	},
	{
		Name:  "max_input_time",
		Type:  2,
		Value: "",
		Ps:    "最大输入时间",
	},
	{
		Name:  "max_file_uploads",
		Type:  2,
		Value: "",
		Ps:    "允许同时上传文件的最大数量",
	},
	{
		Name:  "default_socket_timeout",
		Type:  2,
		Value: "",
		Ps:    "默认Socket超时时间",
	},
	{
		Name:  "error_reporting",
		Type:  3,
		Value: "",
		Ps:    "错误报告级别",
	},
}

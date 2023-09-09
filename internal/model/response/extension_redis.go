package response

type RedisGeneralConfigP struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
	Ps    string `json:"ps"`
	From  string `json:"from"`
}

var DefaultRedisGeneralConfigs = []*RedisGeneralConfigP{
	{
		Name:  "bind",
		Value: "",
		Unit:  "",
		Ps:    "绑定IP(修改绑定IP可能会存在安全隐患)",
		From:  "config",
	},
	{
		Name:  "port",
		Value: "",
		Unit:  "",
		Ps:    "绑定端口",
		From:  "config",
	},
	{
		Name:  "timeout",
		Value: "",
		Unit:  "",
		Ps:    "空闲连接超时时间,0表示不断开",
		From:  "config",
	},
	{
		Name:  "maxclients",
		Value: "10000",
		Unit:  "",
		Ps:    "最大连接数",
		From:  "config",
	},
	{
		Name:  "databases",
		Value: "",
		Unit:  "",
		Ps:    "数据库数量",
		From:  "config",
	},
	{
		Name:  "requirepass",
		Value: "",
		Unit:  "",
		Ps:    "redis密码,留空代表没有设置密码",
		From:  "config",
	},
	{
		Name:  "maxmemory",
		Value: "0",
		Unit:  "",
		Ps:    "MB,最大使用内存，0表示不限制",
		From:  "config",
	},
}

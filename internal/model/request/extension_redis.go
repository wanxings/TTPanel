package request

type RedisSetStatusR struct {
	Action string `json:"action" binding:"required,oneof=start stop restart reload"`
}

type RedisSavePerformanceConfigR struct {
	Bind        string `json:"bind" form:"bind" binding:"required"`
	Port        int    `json:"port" form:"port" binding:"required"`
	Timeout     int    `json:"timeout" form:"timeout"`
	MaxClients  int    `json:"maxclients" form:"maxclients" binding:"required"`
	Databases   int    `json:"databases" form:"databases" binding:"required"`
	RequirePass string `json:"requirepass" form:"requirepass"`
	MaxMemory   int    `json:"maxmemory" form:"maxmemory"`
}

type RedisSavePersistentConfigR struct {
	Aof *AofConfig   `json:"aof"`
	Dir string       `json:"dir"`
	Rdb []*RdbConfig `json:"rdb"`
}

type AofConfig struct {
	AppendFsync string `json:"appendfsync"`
	AppendOnly  string `json:"appendonly"`
}

type RdbConfig struct {
	Changes int `json:"changes"`
	Time    int `json:"time"`
}

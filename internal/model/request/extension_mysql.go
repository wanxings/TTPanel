package request

type MysqlSetStatusR struct {
	Action string `form:"action" json:"action" binding:"required,oneof=stop start restart reload"`
}
type MysqlSetPerformanceConfigR struct {
	BinlogCacheSize      int     `json:"binlog_cache_size"`
	InnodbBufferPoolSize int     `json:"innodb_buffer_pool_size"`
	InnodbLogBufferSize  int     `json:"innodb_log_buffer_size"`
	JoinBufferSize       int     `json:"join_buffer_size"`
	KeyBufferSize        int     `json:"key_buffer_size"`
	MaxConnections       int     `json:"max_connections"`
	MemSize              float64 `json:"memSize"`
	MysqlSet             int     `json:"mysql_set"`
	QueryCacheSize       int     `json:"query_cache_size"`
	QueryCacheType       int     `json:"query_cache_type"`
	ReadBufferSize       int     `json:"read_buffer_size"`
	ReadRndBufferSize    int     `json:"read_rnd_buffer_size"`
	SortBufferSize       int     `json:"sort_buffer_size"`
	TableOpenCache       int     `json:"table_open_cache"`
	ThreadCacheSize      int     `json:"thread_cache_size"`
	ThreadStack          int     `json:"thread_stack"`
	TmpTableSize         int     `json:"tmp_table_size"`
	MaxHeapTableSize     int     `json:"max_heap_table_size"`
}

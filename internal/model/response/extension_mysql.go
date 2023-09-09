package response

type MysqlLoadStatusP struct {
	AbortedClients               string `json:"Aborted_clients"`
	AbortedConnects              string `json:"Aborted_connects"`
	BytesReceived                string `json:"Bytes_received"`
	BytesSent                    string `json:"Bytes_sent"`
	ComCommit                    string `json:"Com_commit"`
	ComRollback                  string `json:"Com_rollback"`
	Connections                  string `json:"Connections"`
	CreatedTmpDiskTables         string `json:"Created_tmp_disk_tables"`
	CreatedTmpTables             string `json:"Created_tmp_tables"`
	InnodbBufferPoolPagesDirty   string `json:"Innodb_buffer_pool_pages_dirty"`
	InnodbBufferPoolReadRequests string `json:"Innodb_buffer_pool_read_requests"`
	InnodbBufferPoolReads        string `json:"Innodb_buffer_pool_reads"`
	KeyReadRequests              string `json:"Key_read_requests"`
	KeyReads                     string `json:"Key_reads"`
	KeyWriteRequests             string `json:"Key_write_requests"`
	KeyWrites                    string `json:"Key_writes"`
	MaxUsedConnections           string `json:"Max_used_connections"`
	OpenTables                   string `json:"Open_tables"`
	OpenedFiles                  string `json:"Opened_files"`
	OpenedTables                 string `json:"Opened_tables"`
	QCacheHits                   string `json:"Qcache_hits"`
	QCacheInserts                string `json:"Qcache_inserts"`
	Questions                    string `json:"Questions"`
	SelectFullJoin               string `json:"Select_full_join"`
	SelectRangeCheck             string `json:"Select_range_check"`
	SortMergePasses              string `json:"Sort_merge_passes"`
	TableLocksWaited             string `json:"Table_locks_waited"`
	ThreadsCached                string `json:"Threads_cached"`
	ThreadsConnected             string `json:"Threads_connected"`
	ThreadsCreated               string `json:"Threads_created"`
	ThreadsRunning               string `json:"Threads_running"`
	Uptime                       string `json:"Uptime"`
	Run                          string `json:"Run"`
	File                         string `json:"File"`
	Position                     string `json:"Position"`
}
type T struct {
	BinlogCacheSize      string `json:"binlog_cache_size"`
	InnodbBufferPoolSize string `json:"innodb_buffer_pool_size"`
	InnodbLogBufferSize  string `json:"innodb_log_buffer_size"`
	JoinBufferSize       string `json:"join_buffer_size"`
	KeyBufferSize        string `json:"key_buffer_size"`
	MaxConnections       string `json:"max_connections"`
	MaxHeapTableSize     string `json:"max_heap_table_size"`
	QueryCacheSize       string `json:"query_cache_size"`
	QueryCacheType       string `json:"query_cache_type"`
	ReadBufferSize       string `json:"read_buffer_size"`
	ReadRndBufferSize    string `json:"read_rnd_buffer_size"`
	SortBufferSize       string `json:"sort_buffer_size"`
	TableOpenCache       string `json:"table_open_cache"`
	ThreadCacheSize      string `json:"thread_cache_size"`
	ThreadStack          string `json:"thread_stack"`
	TmpTableSize         string `json:"tmp_table_size"`
}
type MysqlGeneralConfigP struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
	Ps    string `json:"ps"`
	From  string `json:"from"`
}

var DefaultMysqlGeneralConfigs = []*MysqlGeneralConfigP{
	{
		Name:  "binlog_cache_size",
		Value: "",
		Unit:  "KB",
		Ps:    "* 连接数, 二进制日志缓存大小(4096的倍数)",
	},
	{
		Name:  "innodb_buffer_pool_size",
		Value: "",
		Unit:  "MB",
		Ps:    ", Innodb缓冲区大小",
	},
	{
		Name:  "innodb_log_buffer_size",
		Value: "",
		Unit:  "MB",
		Ps:    ", Innodb日志缓冲区大小",
	},
	{
		Name:  "join_buffer_size",
		Value: "",
		Unit:  "KB",
		Ps:    "* 连接数, 关联表缓存大小",
	},
	{
		Name:  "key_buffer_size",
		Value: "",
		Unit:  "MB",
		Ps:    ", 用于索引的缓冲区大小",
	},
	{
		Name:  "max_connections",
		Value: "",
		Unit:  "",
		Ps:    "最大连接数",
	},
	{
		Name:  "max_heap_table_size",
		Value: "",
		Unit:  "",
		Ps:    "max_heap_table_size",
	},
	{
		Name:  "query_cache_size",
		Value: "",
		Unit:  "MB",
		Ps:    ", 查询缓存,不开启请设为0",
	},
	{
		Name:  "query_cache_type",
		Value: "",
		Unit:  "",
		Ps:    "query_cache_type",
	},
	{
		Name:  "read_buffer_size",
		Value: "",
		Unit:  "KB",
		Ps:    "* 连接数, 读入缓冲区大小",
	},
	{
		Name:  "read_rnd_buffer_size",
		Value: "",
		Unit:  "KB",
		Ps:    "* 连接数, 随机读取缓冲区大小",
	},
	{
		Name:  "sort_buffer_size",
		Value: "",
		Unit:  "KB",
		Ps:    "* 连接数, 每个线程排序的缓冲大小",
	},
	{
		Name:  "table_open_cache",
		Value: "",
		Unit:  "",
		Ps:    "表缓存",
	},
	{
		Name:  "thread_cache_size",
		Value: "",
		Unit:  "",
		Ps:    "线程池大小",
	},
	{
		Name:  "thread_stack",
		Value: "",
		Unit:  "KB",
		Ps:    "* 连接数, 每个线程的堆栈大小",
	},
	{
		Name:  "tmp_table_size",
		Value: "",
		Unit:  "MB",
		Ps:    ", 临时表缓存大小",
	},
}

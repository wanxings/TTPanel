package request

type NginxSetStatusR struct {
	Action string `form:"action" json:"action" binding:"required,oneof=stop start restart reload"`
}

type NginxSetPerformanceConfigR struct {
	WorkerProcesses           string `json:"worker_processes"`
	WorkerConnections         string `json:"worker_connections"`
	KeepaliveTimeout          string `json:"keepalive_timeout"`
	Gzip                      string `json:"gzip"`
	GzipMinLength             string `json:"gzip_min_length"`
	GzipCompLevel             string `json:"gzip_comp_level"`
	ClientMaxBodySize         string `json:"client_max_body_size"`
	ServerNamesHashBucketSize string `json:"server_names_hash_bucket_size"`
	ClientHeaderBufferSize    string `json:"client_header_buffer_size"`
	ClientBodyBufferSize      string `json:"client_body_buffer_size"`
}

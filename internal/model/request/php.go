package request

type VersionR struct {
	Version string `json:"version" form:"version" binding:"required"`
}

type PHPSetStatusR struct {
	Version string `json:"version" form:"version" binding:"required"`
	Action  string `json:"action" form:"action" binding:"required,oneof=stop start restart reload"`
}

type InstallLibR struct {
	Version string `json:"version" form:"version" binding:"required"`
	Name    string `json:"name" form:"name" binding:"required"`
}

type SaveGeneralConfigR struct {
	Version      string `json:"version" form:"version" binding:"required"`
	ShortOpenTag string `json:"short_open_tag" form:"short_open_tag"`
	AspTags      string `json:"asp_tags" form:"asp_tags"`
	DisplayError string `json:"display_errors" form:"display_errors"`
	CgiFixPath   string `json:"cgi_fix_pathinfo" form:"cgi_fix_pathinfo"`
	Timezone     string `json:"date_timezone" form:"date_timezone"`
	FileUpload   string `json:"file_uploads" form:"file_uploads"`
	MemoryLimit  string `json:"memory_limit" form:"memory_limit"`
	PostMaxSize  string `json:"post_max_size" form:"post_max_size"`
	UploadMax    string `json:"upload_max_filesize" form:"upload_max_filesize"`
	MaxExecTime  string `json:"max_execution_time" form:"max_execution_time"`
	MaxInputTime string `json:"max_input_time" form:"max_input_time"`
	MaxFileUp    string `json:"max_file_uploads" form:"max_file_uploads"`
	DefSockTime  string `json:"default_socket_timeout" form:"default_socket_timeout"`
	ErrorReport  string `json:"error_reporting" form:"error_reporting"`
}

type AddDisableFunctionR struct {
	Version string `json:"version" form:"version" binding:"required"`
	Name    string `json:"name" form:"name" binding:"required"`
}

type SavePerformanceConfigR struct {
	Version         string `json:"version" form:"version" binding:"required"`
	Unix            string `json:"unix" form:"unix" binding:"required"`
	Bind            string `json:"bind" form:"bind" binding:"required"`
	Allowed         string `json:"allowed" form:"allowed" binding:"required"`
	Pm              string `json:"pm" form:"pm" binding:"required"`
	MaxChildren     int    `json:"max_children" form:"max_children" binding:"required"`
	StartServers    int    `json:"start_servers" form:"start_servers" binding:"required"`
	MinSpareServers int    `json:"min_spare_servers" form:"min_spare_servers" binding:"required"`
	MaxSpareServers int    `json:"max_spare_servers" form:"max_spare_servers" binding:"required"`
}

type SetConfigR struct {
	Php  string `json:"php" form:"php" binding:"required"`
	Port string `json:"port" form:"port" binding:"required"`
}

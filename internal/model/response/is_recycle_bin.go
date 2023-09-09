package response

type RecycleBinInfo struct {
	Name       string `json:"name"`
	IsDir      bool   `json:"is_dir"`
	Size       int64  `json:"size"`
	DeleteTime int64  `json:"delete_time"`
	SourcePath string `json:"source_path"`
}

package request

import "mime/multipart"

type GetDirR struct {
	Path    string `json:"path" form:"path" binding:"required"`
	Page    int    `json:"page" form:"page" binding:"required"`
	Limit   int    `json:"limit" form:"limit" binding:"required"`
	Query   string `json:"query" form:"query" `
	All     bool   `json:"all" form:"all"`
	Reverse bool   `json:"reverse" form:"reverse"`
	Sort    string `json:"sort" form:"sort"`
}
type SaveFileBodyR struct {
	Type     int    `json:"type" form:"type"`
	Path     string `json:"path" form:"path" binding:"required"`
	Body     string `json:"body" form:"body"`
	Encoding string `json:"encoding" form:"encoding" binding:"required"`
}
type BatchSetFileAccessR struct {
	PathList    []string `json:"path_list" form:"path_list" binding:"required"`
	UserName    string   `json:"user_name" form:"user_name" binding:"required"`
	Permissions string   `json:"permissions" form:"permissions" binding:"required"`
	All         bool     `json:"all" form:"all"`
}
type BatchDeleteDirFileR struct {
	PathList []string `json:"path_list" form:"path_list" binding:"required"`
}
type CommonPathR struct {
	Path string `json:"path" form:"path" binding:"required"`
}
type RenameR struct {
	OldPath string `json:"old_path" form:"old_path" binding:"required"`
	NewPath string `json:"new_path" form:"new_path" binding:"required"`
}
type BatchCheckExistsFilesR struct {
	InitPath  string   `json:"init_path" form:"init_path" binding:"required"`
	FileList  []string `json:"file_list" form:"file_list" binding:"required"`
	CheckPath string   `json:"check_path" form:"check_path" binding:"required"`
}
type BatchCopyMoveR struct {
	InitPath            string   `json:"init_path" form:"init_path" binding:"required"`
	FileList            []string `json:"file_list" form:"file_list" binding:"required"`
	ToPath              string   `json:"to_path" form:"to_path" binding:"required"`
	Action              string   `json:"action" form:"action" binding:"required"`
	FromReserveFileList []string `json:"from_reserve_file_list" form:"from_reserve_file_list"`
	ToReserveFileList   []string `json:"to_reserve_file_list" form:"to_reserve_file_list"`
	AllReserveFileList  []string `json:"all_reserve_file_list" form:"all_reserve_file_list"`
}
type UploadR struct {
	Path      string                `json:"path" form:"path" binding:"required"`
	Cover     bool                  `json:"cover" form:"cover" `
	Blob      *multipart.FileHeader `json:"blob" form:"blob" binding:"required"`
	Name      string                `json:"name" form:"name" binding:"required"`
	Size      int64                 `json:"size" form:"size" binding:"required"`
	Start     int64                 `json:"start" form:"start"`
	ChunkSize int64                 `json:"chunk_size" form:"chunk_size"`
}

type CompressR struct {
	Path         string   `json:"path" form:"path" binding:"required"`
	FileList     []string `json:"file_list" form:"file_list" binding:"required"`
	CompressPath string   `json:"compress_path" form:"compress_path" binding:"required"`
	CompressType string   `json:"compress_type" form:"compress_type" binding:"required"`
	Password     string   `json:"password" form:"password"`
}

type DecompressR struct {
	CompressFilePath  string `json:"compress_file_path" form:"compress_file_path" binding:"required"`
	DecompressDirPath string `json:"decompress_dir_path" form:"decompress_dir_path" binding:"required"`
	Password          string `json:"password" form:"password"`
	IsCover           bool   `json:"is_cover" form:"is_cover"`
}

type GetLogContentR struct {
	Path     string `json:"path" binding:"required"`
	Location string `json:"location" binding:"required,oneof=head tail"`
	Line     uint   `json:"line" binding:"required"`
}

type ClearLogContentR struct {
	Path string `json:"path" binding:"required"`
}

type SearchFileContentR struct {
	DirPath        string   `json:"dir_path" binding:"required"`
	ContainsSubdir bool     `json:"contains_subdir"`
	Suffix         []string `json:"suffix"`
	KeywordReg     []string `json:"keyword_reg"`
	KeywordNormal  []string `json:"keyword_normal"`
	CaseSensitive  bool     `json:"case_sensitive"`
	MinSize        int64    `json:"min_size"`
	MaxSize        int64    `json:"max_size" binding:"required"`
}

type RemoteDownloadR struct {
	Url      string `json:"url" binding:"required"`
	SavePath string `json:"save_path" binding:"required"`
	Replace  bool   `json:"replace"`
}

type RemoteDownloadProcessR struct {
	Key string `json:"key" binding:"required"`
}

type BatchOperateSpecialPermissionR struct {
	Action      string   `json:"action" binding:"required,oneof=+ - ="`
	PathList    []string `json:"path_list" binding:"required"`
	Permissions string   `json:"permissions" binding:"required,oneof=i a s u e"`
	All         bool     `json:"all"`
}

type SetRemarkR struct {
	Path   string `json:"path" binding:"required"`
	Remark string `json:"remark"`
}

type OperateFavoritesR struct {
	Path   string `json:"path" binding:"required"`
	Action string `json:"action" binding:"required,oneof=add delete"`
}

type GenerateDownloadExternalLinkR struct {
	Description string `json:"description"`
	FilePath    string `json:"file_path" binding:"required"`
	ExpireTime  int    `json:"expire_time" binding:"required"`
}

type DeleteDownloadExternalLinkR struct {
	ID int64 `json:"id" binding:"required"`
}

type GetFileTemporaryDownloadLinkR struct {
	FilePath string `json:"file_path" binding:"required"`
}

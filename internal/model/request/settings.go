package request

type SetBasicAuthR struct {
	Status   bool   `json:"status" form:"status"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"username"`
}

type SetPanelPortR struct {
	PanelPort int `json:"panel_port" form:"panel_port" binding:"required"`
}

type SetEntranceR struct {
	Entrance string `json:"entrance" form:"entrance" binding:"required"`
}

type SetUserR struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type SetEntranceErrorCodeR struct {
	ErrorCode int `json:"error_code" form:"error_code" binding:"required"`
}

type SetPanelNameR struct {
	PanelName string `json:"panel_name" form:"panel_name" binding:"required"`
}

type SetPanelIPR struct {
	PanelIP string `json:"panel_ip" form:"panel_ip" binding:"required"`
}

type SetDefaultWebsiteDirectoryR struct {
	DefaultWebsiteDirectory string `json:"default_website_directory" form:"default_website_directory" binding:"required"`
}

type SetDefaultBackupDirectoryR struct {
	DefaultBackupDirectory string `json:"default_backup_directory" form:"default_backup_directory" binding:"required"`
}

type SetPanelApiR struct {
	Status    bool     `json:"status" form:"status"`
	Key       string   `json:"key" form:"key" binding:"required"`
	Whitelist []string `json:"whitelist" form:"whitelist"`
}

type SetLanguageR struct {
	Language string `json:"language" form:"language" binding:"required"`
}

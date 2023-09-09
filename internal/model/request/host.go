package request

type AddHostCategoryR struct {
	Name   string `json:"name" form:"name" binding:"required"`
	Remark string `json:"remark" form:"remark"`
}

type EditHostCategoryR struct {
	ID     int64  `json:"id" form:"id" binding:"required"`
	Name   string `json:"name" form:"name" binding:"required"`
	Remark string `json:"remark" form:"remark"`
}

type DeleteHostCategoryR struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}

type DeleteHostR struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}

type AddHostR struct {
	Name       string `json:"name" form:"name" binding:"required"`
	Address    string `json:"address" form:"address" binding:"required"`
	User       string `json:"user" form:"user" binding:"required"`
	Password   string `json:"password" form:"password"`
	PrivateKey string `json:"private_key" form:"private_key"`
	Port       int    `json:"port" form:"port" binding:"required"`
	CId        int64  `json:"cid" form:"cid"`
	Remark     string `json:"remark" form:"remark"`
}

type HostListR struct {
	Query string `json:"query" form:"query"`
	Limit int    `json:"limit" form:"limit" binding:"required"`
	CId   int64  `json:"cid" form:"cid"`
	Page  int    `json:"page" form:"page" binding:"required"`
}

type TerminalR struct {
	HostId int64 `json:"host_id" form:"host_id" binding:"required"`
	Cols   int   `json:"cols" form:"cols" binding:"required"`
	Rows   int   `json:"rows" form:"rows" binding:"required"`
}

type CommandLogListR struct {
	Query    string `json:"query" form:"query"`
	HostName string `json:"host_name" form:"host_name"`
	UID      int64  `json:"uid" form:"uid"`
	UserName string `json:"username" form:"username"`
	Limit    int    `json:"limit" form:"limit" binding:"required"`
	Page     int    `json:"page" form:"page" binding:"required"`
}

type ShortcutCommandListR struct {
	Query string `json:"query" form:"query"`
	Limit int    `json:"limit" form:"limit" binding:"required"`
	Page  int    `json:"page" form:"page" binding:"required"`
}

type AddShortcutCommandR struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Cmd         string `json:"cmd" form:"cmd" binding:"required"`
	Description string `json:"description" form:"description"`
}

type DeleteShortcutCommandR struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}

type EditShortcutCommandR struct {
	ID          int64  `json:"id" form:"id" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	Cmd         string `json:"cmd" form:"cmd" binding:"required"`
	Description string `json:"description" form:"description"`
}

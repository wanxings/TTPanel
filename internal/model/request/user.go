package request

type Login struct {
	UserName string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type CreateTemporaryUserR struct {
	ExpireTime int    `json:"expire_time" form:"expire_time" binding:"required"`
	Remark     string `json:"remark" form:"remark"`
}

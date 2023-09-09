package request

type KillProcessR struct {
	Pid   int32 `json:"pid" binding:"required"`
	Force bool  `json:"force"`
}

type DeleteServiceR struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type SetRunLevelR struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Level int    `json:"level" form:"level" binding:"required"`
}

type DeleteLinuxUserR struct {
	Name string `json:"name" form:"name" binding:"required"`
}

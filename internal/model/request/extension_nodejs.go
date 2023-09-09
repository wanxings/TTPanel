package request

type NodejsSetRegistrySourcesR struct {
	Name string `json:"name" binding:"required"`
}

type NodejsSetVersionUrlR struct {
	Url string `json:"url" binding:"required"`
}

type NodejsInstallR struct {
	Version string `json:"version" binding:"required"`
}

type NodejsUninstallR struct {
	Version string `json:"version" binding:"required"`
}

type NodejsSetDefaultEnvR struct {
	Version string `json:"version"`
}

type NodejsNodeModulesListR struct {
	Version string `json:"version" binding:"required"`
}

type NodejsOperationNodeModulesR struct {
	Operation string `json:"operation" binding:"required"`
	Version   string `json:"version" binding:"required"`
	Modules   string `json:"modules" binding:"required"`
}

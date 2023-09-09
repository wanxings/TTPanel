package request

type AddRepositoryR struct {
	Name      string `form:"name" json:"name" binding:"required"`
	Username  string `form:"username" json:"username"`
	Password  string `form:"password" json:"password"`
	Namespace string `form:"namespace" json:"namespace"`
	Remark    string `form:"remark" json:"remark" `
	Url       string `form:"url" json:"url" binding:"required"`
}

type RepositoryListR struct {
	Query string `form:"query" json:"query"`
	Limit int    `form:"limit" json:"limit" binding:"required"`
	Page  int    `form:"page" json:"page" binding:"required"`
}

type EditRepositoryR struct {
	Id        int64  `form:"query" json:"id" binding:"required"`
	Url       string `form:"query" json:"url" binding:"required"`
	Username  string `form:"query" json:"username"`
	Password  string `form:"query" json:"password"`
	Name      string `form:"query" json:"name" binding:"required"`
	Namespace string `form:"query" json:"namespace" binding:"required"`
	Remark    string `form:"query" json:"remark"`
}

type DeleteRepositoryR struct {
	Id int64 `form:"id" json:"id" binding:"required"`
}

type DockerPullImageR struct {
	SourceId int64  `form:"source_id" json:"source_id" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
}

type DockerAppConfigR struct {
	ServerName string `form:"server_name" json:"server_name" binding:"required"`
}

type DockerDeployAppR struct {
	AppName       string         `json:"app_name" binding:"required"`
	ServerName    string         `json:"server_name" binding:"required"`
	AppPath       string         `json:"app_path" binding:"required"`
	Params        map[string]any `json:"params" binding:"required"`
	DockerCompose string         `json:"docker_compose" binding:"required"`
}

type ImageListR struct {
	Query string `form:"query" json:"query"`
	Limit int    `form:"limit" json:"limit" binding:"required"`
	Page  int    `form:"page" json:"page" binding:"required"`
}

type DockerImportImageR struct {
	Path string `form:"path" json:"path" binding:"required"`
}

type DockerBuildImageR struct {
	Path           string            `form:"path" json:"path"`
	DockerFileBody string            `form:"docker_file_body" json:"docker_file_body"`
	Tags           []string          `form:"tags" json:"tags" binding:"required"`
	Labels         map[string]string `form:"labels" json:"labels" binding:"required"`
}

type DockerPushImageR struct {
	RepositoryId int64  `form:"repository_id" json:"repository_id" binding:"required"`
	ImageId      string `form:"image_id" json:"image_id" binding:"required"`
	Tag          string `form:"tag" json:"tag" binding:"required"`
}

type DockerExportImageR struct {
	ImageId string `form:"image_id" json:"image_id" binding:"required"`
	Path    string `form:"path" json:"path" binding:"required"`
	Name    string `form:"name" json:"name" binding:"required"`
}

type DockerBatchDeleteImageR struct {
	ImageIds []string `form:"image_ids" json:"image_ids" binding:"required"`
}

type CreateComposeTemplateR struct {
	Type   string `form:"type" json:"type" binding:"required,oneof=create local"`
	Name   string `form:"name" json:"name" binding:"required"`
	Body   string `form:"body" json:"body"`
	Path   string `form:"path" json:"path"`
	Remark string `form:"remark" json:"remark"`
}

type ComposeTemplateListR struct {
	Query string `form:"query" json:"query"`
	Limit int    `form:"limit" json:"limit" binding:"required"`
	Page  int    `form:"page" json:"page" binding:"required"`
}

type BatchDeleteComposeTemplateR struct {
	Ids []int64 `form:"ids" json:"ids" binding:"required"`
}

type EditComposeTemplateR struct {
	Id     int64  `form:"id" json:"id" binding:"required"`
	Body   string `form:"body" json:"body" binding:"required"`
	Name   string `form:"name" json:"name"  binding:"required"`
	Remark string `form:"remark" json:"remark"`
}

type ComposeTemplatePullImageR struct {
	Id int64 `form:"id" json:"id" binding:"required"`
}

type CreateComposeR struct {
	Name          string `json:"name" binding:"required"`
	Path          string `json:"path" binding:"required"`
	DockerCompose string `json:"docker_compose"`
	Remark        string `json:"remark"`
}

type ComposeListR struct {
	Query string `form:"query" json:"query"`
	Limit int    `form:"limit" json:"limit" binding:"required"`
	Page  int    `form:"page" json:"page" binding:"required"`
}

type ComposeConfigR struct {
	ComposePath string `json:"compose_path" binding:"required"`
}

type OperateComposeR struct {
	ID       int64    `json:"id" binding:"required"`
	Action   string   `json:"action" binding:"required,oneof=up down start stop restart"`
	Services []string `json:"services"`
}

type SaveComposeConfigR struct {
	ComposePath   string         `json:"compose_path"`
	Params        map[string]any `json:"params"`
	DockerCompose string         `json:"docker_compose"`
}

type DeleteComposeR struct {
	Id      int64 `form:"id" json:"id" binding:"required"`
	DelPath bool  `form:"del_path" json:"del_path"`
}

type CreateContainerR struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Image string `json:"image" form:"image" binding:"required"`
	Port  []struct {
		HostPort      int    `json:"host_port" form:"host_port" binding:"required"`
		ContainerPort int    `json:"container_port" form:"container_port" binding:"required"`
		Protocol      string `json:"protocol" form:"protocol" binding:"required,oneof=tcp udp"`
	} `json:"port" form:"port"`
	Cmd        []string `json:"cmd" form:"cmd"`
	AutoRemove bool     `json:"auto_remove" form:"auto_remove"`
	NanoCPUs   int      `json:"nano_cpus" form:"nano_cpus" binding:"required"`
	Memory     int64    `json:"memory" form:"memory" binding:"required"`
	Volumes    []struct {
		HostDir            string `json:"host_dir" form:"host_dir" binding:"required"`
		HostDirPermissions string `json:"host_dir_permissions" form:"host_dir_permissions" binding:"required"`
		ContainerDir       string `json:"container_dir" form:"container_dir" binding:"required"`
	} `json:"volumes" form:"volumes"`
	Labels        map[string]string `json:"labels" form:"labels"`
	Env           []string          `json:"env" form:"env"`
	RestartPolicy string            `json:"restart_policy" form:"restart_policy" binding:"required,oneof=no on-failure always unless-stopped"`
}

type ContainerListR struct {
	Filters  string `form:"filters" json:"filters"`
	GetStats bool   `form:"get_stats" json:"get_stats"`
}

type OperateContainerR struct {
	Name    string `form:"name" json:"name" binding:"required"`
	NewName string `form:"new_name" json:"new_name"`
	Action  string `form:"action" json:"action" binding:"required,oneof=start stop reboot kill pause recover remove rename"`
}

type ContainerLogR struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Since string `form:"since" json:"since"`
	Until string `form:"until" json:"until"`
	Tail  int    `form:"tail" json:"tail"`
}

type ContainerMonitorR struct {
	Id string `form:"id" json:"id" binding:"required"`
}

type ContainerSSHR struct {
	ContainerId string `form:"container_id" json:"container_id" binding:"required"`
	User        string `form:"user" json:"user" binding:"required"`
	Command     string `form:"command" json:"command" binding:"required"`
	Cols        int    `json:"cols" form:"cols" binding:"required"`
	Rows        int    `json:"rows" form:"rows" binding:"required"`
}

type CreateNetworkingR struct {
	Name    string            `form:"name" json:"name" binding:"required"`
	Driver  string            `form:"driver" json:"driver" binding:"required"`
	Options map[string]string `form:"options" json:"options" binding:"required"`
	Subnet  string            `form:"subnet" json:"subnet" binding:"required"`
	Gateway string            `form:"gateway" json:"gateway" binding:"required"`
	IpRange string            `form:"ip_range" json:"ip_range" binding:"required"`
	Labels  map[string]string `form:"labels" json:"labels" binding:"required"`
}

type NetworkingListR struct {
	Query string `form:"name" json:"query"`
	Limit int    `form:"name" json:"limit" binding:"required"`
	Page  int    `form:"name" json:"page" binding:"required"`
}

type DeleteNetworkingR struct {
	Id string `form:"id" json:"id" binding:"required"`
}

type CreateVolumeR struct {
	Name       string            `form:"name"  json:"name" binding:"required"`
	Driver     string            `form:"driver"  json:"driver" binding:"required"`
	DriverOpts map[string]string `form:"driver_opts" json:"driver_opts" binding:"required"`
	Labels     map[string]string `form:"labels" json:"labels" binding:"required"`
}

type VolumeListR struct {
	Name string `form:"name" json:"name"`
}

type DeleteVolumeR struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Force bool   `form:"force" json:"force"`
}

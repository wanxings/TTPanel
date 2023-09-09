package response

type ExtensionsInfoResponse struct {
	Description *Description
	VersionList []*Version `json:"versions" mapstructure:"versions" yaml:"versions"`
}

type DockerBaseStatistics struct {
	Install   bool `json:"install"`
	Status    bool `json:"status"`
	Container struct {
		Total      int `json:"total"`
		Running    int `json:"running"`
		Exited     int `json:"exited"`
		Paused     int `json:"paused"`
		Removing   int `json:"removing"`
		Restarting int `json:"restarting"`
		Created    int `json:"created"`
	} `json:"container"`
	Compose struct {
		Total int64 `json:"total"`
	} `json:"compose"`
	Image struct {
		Total int   `json:"total"`
		Size  int64 `json:"size"`
	} `json:"image"`
	Volume struct {
		Total int `json:"total"`
	} `json:"volume"`
	Network struct {
		Total int `json:"total"`
	} `json:"network"`
}

type Description struct {
	Name         string `json:"name"  yaml:"name" mapstructure:"name"`
	Title        string `json:"title" yaml:"title" mapstructure:"title"`
	TitleEn      string `json:"title_en" yaml:"title_en" mapstructure:"title_en"`
	Website      string `json:"website" yaml:"website" mapstructure:"website"`
	Ps           string `json:"ps" yaml:"ps" mapstructure:"ps"`
	PsEn         string `json:"ps_en" yaml:"ps_en" mapstructure:"ps_en"`
	VersionShell string `json:"version_shell" yaml:"version_shell" mapstructure:"version_shell"`
	ServerPath   string `json:"server_path" yaml:"server_path" mapstructure:"server_path"`
	InitShell    string `json:"init_shell" yaml:"init_shell" mapstructure:"init_shell"`
	Mutex        string `json:"mutex" yaml:"mutex" mapstructure:"mutex"`
	Dependent    string `json:"dependent" yaml:"dependent" mapstructure:"dependent"`
	Version      string `json:"version"`
	Install      bool   `json:"install"`
	Status       bool   `json:"status"`
	ShowIndex    bool   `json:"show_index"`
}
type Version struct {
	Name     string `json:"name" mapstructure:"name" yaml:"name"`
	Version  string `json:"version" mapstructure:"version" yaml:"version"`
	SVersion string `json:"s_version" mapstructure:"s_version" yaml:"s_version"`
	MVersion string `json:"m_version" mapstructure:"m_version" yaml:"m_version"`
	MemLimit int    `json:"mem_limit" mapstructure:"mem_limit" yaml:"mem_limit"`
	CpuLimit int    `json:"cpu_limit" mapstructure:"cpu_limit" yaml:"cpu_limit"`
}

package response

type NodejsConfigResponse struct {
	RegistrySources *NodejsRegistrySources `json:"registry_sources"`
	DownloadUrl     *NodejsDownloadUrls    `json:"download_url"`
	VersionUrl      *NodejsVersionUrl      `json:"version_url"`
	CliVersion      string                 `json:"cli_version"`
}

type NodejsRegistrySources struct {
	List map[string]string `json:"list"`
	Use  string            `json:"use"`
}

type NodejsDownloadUrls struct {
	List map[string]string `json:"list"`
	Use  string            `json:"use"`
}

type NodejsVersionUrl struct {
	List           []string `json:"list"`
	Use            string   `json:"use"`
	LastUpdateTime int64    `json:"last_update_time"`
}

type NodejsVersion struct {
	Version  string      `json:"version"`
	Date     string      `json:"date"`
	Files    []string    `json:"files"`
	Npm      string      `json:"npm"`
	V8       string      `json:"v8"`
	Uv       string      `json:"uv"`
	Zlib     string      `json:"zlib"`
	Openssl  string      `json:"openssl"`
	Modules  string      `json:"modules"`
	Lts      interface{} `json:"lts"`
	Security bool        `json:"security"`
	Install  bool        `json:"install"`
}

type NodejsNodeModulesInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	License     string `json:"license"`
	Homepage    string `json:"homepage"`
}

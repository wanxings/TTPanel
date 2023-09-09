package response

import (
	"github.com/docker/docker/api/types"
	"time"
)

type DockerImageListP struct {
	Id         string             `json:"id"`
	CreateTime int64              `json:"create_time"`
	Tags       []string           `json:"tags"`
	Size       int64              `json:"size"`
	Detail     types.ImageSummary `json:"detail"`
}

type ContainerStatsP struct {
	CpuPercent float64   `json:"cpuPercent"`
	Memory     float64   `json:"memory"`
	Cache      float64   `json:"cache"`
	IoRead     float64   `json:"ioRead"`
	IoWrite    float64   `json:"ioWrite"`
	NetworkRX  float64   `json:"networkRX"`
	NetworkTX  int       `json:"networkTX"`
	ShotTime   time.Time `json:"shotTime"`
}

type ContainerInfo struct {
	types.Container
	Stats *types.StatsJSON `json:"stats"`
}

type DockerAppListP struct {
	CategoryList struct {
		Field1 struct {
			Zh string `json:"zh"`
			En string `json:"en"`
		} `json:"1"`
		Field2 struct {
			Zh string `json:"zh"`
			En string `json:"en"`
		} `json:"2"`
		Field3 struct {
			Zh string `json:"zh"`
			En string `json:"en"`
		} `json:"3"`
	} `json:"category_list"`
	AppList []DockerAppInfo `json:"app_list"`
}

type DockerAppInfo struct {
	Id               int    `json:"id"`
	Category         string `json:"category"`
	ServerName       string `json:"server_name"`
	Sort             int    `json:"sort"`
	TitleZh          string `json:"title_zh"`
	TitleEn          string `json:"title_en"`
	Home             string `json:"home"`
	Github           string `json:"github"`
	Document         string `json:"document"`
	CreateAuthor     string `json:"create_author"`
	CreateAuthorHome string `json:"create_author_home"`
	RevisionStaff    []struct {
		Name string `json:"name"`
		Home string `json:"home"`
	} `json:"revision_staff"`
	CreateTime    int    `json:"create_time"`
	UpdateTime    int    `json:"update_time"`
	Version       string `json:"version"`
	IconBase64    string `json:"icon_base64"`
	DescriptionZh string `json:"description_zh"`
	DescriptionEn string `json:"description_en"`
}

type DockerAppConfigP struct {
	DockerCompose string              `json:"docker_compose"`
	Params        []DockerAppEnvParam `json:"params"`
}

type DockerAppEnvParam struct {
	Id            int    `json:"id"`
	Type          string `json:"type"`
	Key           string `json:"key"`
	TestReg       string `json:"test_reg"`
	Value         any    `json:"value"`
	TitleZh       string `json:"title_zh"`
	TitleEn       string `json:"title_en"`
	DescriptionZh string `json:"description_zh"`
	DescriptionEn string `json:"description_en"`
}

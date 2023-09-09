package request

type SaveDnsConfigR struct {
	DNSList []string `json:"dns_list" form:"dns_list" binding:"required"`
}

type GetTimeZoneConfigR struct {
	MainZone string `json:"main_zone" form:"main_zone"`
}

type SetTimeZoneR struct {
	MainZone string `json:"main_zone" form:"main_zone" binding:"required"`
	SubZone  string `json:"sub_zone" form:"sub_zone" binding:"required"`
}

type AddHostsR struct {
	IP     string `json:"ip" form:"ip" binding:"required"`
	Domain string `json:"domain" form:"domain" binding:"required"`
}

type RemoveHostsR struct {
	Domain string `json:"domain" form:"domain" binding:"required"`
}

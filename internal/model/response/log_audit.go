package response

type LogFileOccupancy struct {
	LogoPath    string `json:"log_path"`
	Description string `json:"description"`
	Count       int    `json:"count"`
	Size        int64  `json:"size"`
}

type SSHLoginLog struct {
	IP            string `json:"ip"`
	Port          string `json:"port"`
	IPAttribution string `json:"ip_attribution"`
	User          string `json:"user"`
	Status        int    `json:"status"`
	LoginTime     int64  `json:"login_time"`
}

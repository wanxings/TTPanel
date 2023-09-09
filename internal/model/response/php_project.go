package response

import (
	"TTPanel/internal/model"
)

type PHPProject struct {
	*model.Project
	PHPVersion string `json:"php_version"`
	//TTWaf      safeResp.TTWafProjectConfig `json:"waf"`
	TTWafStatus    bool        `json:"ttwaf_status"`
	SSL            *SSLDetails `json:"ssl"`
	RunPath        string      `json:"run_path"`
	RunPathDirList []string    `json:"run_path_dir_list"`
	UserIni        bool        `json:"user_ini"`
}

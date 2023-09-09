package util

import (
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
)

// EncodeAuthToBase64 将认证信息编码为 Base64 字符串
func EncodeAuthToBase64(authConfig types.AuthConfig) string {
	authJson, _ := json.Marshal(authConfig)
	return base64.URLEncoding.EncodeToString(authJson)
}

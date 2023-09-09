package util

import (
	"net"
	"regexp"
	"strings"
)

// CheckIP 校验IP或IP+掩码是否合法，true为合法，false为不合法
func CheckIP(ip string) bool {
	ipPattern := `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	maskedIPPattern := `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/\d+$`

	match, _ := regexp.MatchString(ipPattern, ip)
	if match {
		return true
	}
	match, _ = regexp.MatchString(maskedIPPattern, ip)
	return match
}

// IsPublicIP 判断是否为公网IP true为公网IP，false为非公网IP
func IsPublicIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	if parsedIP.IsLoopback() || parsedIP.IsLinkLocalMulticast() || parsedIP.IsLinkLocalUnicast() {
		return false
	}

	if ipv4 := parsedIP.To4(); ipv4 != nil {
		switch true {
		case ipv4[0] == 10:
			return false
		case ipv4[0] == 172 && ipv4[1] >= 16 && ipv4[1] <= 31:
			return false
		case ipv4[0] == 192 && ipv4[1] == 168:
			return false
		default:
			return true
		}
	}

	return false
}

// GetIPType 获取IP类型，返回4为IPv4，返回6为IPv6，返回0为未知类型
func GetIPType(ipString string) int {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return 0
	}

	if strings.Contains(ipString, ":") {
		if ip.To4() == nil {
			return 6 // IPv6
		}
	} else {
		return 4 // IPv4
	}

	return 0
}

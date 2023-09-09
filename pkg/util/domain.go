package util

import (
	"golang.org/x/net/idna"
	"regexp"
	"strings"
)

// CheckDomain 校验域名
func CheckDomain(domain string) bool {
	if StrIsEmpty(domain) {
		return false
	}
	// 使用 idna 转换域名
	domain = strings.ToLower(domain)
	domain, err := idna.ToASCII(domain)
	if err != nil {
		return false
	}

	// 使用正则表达式校验域名是否符合要求
	regex := regexp.MustCompile(`^([\w\-\*]{1,100}\.){1,10}([\w\-]{1,11}|[\w\-]{1,12}\.[\w\-]{1,13})$`)
	if !regex.MatchString(domain) {
		return false
	}
	return true
}

// GetDomain 在域名:端口中获取域名
func GetDomain(siteMenu string) string {
	domain := strings.Split(siteMenu, ":")[0]
	return strings.TrimSpace(domain)
}

// FormatDomain 格式化域名
func FormatDomain(domain string) string {
	domain = strings.ToLower(domain)
	domain, _ = idna.ToASCII(domain)
	return domain
}

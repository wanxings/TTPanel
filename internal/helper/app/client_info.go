package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
)

type ClientInfo struct {
	Ip      string
	Browser int
	System  int
}

func GetClientInfo(c *gin.Context) *ClientInfo {
	return &ClientInfo{
		Ip:      GetIp(c),
		Browser: GetBrowser(c),
		System:  GetSystem(c),
	}
}
func GetIp(c *gin.Context) string {
	return c.ClientIP()
}
func GetSystem(c *gin.Context) int {
	userAgent := c.GetHeader("User-Agent")
	fmt.Println(regexp.MatchString("(?i)win", userAgent))
	match := false
	if match, _ = regexp.MatchString("(?i)win", userAgent); match == true {
		return 1
	} else if match, _ = regexp.MatchString("(?i)mac", userAgent); match == true {
		return 2
	} else if match, _ = regexp.MatchString("(?i)Android", userAgent); match == true {
		return 3
	} else if match, _ = regexp.MatchString("(?i)Ios", userAgent); match == true {
		return 4
	} else if match, _ = regexp.MatchString("(?i)linux", userAgent); match == true {
		return 5
	} else if match, _ = regexp.MatchString("(?i)unix", userAgent); match == true {
		return 6
	} else {
		return 7
	}
}
func GetBrowser(c *gin.Context) int {
	userAgent := c.GetHeader("User-Agent")
	match := false
	if match, _ = regexp.MatchString("(?i)MSIE", userAgent); match == true {
		return 1
	} else if match, _ = regexp.MatchString("(?i)Trident", userAgent); match == true {
		return 1
	} else if match, _ = regexp.MatchString("(?i)Firefox", userAgent); match == true {
		return 2
	} else if match, _ = regexp.MatchString("(?i)Chrome", userAgent); match == true {
		return 3
	} else if match, _ = regexp.MatchString("(?i)Safari", userAgent); match == true {
		return 4
	} else if match, _ = regexp.MatchString("(?i)Opera", userAgent); match == true {
		return 5
	} else {
		return 6
	}
}

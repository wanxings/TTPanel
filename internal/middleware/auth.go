package middleware

import (
	"TTPanel/internal/helper/constant"
	"github.com/gin-gonic/gin"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", "nginx")
		c.Header("Connection", "keep-alive")
		//判断验证方式
		if jwtToken := c.GetHeader("Authorization"); jwtToken != "" { //jwt验证
			c.Set("authMode", constant.AuthModeByJwt)
			c.Next()
		}
		if apiToken := c.GetHeader("ApiToken"); apiToken != "" { //apiToken验证
			c.Set("authMode", constant.AuthModeByApiToken)
			c.Next()
		}
		if temporaryToken := c.Query("temporaryToken"); temporaryToken != "" { //临时token验证
			c.Set("authMode", constant.AuthModeByTemporaryToken)
			c.Next()
		}
		c.Set("authMode", constant.AuthModeByJwt)
		//Todo:后续增加session验证
		c.Next()
	}
}

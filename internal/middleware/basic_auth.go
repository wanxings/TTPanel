package middleware

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if mode, exist := c.Get("authMode"); exist && mode != constant.AuthModeByJwt {
			c.Next()
			return
		}
		//是否开启了basicAuth
		if global.Config.System.BasicAuth.Status == true {
			username, password, ok := c.Request.BasicAuth()
			if !ok || global.Config.System.BasicAuth.Username != username || global.Config.System.BasicAuth.Password != password {
				c.Header("www-authenticate", `Basic realm="Authorization Required"`)
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		c.Next()
	}
}

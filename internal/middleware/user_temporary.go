package middleware

import (
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UserTemporaryAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if mode, exist := c.Get("authMode"); exist && mode != constant.AuthModeByTemporaryToken {
			c.Next()
			return
		}

		session := sessions.Default(c)
		//session.Delete("id")
		id := session.Get("id")
		username := session.Get("username")
		if username == nil {
			response := app.NewResponse(c)
			response.ToErrorResponse(errcode.UnauthorizedTokenTimeout)
			c.Abort()
			return
		}
		c.Set("UID", id)
		c.Set("USERNAME", username)
		c.Next()
	}
}

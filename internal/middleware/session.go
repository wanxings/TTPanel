package middleware

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SessionInit() gin.HandlerFunc {
	store := cookie.NewStore([]byte(global.Config.System.SessionSecret)) //秘钥，需要随机生成
	store.Options(sessions.Options{
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   60 * 10, //初始值为10分钟，后续会根据用户操作刷新
	})
	return sessions.Sessions("TTPanel", store)
}
func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Next()

		if mode, exist := c.Get("authMode"); exist && mode != constant.AuthModeBySession {
			c.Next()
			return
		}

		session := sessions.Default(c)
		//session.Delete("id")
		fmt.Printf("SessionAuth->sessionID:%s", session.ID())
		id := session.Get("id")
		adminToken, ok := global.GoCache.Get("admin_token")
		username := session.Get("username")
		if username == nil || !ok || adminToken != session.Get("admin_token") {
			response := app.NewResponse(c)
			response.ToErrorResponse(errcode.UnauthorizedTokenTimeout)
			session.Clear()
			_ = session.Save()
			c.Abort()
			return
		}
		//刷新过期时间
		session.Options(sessions.Options{MaxAge: global.Config.System.SessionExpire})
		_ = session.Save()
		c.Set("UID", id)
		c.Set("USERNAME", username)
		c.Next()
	}
}

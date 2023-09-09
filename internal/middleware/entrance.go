package middleware

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Entrance 面板入口验证，不论authMode是什么，在这里都要进行验证，如果是session验证，验证通过后，会设置session的entrance_auth为True
func Entrance() gin.HandlerFunc {
	return func(c *gin.Context) {
		mode, ok := c.Get("authMode")
		fmt.Println(mode)
		if ok && mode == constant.AuthModeByJwt { //如果是jwt验证
			session := sessions.Default(c)
			// 是否验证过安全入口
			fmt.Println(session.Get("entrance_auth"))
			global.Log.Debugf("Entrance->session.Get->entrance_auth:%v \n", session.Get("entrance_auth"))
			if session.Get("entrance_auth") == "True" {
				global.Log.Debugf("Entrance verified\n")
				fmt.Println("Entrance verified")
				session.Options(sessions.Options{MaxAge: global.Config.System.SessionExpire})
				_ = session.Save()
				if c.Request.URL.Path == global.Config.System.Entrance {
					c.Redirect(http.StatusFound, "/")
				} else {
					c.Next()
				}
				return
			}
			fmt.Printf("c.Request.URL.Path:%v \n", c.Request.URL.Path)
			fmt.Printf("global.Config.System.Entrance:%v \n", global.Config.System.Entrance)
			if c.Request.URL.Path == global.Config.System.Entrance {
				global.Log.Debugf("Entrance verification passed\n")
				session.Set("entrance_auth", "True")
				_ = session.Save()
				c.Redirect(http.StatusFound, "/")
				return
			}
		}

		if ok && mode == constant.AuthModeBySession { //如果是session验证
			//Todo:入口验证session
		}
		if ok && mode == constant.AuthModeByTemporaryToken { //如果是临时用户验证
			//Todo:入口验证临时用户
			fmt.Println("Entrance->temporaryToken验证")
		}
		if ok && mode == constant.AuthModeByApiToken && global.Config.System.PanelApi.Status { //如果是apiToken验证
			token := c.GetHeader("ApiToken")
			//验证token是否合法
			if token == global.Config.System.PanelApi.Key {
				for _, ip := range global.Config.System.PanelApi.Whitelist {
					if ip == "*" || ip == "all" || ip == c.ClientIP() {
						//验证通过
						c.Set("UID", 0)
						c.Set("USERNAME", "API")
						c.Next()
						return
					}
				}
			}
		}

		// 入口不正确响应状态码，如果是200,响应默认提示
		fmt.Printf("global.Config.System.EntranceErrorCode:%v \n", global.Config.System.EntranceErrorCode)
		if global.Config.System.EntranceErrorCode == 200 {
			c.File(fmt.Sprintf("%s/data/panel_entry_prompt.html", global.Config.System.PanelPath))
		} else {
			c.AbortWithStatus(global.Config.System.EntranceErrorCode)
		}
		c.Abort()
		return
	}
}

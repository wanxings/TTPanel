package middleware

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model"
	"TTPanel/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

// JWT 验证
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if mode, exist := c.Get("authMode"); exist && mode != constant.AuthModeByJwt {
			c.Next()
			return
		}
		var (
			token string
			rCode = errcode.Success
		)
		if tokenGet, ok := c.GetQuery("Authorization"); ok {
			token = tokenGet
			global.Log.Debugf("GetQuery->Authorization:%s\n", token)
		} else {
			token = c.GetHeader("Authorization")
			global.Log.Debugf("GetHeader->Authorization:%s\n", token)
		}
		token = strings.TrimPrefix(token, "Bearer ")
		// 验证前端传过来的token格式，不为空
		if util.StrIsEmpty(token) {
			rCode = errcode.InvalidParams
			global.Log.Debugf("Authorization is None:%s\n", token)
		} else {
			claims, err := app.ParseToken(token)
			if err != nil {
				global.Log.Errorf("JWT->app.ParseToken Error:%s\n", err.Error())
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					rCode = errcode.UnauthorizedTokenError
				default:
					rCode = errcode.UnauthorizedTokenError
				}
			} else {
				c.Set("UID", claims.UID)
				c.Set("USERNAME", claims.USERNAME)
				c.Set("jwt_token", token)

				// 加载用户信息
				user := &model.User{
					ID: claims.UID,
				}
				user, _ = user.Get(global.PanelDB)
				//c.Set("USER", user)

				// 强制下线机制
				//if (global.Config.System.JwtIssuer + ":" + user.Salt) != claims.Issuer {
				//	rCode = errcode.UnauthorizedTokenError
				//}
			}
		}

		//判断是否是API接口访问
		if is, ok := c.Get("IsAPI"); ok && is == "True" {
			global.Log.Debugf("Authorization IsAPI:%s\n", is)
			rCode = errcode.Success
		}

		if rCode != errcode.Success {
			response := app.NewResponse(c)
			response.ToErrorResponse(rCode)
			c.Abort()
			return
		}

		c.Next()
	}
}

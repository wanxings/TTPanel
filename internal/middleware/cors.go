package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	//corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowHeaders("ApiToken")
	return cors.New(corsConfig)
}

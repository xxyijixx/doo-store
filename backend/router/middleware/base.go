package middleware

import (
	"doo-store/backend/core/api/v1/helper"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Base() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("base middleware", c.Request.URL.Path)
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token, X-Xsrf-Token, Language")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 记录登录信息
		c.Set("token", helper.Token(c))
		c.Set("client_version", helper.Version(c))

		c.Next()
	}
}

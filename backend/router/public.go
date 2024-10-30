package router

import (
	v1 "doo-store/backend/core/api/v1"

	"github.com/gin-gonic/gin"
)

type PublicRouter struct {
}

func (a *PublicRouter) InitRouter(Router *gin.RouterGroup) {
	publicRouter := Router.Group("public")
	baseApi := v1.ApiGroupApp.BaseApi
	{
		publicRouter.GET("/health", baseApi.HealthCheck)
	}
}

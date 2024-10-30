package router

import (
	v1 "doo-store/backend/core/api/v1"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
}

func (a *AppRouter) InitRouter(Router *gin.RouterGroup) {
	appRouter := Router.Group("apps")
	baseApi := v1.ApiGroupApp.BaseApi
	{
		appRouter.GET("", baseApi.AppPage)
		appRouter.GET("/sync", baseApi.AppSync)
		appRouter.GET("/:key/detail", baseApi.AppDetailByKey)
		appRouter.GET("/:key/detail/:version", baseApi.AppDeatilByKeyAndVersoin)
		appRouter.POST("/:key/detail/:version", baseApi.AppInstall)
	}
}

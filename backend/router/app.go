package router

import (
	v1 "doo-store/backend/core/api/v1"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
}

func (a *AppRouter) InitRouter(Router *gin.RouterGroup) {
	appRouter := Router.Group("apps")
	baseApi := v1.Api
	{
		appRouter.GET("", baseApi.AppPage)
		appRouter.POST("/:key", baseApi.AppInstall)
		appRouter.PUT("/:key", baseApi.AppInstallOperate)
		appRouter.DELETE("/:key", baseApi.AppUnInstall)
		appRouter.GET("/:key/detail", baseApi.AppDetailByKey)

		appRouter.GET("/installed", baseApi.AppInstalledPage)
		appRouter.GET("/installed/:id/params", baseApi.AppInstalledParams)
		appRouter.PUT("/installed/:id/params", baseApi.AppInstalledUpdateParams)
		appRouter.GET("/installed/:id/logs", baseApi.AppLogs)
		appRouter.GET("/tags", baseApi.AppTags)

		appRouter.GET("/plugin/info", baseApi.GetPluginInfo)

		appRouter.POST("/manage/upload", baseApi.AppUpload)
	}
}

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
		appRouter.GET("", baseApi.ListApps)
		appRouter.POST("/:key", baseApi.InstallApp)
		appRouter.PUT("/:key", baseApi.UpdateAppInstall)
		appRouter.DELETE("/:key", baseApi.UninstallApp)
		appRouter.GET("/:key/detail", baseApi.GetAppDetail)

		appRouter.GET("/installed", baseApi.ListInstalledApps)
		appRouter.GET("/installed/:id/params", baseApi.GetAppParams)
		appRouter.PUT("/installed/:id/params", baseApi.UpdateAppParams)
		appRouter.GET("/installed/:id/logs", baseApi.GetAppLogs)
		appRouter.GET("/tags", baseApi.ListAppTags)

		appRouter.GET("/plugin/info", baseApi.GetInstalledAppInfo)

		appRouter.POST("/manage/upload", baseApi.UploadApp)
	}
}

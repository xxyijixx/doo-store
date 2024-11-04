package router

import (
	"doo-store/backend/i18n"
	entryRouter "doo-store/backend/router"
	"doo-store/backend/router/middleware"
	"doo-store/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

var (
	Router *gin.Engine
)

func Routers() *gin.Engine {
	Router = gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	swaggerRouter := Router.Group("swagger")
	swaggerRouter.GET("/*any", gs.WrapHandler(swaggerFiles.Handler))

	PrivateGroup := Router.Group("/api/v1")
	PrivateGroup.Use(middleware.Base())
	PrivateGroup.Use(i18n.GinI18nLocalize())
	for _, router := range entryRouter.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}

	return Router
}

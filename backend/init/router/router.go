package router

import (
	"doo-store/backend/i18n"
	entryRouter "doo-store/backend/router"
	"doo-store/backend/router/middleware"
	"doo-store/docs"
	"doo-store/web"
	"fmt"
	"net/http"
	"strings"

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

	Router.Static("/store", "./web/dist")

	// t, err := template.New("index").Parse(string(web.IndexByte))
	// if err != nil {
	// 	common.PrintError(fmt.Sprintf("模板解析失败: %s", err.Error()))
	// 	os.Exit(1)
	// }
	// Router.SetHTMLTemplate(t)
	swaggerRouter := Router.Group("swagger")
	swaggerRouter.GET("/*any", gs.WrapHandler(swaggerFiles.Handler))

	PrivateGroup := Router.Group("/api/v1")
	PrivateGroup.Use(middleware.Base())
	PrivateGroup.Use(i18n.GinI18nLocalize())
	for _, router := range entryRouter.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}

	Router.NoRoute(func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		fmt.Println("获取当前的UrlPath", urlPath)
		if strings.HasPrefix(urlPath, "/assets") {
			assets := strings.Replace(urlPath, "/assets", "/assets", -1)
			c.FileFromFS("dist"+assets, http.FS(web.Assets))
			return
		}
		if strings.HasPrefix(urlPath, "/src/assets") {
			assets := strings.Replace(urlPath, "/store/src/assets", "/assets", -1)
			c.FileFromFS("src"+assets, http.FS(web.SrcAssets))
			return
		}
		if strings.HasSuffix(urlPath, "/favicon.ico") {
			c.FileFromFS("/favicon.ico", http.FS(web.Favicon))
			return
		}
		c.JSON(http.StatusNotFound, gin.H{})
	})

	// for _, router := range entryRouter.WebRouterApp {
	// 	router.InitRouter(Router.Group("/store"))
	// }

	return Router
}

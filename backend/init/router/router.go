package router

import (
	"doo-store/backend/i18n"
	entryRouter "doo-store/backend/router"
	"doo-store/backend/router/middleware"
	"doo-store/backend/utils/common"
	"doo-store/docs"
	"doo-store/web"
	"fmt"
	"html/template"
	"net/http"
	"os"
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

	// Router.Use()
	Router.Use(middleware.Base())
	Router.Use(i18n.GinI18nLocalize())

	t, err := template.New("index").Parse(string(web.IndexByte))
	if err != nil {
		common.PrintError(fmt.Sprintf("模板解析失败: %s", err.Error()))
		os.Exit(1)
	}
	Router.SetHTMLTemplate(t)
	swaggerRouter := Router.Group("swagger")
	swaggerRouter.GET("/*any", gs.WrapHandler(swaggerFiles.Handler))

	PrivateGroup := Router.Group("/api/v1")

	for _, router := range entryRouter.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}

	Router.Any("/store/*path", func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		if strings.HasPrefix(urlPath, "/store/assets") {
			assets := strings.Replace(urlPath, "/store/assets", "/assets", -1)
			c.FileFromFS("dist"+assets, http.FS(web.Assets))
			return
		}
		if strings.HasPrefix(urlPath, "/store/src/assets") {
			assets := strings.Replace(urlPath, "/store/src/assets", "/assets", -1)
			c.FileFromFS("src"+assets, http.FS(web.SrcAssets))
			return
		}
		if strings.HasSuffix(urlPath, "/favicon.ico") {
			c.FileFromFS("/favicon.ico", http.FS(web.Favicon))
			return
		}
		c.HTML(http.StatusOK, "index", gin.H{
			"CODE": "",
			"MSG":  "",
		})
	})

	// 页面输出
	return Router
}

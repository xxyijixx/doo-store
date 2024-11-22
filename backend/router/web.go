package router

import (
	"doo-store/web"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebRouter struct {
}

func (a *WebRouter) InitRouter(Router *gin.RouterGroup) {
	// webRouter := Router.Group()
	// 静态资源
	Router.Any("/*any", func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		if strings.HasPrefix(urlPath, "/store/assets") {
			assets := strings.Replace(urlPath, "/store/assets", "/assets", -1)
			// assets = strings.Replace(assets, "/apps/okr/assets", "/assets", -1)
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
	})

}

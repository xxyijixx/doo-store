package v1

import (
	"doo-store/backend/core/api/v1/helper"
	"doo-store/backend/core/dto/request"
	"fmt"

	"github.com/gin-gonic/gin"
)

// @Summary app sync
// @Schemes
// @Description
// @Tags
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /apps/sync [get]
func (*BaseApi) AppSync(c *gin.Context) {

}

// @Summary app page
// @Schemes
// @Description
// @Tags
// @Produce json
// @Param page query integer false "page" default(0)
// @Param pageSize query integer false "pageSize" default(10)
// @Success 200 {string} string "ok"
// @Router /apps [get]
func (*BaseApi) AppPage(c *gin.Context) {
	var req request.AppSearch
	err := helper.CheckBindQueryAndValidate(&req, c)
	if err != nil {
		return
	}
	data, err := appService.AppPage(req)
	if err != nil {
		return
	}
	helper.SuccessWithData(c, data)
}

// @Summary app detail
// @Schemes
// @Description
// @Tags
// @Produce json
// @Param key path string true "key"
// @Success 200 {string} string "ok"
// @Router /apps/{key}/detail [get]
func (*BaseApi) AppDetailByKey(c *gin.Context) {
	key := c.Param("key")
	data, err := appService.AppDetailByKey(key)
	if err != nil {
		return
	}
	helper.SuccessWithData(c, data)
}

// @Summary app detail
// @Schemes
// @Description
// @Tags
// @Produce json
// @Param key path string true "key"
// @Param version path string true "version"
// @Success 200 {string} string "ok"
// @Router /apps/{key}/detail/{version} [get]
func (*BaseApi) AppDeatilByKeyAndVersoin(c *gin.Context) {
	key := c.Param("key")
	version := c.Param("version")
	data, err := appService.AppDetailByKeyAndVersion(key, version)
	if err != nil {
		return
	}
	helper.SuccessWithData(c, data)
}

// @Summary app install
// @Schemes
// @Description
// @Tags
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Param version path string true "version"
// @Param data body request.AppInstall true "RequestBody"
// @Success 200 {string} string "ok"
// @Router /apps/{key}/detail/{version} [post]
func (*BaseApi) AppInstall(c *gin.Context) {
	key := c.Param("key")
	version := c.Param("version")
	var req request.AppInstall
	err := helper.CheckBindAndValidate(&req, c)
	if err != nil {
		return
	}
	req.Key = key
	req.Version = version
	err = appService.AppInstall(req)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	helper.SuccessWithData(c, "安装成功")
}

// @Summary app uninstall
// @Schemes
// @Description
// @Tags
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Param version path string true "version"
// @Param data body request.AppUnInstall true "RequestBody"
// @Success 200 {string} string "ok"
// @Router /apps/{key}/detail/{version} [delete]
func (*BaseApi) AppUnInstall(c *gin.Context) {
	key := c.Param("key")
	version := c.Param("version")
	var req request.AppUnInstall
	err := helper.CheckBindAndValidate(&req, c)
	if err != nil {
		return
	}
	req.Key = key
	req.Version = version
	err = appService.AppUnInstall(req)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	helper.SuccessWithData(c, "安装成功")
}

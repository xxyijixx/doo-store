package v1

import (
	"doo-store/backend/core/api/v1/helper"
	"doo-store/backend/core/dto/request"

	"github.com/gin-gonic/gin"
)

// @Summary app page
// @Schemes
// @Description
// @Security BearerAuth
// @Tags
// @Produce json
// @Param page query integer false "page" default(0)
// @Param page_size query integer false "page_size" default(10)
// @Success 200 {object} dto.Response "success"
// @Router /apps [get]
func (*BaseApi) AppPage(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	var req request.AppSearch
	err = helper.CheckBindQueryAndValidate(&req, c)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	data, err := appService.AppPage(req)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	helper.SuccessWith(c, data)
}

// @Summary app detail
// @Schemes
// @Description
// @Security BearerAuth
// @Tags
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} dto.Response "success"
// @Router /apps/{key}/detail [get]
func (*BaseApi) AppDetailByKey(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	key := c.Param("key")
	data, err := appService.AppDetailByKey(key)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	helper.SuccessWith(c, data)
}

// @Summary app install
// @Schemes
// @Description
// @Security BearerAuth
// @Tags
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Param data body request.AppInstall true "RequestBody"
// @Success 200 {object} dto.Response "success"
// @Router /apps/{key} [post]
func (*BaseApi) AppInstall(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	key := c.Param("key")
	var req request.AppInstall
	err = helper.CheckBindAndValidate(&req, c)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	req.Key = key
	err = appService.AppInstall(req)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	helper.SuccessWith(c, "安装成功")
}

// @Summary app update
// @Schemes
// @Description
// @Security BearerAuth
// @Tags
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Param data body request.AppInstalledOperate true "RequestBody"
// @Success 200 {object} dto.Response "success"
// @Router /apps/{key} [put]
func (*BaseApi) AppInstallOperate(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	key := c.Param("key")
	var req request.AppInstalledOperate
	err = helper.CheckBindAndValidate(&req, c)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	req.Key = key
	err = appService.AppInstallOperate(req)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	helper.SuccessWith(c, "操作成功")
}

// @Summary app uninstall
// @Schemes
// @Description
// @Security BearerAuth
// @Tags
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Param data body request.AppUnInstall true "RequestBody"
// @Success 200 {object} dto.Response "success"
// @Router /apps/{key} [delete]
func (*BaseApi) AppUnInstall(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	key := c.Param("key")
	var req request.AppUnInstall
	err = helper.CheckBindAndValidate(&req, c)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	req.Key = key
	err = appService.AppUnInstall(req)
	if err != nil {
		helper.Error(c, err.Error())
		return
	}
	helper.SuccessWith(c, "卸载成功")
}

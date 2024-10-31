package v1

import (
	"doo-store/backend/core/api/v1/helper"
	"doo-store/backend/core/dto/request"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary app page
// @Schemes
// @Description
// @Security BearerAuth
// @Tags
// @Produce json
// @Param page query integer false "page" default(0)
// @Param pageSize query integer false "pageSize" default(10)
// @Success 200 {string} string "ok"
// @Router /apps [get]
func (*BaseApi) AppPage(c *gin.Context) {
	token, tokenExist := c.Get("token")
	if !tokenExist {
		logrus.Debug("token not exist")
		return
	}
	logrus.Debug("token", token)
	t := token.(string)
	info, _ := dootaskService.GetUserInfo(t)
	fmt.Println("dootask", info)

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
// @Security BearerAuth
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

// @Summary app install
// @Schemes
// @Description
// @Security BearerAuth
// @Tags
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Param data body request.AppInstall true "RequestBody"
// @Success 200 {string} string "ok"
// @Router /apps/{key} [post]
func (*BaseApi) AppInstall(c *gin.Context) {
	key := c.Param("key")
	var req request.AppInstall
	err := helper.CheckBindAndValidate(&req, c)
	if err != nil {
		return
	}
	req.Key = key
	err = appService.AppInstall(req)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	helper.SuccessWithData(c, "安装成功")
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
// @Success 200 {string} string "ok"
// @Router /apps/{key} [put]
func (*BaseApi) AppInstallOperate(c *gin.Context) {
	key := c.Param("key")
	var req request.AppInstalledOperate
	err := helper.CheckBindAndValidate(&req, c)
	if err != nil {
		return
	}
	req.Key = key
	err = appService.AppInstallOperate(req)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	helper.SuccessWithData(c, "安装成功")
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
// @Success 200 {string} string "ok"
// @Router /apps/{key} [delete]
func (*BaseApi) AppUnInstall(c *gin.Context) {
	key := c.Param("key")
	var req request.AppUnInstall
	err := helper.CheckBindAndValidate(&req, c)
	if err != nil {
		return
	}
	req.Key = key
	err = appService.AppUnInstall(req)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	helper.SuccessWithData(c, "安装成功")
}

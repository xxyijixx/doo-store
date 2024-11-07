package v1

import (
	"doo-store/backend/constant"
	"doo-store/backend/core/api/v1/helper"
	"doo-store/backend/core/dto"
	"doo-store/backend/core/dto/request"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 获取插件列表
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param page query integer true "page" default(1)
// @Param page_size query integer true "page_size" default(10)
// @Param class query string false "class"
// @Param name query string false "name"
// @Param id query integer false "id"
// @Param description query string false "description"
// @Success 200 {object} dto.Response "success"
// @Router /apps [get]
func (*BaseApi) AppPage(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	var req request.AppSearch
	err = helper.CheckBindQueryAndValidate(&req, c)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	fmt.Println("req", req)
	data, err := appService.AppPage(dto.ServiceContext{C: c}, req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, data)
}

// @Summary 获取插件详情
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param key path string true "key"
// @Success 200 {object} dto.Response "success"
// @Router /apps/{key}/detail [get]
func (*BaseApi) AppDetailByKey(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	key := c.Param("key")
	data, err := appService.AppDetailByKey(dto.ServiceContext{C: c}, key)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, data)
}

// @Summary 插件安装
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Accept json
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param key path string true "key"
// @Param data body request.AppInstall true "RequestBody"
// @Success 200 {object} dto.Response "success"
// @Router /apps/{key} [post]
func (*BaseApi) AppInstall(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	key := c.Param("key")
	var req request.AppInstall
	err = helper.CheckBindAndValidate(&req, c)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	req.Key = key

	// 校验CPUS和MemoryLimit
	re := regexp.MustCompile(`^(\d+(\.\d+)?(B|b|K|k|M|m|G|g|T|t)?)$|^\d+(\.\d+)?$`)
	if !re.MatchString(req.MemoryLimit) {
		helper.ErrorWith(c, constant.ErrInvalidParameter, nil)
		return
	}
	re = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	if !re.MatchString(req.CPUS) {
		helper.ErrorWith(c, constant.ErrInvalidParameter, nil)
		return
	}

	err = appService.AppInstall(dto.ServiceContext{C: c}, req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, "安装成功")
}

// @Summary app update
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Accept json
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param key path string true "key"
// @Param data body request.AppInstalledOperate true "RequestBody"
// @Success 200 {object} dto.Response "success"
// @Router /apps/{key} [put]
func (*BaseApi) AppInstallOperate(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	key := c.Param("key")
	var req request.AppInstalledOperate
	err = helper.CheckBindAndValidate(&req, c)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	req.Key = key
	err = appService.AppInstallOperate(dto.ServiceContext{C: c}, req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, "操作成功")
}

// @Summary 插件卸载
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Accept json
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param key path string true "key"
// @Param data body request.AppUnInstall true "RequestBody"
// @Success 200 {object} dto.Response "success"
// @Router /apps/{key} [delete]
func (*BaseApi) AppUnInstall(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	key := c.Param("key")
	var req request.AppUnInstall
	err = helper.CheckBindAndValidate(&req, c)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	req.Key = key
	err = appService.AppUnInstall(dto.ServiceContext{C: c}, req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, "卸载成功")
}

// @Summary 获取已安装插件列表
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param page query integer true "page" default(1)
// @Param page_size query integer true "page_size" default(10)
// @Param class query string false "分类"
// @Success 200 {object} dto.Response "success"
// @Router /apps/installed [get]
func (*BaseApi) AppInstalledPage(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	var req request.AppInstalledSearch
	err = helper.CheckBindQueryAndValidate(&req, c)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	data, err := appService.AppInstalledPage(dto.ServiceContext{C: c}, req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, data)
}

// @Summary 获取插件参数信息
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param id path integer true "id"
// @Success 200 {object} dto.Response "success"
// @Router /apps/installed/{id}/params [get]
func (*BaseApi) AppInstalledParams(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := appService.Params(dto.ServiceContext{C: c}, int64(id))
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, data)
}

// @Summary 修改插件参数信息
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param id path integer true "id"
// @Param data body request.AppInstall true "RequestBody"
// @Success 200 {object} dto.Response "success"
// @Router /apps/installed/{id}/params [put]
func (*BaseApi) AppInstalledUpdateParams(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.AppInstall
	err = helper.CheckBindAndValidate(&req, c)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	req.InstalledId = int64(id)

	// 校验CPUS和MemoryLimit
	re := regexp.MustCompile(`^(\d+(\.\d+)?(B|b|K|k|M|m|G|g|T|t)?)$|^\d+(\.\d+)?$`)
	if !re.MatchString(req.MemoryLimit) {
		helper.ErrorWith(c, constant.ErrInvalidParameter, nil)
		return
	}
	re = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	if !re.MatchString(req.CPUS) {
		helper.ErrorWith(c, constant.ErrInvalidParameter, nil)
		return
	}

	err = appService.UpdateParams(dto.ServiceContext{C: c}, req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, nil)
}

// @Summary 获取插件分类信息
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Success 200 {object} dto.Response "success"
// @Router /apps/tags [get]
func (*BaseApi) AppTags(c *gin.Context) {
	err := checkAuth(c, false)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	data, err := appService.AppTags(dto.ServiceContext{C: c})
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, data)
}

func (*BaseApi) AppLogs(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.AppLogsSearch
	err = helper.CheckBindQueryAndValidate(&req, c)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	req.Id = int64(id)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	defer conn.Close()
	_, err = appService.GetLogs(dto.ServiceContext{C: c}, conn, req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
	}
}

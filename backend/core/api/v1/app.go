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
// @Success 200 {object} dto.Response{data=dto.PageResult{items=[]model.App}} "success"
// @Router /apps [get]
func (*BaseApi) ListApps(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	var req request.AppSearch
	if err := helper.ValidateQueryParams(c, &req); err != nil {
		fmt.Printf("请求参数验证失败：%v\n", err)
		helper.ErrorWith(c, err.Error(), nil)
		return
	}

	result, err := appService.ListApps(dto.NewServiceContext(c), req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, result)
}

// @Summary 获取插件详情
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param key path string true "key"
// @Success 200 {object} dto.Response{data=response.AppDetail} "success"
// @Router /apps/{key}/detail [get]
func (*BaseApi) GetAppDetail(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	key := c.Param("key")
	result, err := appService.GetAppDetail(dto.NewServiceContext(c), key)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, result)
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
func (*BaseApi) InstallApp(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	var req request.AppInstall
	if err := helper.ValidateJSONRequest(c, &req); err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	req.Key = c.Param("key")

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

	err = appService.InstallApp(dto.NewServiceContext(c), req)
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
func (*BaseApi) UpdateAppInstall(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	var req request.AppInstalledOperate
	if err := helper.ValidateJSONRequest(c, &req); err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	req.Key = c.Param("key")

	err = appService.UpdateAppInstall(dto.NewServiceContext(c), req)
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
func (*BaseApi) UninstallApp(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	var req request.AppUnInstall
	// err = helper.ValidateJSONRequest(&req, c)
	// if err != nil {
	// 	helper.ErrorWith(c, err.Error(), nil)
	// 	return
	// }
	req.Key = c.Param("key")

	err = appService.UninstallApp(dto.NewServiceContext(c), req)
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
// @Param name query string false "name"
// @Param description query string false "description"
// @Success 200 {object} dto.Response{data=dto.PageResult{items=[]object}} "success"
// @Router /apps/installed [get]
func (*BaseApi) ListInstalledApps(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	var req request.AppInstalledSearch
	if err := helper.ValidateQueryParams(c, &req); err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}

	result, err := appService.ListInstalledApps(dto.NewServiceContext(c), req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, result)
}

// @Summary 获取插件参数信息
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param id path integer true "id"
// @Success 200 {object} dto.Response{data=response.AppInstalledParamsResp} "success"
// @Router /apps/installed/{id}/params [get]
func (*BaseApi) GetAppParams(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	result, err := appService.GetAppParams(dto.NewServiceContext(c), int64(id))
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, result)
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
func (*BaseApi) UpdateAppParams(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.AppInstall
	if err := helper.ValidateJSONRequest(c, &req); err != nil {
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

	result, err := appService.UpdateAppParams(dto.NewServiceContext(c), req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, result)
}

// @Summary 获取插件分类信息
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Success 200 {object} dto.Response{data=[]model.Tag} "success"
// @Router /apps/tags [get]
func (*BaseApi) ListAppTags(c *gin.Context) {
	err := checkAuth(c, false)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	result, err := appService.ListAppTags(dto.NewServiceContext(c))
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, result)
}

// @Summary 获取插件日志信息
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param since query integer false "开始时间(Unix时间戳，秒)"
// @Param until query integer false "结束时间(Unix时间戳，秒)"
// @Param tail query integer true "查询条数" default(1000)
// @Param id path integer true "id"
// @Success 200 {object} dto.Response "success"
// @Router /apps/installed/{id}/logs [get]
func (*BaseApi) GetAppLogs(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.AppLogsSearch
	if err := helper.ValidateQueryParams(c, &req); err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	req.Id = int64(id)
	if req.Tail <= 0 || req.Tail >= 10000 {
		req.Tail = 1000
	}
	result, err := appService.GetAppLogs(dto.NewServiceContext(c), req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, result)
}

// @Summary 上传插件
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param data body request.PluginUpload true "RequestBody"
// @Success 200 {object} dto.Response "success"
// @Router /apps/manage/upload [post]
func (*BaseApi) UploadApp(c *gin.Context) {
	err := checkAuth(c, true)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}

	var req request.PluginUpload
	if err := helper.ValidateJSONRequest(c, &req); err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	fmt.Printf("请求参数：\n%+v\n", req)
	err = appService.UploadApp(dto.NewServiceContext(c), req)
	if err != nil {
		helper.ErrorWith(c, err.Error(), nil)
		return
	}
	helper.SuccessWith(c, nil)
}

// @Summary 获取已安装的插件信息(仅需要登录)
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param key query string true "key"
// @Success 200 {object} object{ret=string,msg=string,data=map[string]any{}}) "success"
// @Router /apps/plugin/info [get]
func (*BaseApi) GetInstalledAppInfo(c *gin.Context) {
	// 身份校验，仅需要登录
	err := checkAuth(c, false)
	if err != nil {
		helper.ErrorWithRet(c, err.Error(), nil)
		return
	}
	var req request.GetInstalledPluginInfo
	if err := helper.ValidateQueryParams(c, &req); err != nil {
		helper.ErrorWithRet(c, err.Error(), nil)
		return
	}
	result, err := appService.GetInstalledAppInfo(dto.NewServiceContext(c), req)
	if err != nil {
		helper.ErrorWithRet(c, err.Error(), nil)
		return
	}
	helper.SuccessWithRet(c, result)
}

// @Summary 获取所有已安装的插件信息(仅需要登录)
// @Schemes
// @Description
// @Security BearerAuth
// @Tags app
// @Produce json
// @Param language header string false "i18n" default(zh)
// @Param key query string true "key"
// @Success 200 {object} object{ret=string,msg=string,data=map[string]any{}}) "success"
// @Router /apps/running [get]
func (*BaseApi) ListRunningApps(c *gin.Context) {
	// TODO 暂时去除身份校验
	// err := checkAuth(c, false)
	// if err != nil {
	// 	helper.ErrorWithRet(c, err.Error(), nil)
	// 	return
	// }
	result, err := appService.ListRunningAppKeys(dto.NewServiceContext(c))
	if err != nil {
		helper.ErrorWithRet(c, err.Error(), nil)
		return
	}
	helper.SuccessWithRet(c, map[string]any{"list": result})
}

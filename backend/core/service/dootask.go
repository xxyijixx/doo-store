package service

import (
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/core/dto"
	"doo-store/backend/utils/common"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// 只需要用户的基础信息

type DootaskService struct {
	client *common.HTTPClient
}

type IDootaskService interface {
	GetUserInfo(token string) (*dto.UserInfoResp, error)
	GetVersoinInfo() (*dto.VersionInfoResp, error)
}

func NewIDootaskService() IDootaskService {
	return &DootaskService{
		client: common.NewHTTPClient(5 * time.Second),
	}
}

// GetUserInfo 获取用户信息
func (d *DootaskService) GetUserInfo(token string) (*dto.UserInfoResp, error) {
	if token == "" {
		return nil, errors.New(constant.ErrNoPermission)
	}
	// url := fmt.Sprintf("%s%s?token=%s", constant.DooTaskUrl, "/api/users/info", token)
	url := fmt.Sprintf("%s%s?token=%s", config.EnvConfig.DOOTASK_URL, "/api/users/info", token)
	result, err := d.client.Get(url)
	if err != nil {
		return nil, err
	}

	info, err := d.UnmarshalAndCheckResponse(result)
	if err != nil {
		return nil, err
	}
	userInfo := new(dto.UserInfoResp)
	if err := common.MapToStruct(info, userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetVersionInfo 获取版本信息
func (d *DootaskService) GetVersoinInfo() (*dto.VersionInfoResp, error) {
	// url := fmt.Sprintf("%s%s", constant.DooTaskUrl, "/api/system/version")
	url := fmt.Sprintf("%s%s", config.EnvConfig.DOOTASK_URL, "/api/system/version")
	result, err := d.client.Get(url)
	if err != nil {
		return nil, err
	}
	versionInfo := &dto.VersionInfoResp{}

	if err := common.StrToStruct(string(result), &versionInfo); err != nil {
		return nil, err
	}
	return versionInfo, nil
}

// 解码并检查返回数据
func (d *DootaskService) UnmarshalAndCheckResponse(resp []byte) (map[string]interface{}, error) {
	var ret map[string]interface{}
	if err := json.Unmarshal(resp, &ret); err != nil {
		// return nil, e.NewErrorWithDetail(constant.ErrDooTaskUnmarshalResponse, err, nil)
		return nil, errors.New(constant.ErrDooTaskUnmarshalResponse)
	}

	retCode, ok := ret["ret"].(float64)
	if !ok {
		return nil, errors.New(constant.ErrDooTaskResponseFormat)
	}

	if retCode != 1 {
		msg, ok := ret["msg"].(string)
		if !ok {
			return nil, errors.New(constant.ErrDooTaskRequestFailed)
		}
		// return nil, e.NewErrorWithDetail(constant.ErrDooTaskRequestFailedWithErr, msg, nil)
		return nil, errors.New(constant.ErrDooTaskRequestFailedWithErr + msg)
	}

	data, ok := ret["data"].(map[string]interface{})
	if !ok {
		dataList, isList := ret["data"].([]interface{})
		if !isList {
			return nil, errors.New(constant.ErrDooTaskDataFormat)
		}

		data = make(map[string]interface{})
		data["list"] = dataList
	}

	return data, nil
}

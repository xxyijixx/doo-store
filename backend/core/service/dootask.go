package service

import (
	"doo-store/backend/constant"
	"doo-store/backend/core/dto"
	"doo-store/backend/utils/common"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
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

func (d *DootaskService) GetUserInfo(token string) (*dto.UserInfoResp, error) {
	if token == "" {
		return nil, errors.New(constant.ErrNoPermission)
	}
	url := fmt.Sprintf("%s%s?token=%s", constant.DooTaskUrl, "/api/users/info", token)
	logrus.Debugf("dootask get user info url: %s", url)
	result, err := d.client.Get(url)
	if err != nil {
		return nil, err
	}

	fmt.Println("DooTask Result", string(result))

	info, err := d.UnmarshalAndCheckResponse(result)
	if err != nil {
		return nil, err
	}
	userInfo := new(dto.UserInfoResp)
	if err := common.MapToStruct(info, userInfo); err != nil {
		return nil, err
	}
	// okr普通人员是否拥有管理员有权限
	// userInfo.OkrAdminOwner = OkrService.GetSettingSuperiorUserId() == userInfo.Userid
	return userInfo, nil
}

func (d *DootaskService) GetVersoinInfo() (*dto.VersionInfoResp, error) {
	url := fmt.Sprintf("%s%s", constant.DooTaskUrl, "/api/system/version")
	logrus.Debugf("dootask get system info url: %s", url)
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
		// return nil, e.WithDetail(constant.ErrDooTaskUnmarshalResponse, err, nil)
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
		// return nil, e.WithDetail(constant.ErrDooTaskRequestFailedWithErr, msg, nil)
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

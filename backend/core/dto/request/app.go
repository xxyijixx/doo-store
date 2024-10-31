package request

import "doo-store/backend/core/dto"

type AppSearch struct {
	dto.PageInfo
	Name string `json:"name"`
}

type AppInstall struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
	Key    string                 `json:"-"`
}

type AppUnInstall struct {
	Key string `json:"-"`
}

type AppInstalledOperate struct {
	Action string                 `json:"action"`
	Key    string                 `json:"-"`
	Params map[string]interface{} `json:"params"`
}

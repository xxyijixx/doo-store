package request

import "doo-store/backend/core/dto"

type AppSearch struct {
	dto.PageInfo
	Name string `json:"name"`
}

type AppInstall struct {
	Name    string                 `json:"name"`
	Params  map[string]interface{} `json:"params"`
	Version string                 `json:"-"`
	Key     string                 `json:"-"`
}

type AppUnInstall struct {
	Name    string `json:"name"`
	Version string `json:"-"`
	Key     string `json:"-"`
}

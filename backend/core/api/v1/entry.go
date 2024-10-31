package v1

import (
	"doo-store/backend/core/service"
)

type BaseApi struct {
}

var Api = new(BaseApi)

type DootaskInfo struct {
	Token    string
	Userinfo interface{}
}

var (
	appService     = service.NewIAppService()
	dootaskService = service.NewIDootaskService()
)

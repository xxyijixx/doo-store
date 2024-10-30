package v1

import "doo-store/backend/core/service"

type BaseApi struct{}

type ApiGroup struct {
	BaseApi
}

var ApiGroupApp = new(ApiGroup)

var (
	appService = service.NewIAppService()
)

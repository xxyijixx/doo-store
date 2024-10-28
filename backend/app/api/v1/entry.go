package v1

type BaseApi struct{}

type ApiGroup struct {
	BaseApi
}

var ApiGroupApp = new(ApiGroup)

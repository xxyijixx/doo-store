package v1

import (
	"doo-store/backend/core/service"
	"net/http"

	"github.com/gorilla/websocket"
)

type BaseApi struct {
}

var Api = new(BaseApi)

type DootaskInfo struct {
	Token    string
	Userinfo interface{}
}

// 设置 WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	appService     = service.NewIAppService()
	dootaskService = service.NewIDootaskService()
)

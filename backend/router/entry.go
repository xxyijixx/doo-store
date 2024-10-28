package router

import "github.com/gin-gonic/gin"

type CommonRouter interface {
	InitRouter(Router *gin.RouterGroup)
}

func commonGroups() []CommonRouter {
	return []CommonRouter{
		&PublicRouter{},
	}
}

func RouterGroups() []CommonRouter {
	return commonGroups()
}

var RouterGroupApp = RouterGroups()

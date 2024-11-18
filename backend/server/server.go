package server

import (
	"doo-store/backend/init/app"
	"doo-store/backend/init/router"
)

func Start() {
	app.Init()
	// validate.Load()
	rootRouter := router.Routers()

	rootRouter.Run(":8080")
}

package server

import (
	"doo-store/backend/init/app"
	"doo-store/backend/init/router"
	"doo-store/backend/init/task"
)

func Start() {
	app.Init()
	task.Init()
	// validate.Load()
	rootRouter := router.Routers()

	rootRouter.Run(":8080")
}

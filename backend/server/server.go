package server

import (
	"doo-store/backend/init/app"
	"doo-store/backend/init/router"
	"doo-store/backend/init/task"
)

func Start() {
	app.Init()
	task.Init()
	rootRouter := router.Init()

	rootRouter.Run(":8080")
}

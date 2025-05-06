package server

import (
	"doo-store/backend/init/app"
	"doo-store/backend/init/redis"
	"doo-store/backend/init/router"
	"doo-store/backend/init/task"
	"os"

	"github.com/gin-gonic/gin"
)

func Start() {
	app.Init()
	task.Init()
	redis.Init()
	rootRouter := router.Init()
	if os.Getenv("ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	rootRouter.Run(":8080")
}

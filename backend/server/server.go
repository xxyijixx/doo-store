package server

import "doo-store/backend/init/router"

func Start() {
	rootRouter := router.Routers()

	rootRouter.Run(":8080")
}

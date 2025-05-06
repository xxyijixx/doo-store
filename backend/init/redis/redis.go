package redis

import (
	"doo-store/backend/utils/redis"
	log "github.com/sirupsen/logrus"
)

func Init() {
	log.Info("初始化Redis连接...")
	if err := redis.Init(); err != nil {
		log.Warnf("Redis初始化失败，某些功能可能无法正常工作: %v", err)
	}
}
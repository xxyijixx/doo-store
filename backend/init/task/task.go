package task

import (
	"context"
	"doo-store/backend/task"
	"time"
)

func Init() {
	task.InitializeGlobalManager(100, 3)

	// 初始化Docker容器监控
	monitor, err := task.NewDockerMonitor(context.Background())
	if err != nil {
		panic(err)
	}
	monitor.StartMonitoring(30 * time.Second)
}

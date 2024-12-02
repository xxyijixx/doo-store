package task

import "doo-store/backend/task"

func Init() {
	task.InitializeGlobalManager(100, 3)
	
	// 初始化Docker容器监控
	if err := task.InitDockerMonitoring(); err != nil {
		panic(err)
	}
}

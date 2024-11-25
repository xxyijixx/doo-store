package task

import "doo-store/backend/task"

func Init() {
	task.InitializeGlobalManager(100, 3)
}

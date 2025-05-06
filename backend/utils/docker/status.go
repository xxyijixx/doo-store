package docker

const (
	ContainerStatusCreated    = "created"
	ContainerStatusRunning    = "running"
	ContainerStatusPaused     = "paused"
	ContainerStatusRestarting = "restarting"
	ContainerStatusRemoving   = "removing"
	ContainerStatusExited     = "exited"
	ContainerStatusDead       = "dead"

	// 自定义状态
	CustomContainerStatusInit = "init" // 容器未初始化
)

package constant

const (
	Running    = "Running"
	UnHealthy  = "UnHealthy"
	Error      = "Error"
	Stopped    = "Stopped"
	Installing = "Installing"
	Paused     = "Paused"
	UpErr      = "UpErr"

	AppNormal   = "Normal"
	AppTakeDown = "TakeDown"

	CPUS          = "CPUS"
	MemoryLimit   = "MEMORY_LIMIT"
	HostIP        = "HOST_IP"
	ContainerName = "CONTAINER_NAME"
)

type AppOperate string

var (
	Start   AppOperate = "start"
	Stop    AppOperate = "stop"
	Restart AppOperate = "restart"
	Delete  AppOperate = "delete"
	Backup  AppOperate = "backup"
	Update  AppOperate = "update"
	Upgrade AppOperate = "upgrade"
)

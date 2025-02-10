package constant

const (
	Running    = "Running"
	UnHealthy  = "UnHealthy"
	Restarting = "Restarting"
	Error      = "Error"
	Dead       = "Dead"
	Stopped    = "Stopped"
	Installing = "Installing"
	Paused     = "Paused"
	UpErr      = "UpErr"
	Unknown    = "Unknown"

	AppNormal   = "Normal"
	AppUnused   = "Unused"
	AppTakeDown = "TakeDown"
	AppInUse    = "InUse"

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

package model

type App struct {
	BaseModel
	Name           string `json:"name" gorm:"size:60;not null;default:''"`
	Key            string `json:"key" gorm:"size:60;not null;default:'';unique"`
	Icon           string `json:"icon"`
	Description    string `json:"description" gorm:"size:255;not null;default:''"`
	Github         string `json:"github"`
	Class          string `json:"class" gorm:"size:60;not null;default:''"`
	DependsVersion string `json:"depends_version"`
	Sort           int    `json:"sort" gorm:"default:999"`
	Status         string `json:"status" gorm:"size:20;not null;default:''"`
}

func (*App) TableName() string {
	return TableName("apps")
}


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

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
	// 插件状态
	PluginStatusRunning    = "Running"
	PluginStatusUnHealthy  = "UnHealthy"
	PluginStatusRestarting = "Restarting"
	PluginStatusError      = "Error"
	PluginStatusDead       = "Dead"
	PluginStatusStopped    = "Stopped"
	PluginStatusInstalling = "Installing"
	PluginStatusPaused     = "Paused"
	PluginStatusUpErr      = "UpErr"
	PluginStatusUnknown    = "Unknown"

	// 插件安装状态
	AppNormal   = "Normal"
	AppUnused   = "Unused"
	AppTakeDown = "TakeDown"
	AppInUse    = "InUse"

	// 环境变量
	CPUS          = "CPUS"
	MemoryLimit   = "MEMORY_LIMIT"
	HostIP        = "HOST_IP"
	ContainerName = "CONTAINER_NAME"
)

// 插件操作
type PluginAction string

var (
	PluginActionStart   PluginAction = "start"
	PluginActionStop    PluginAction = "stop"
	PluginActionRestart PluginAction = "restart"
	PluginActionDelete  PluginAction = "delete"
	PluginActionBackup  PluginAction = "backup"
	PluginActionUpdate  PluginAction = "update"
	PluginActionUpgrade PluginAction = "upgrade"
)

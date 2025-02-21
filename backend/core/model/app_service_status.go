package model

type AppServiceStatus struct {
	BaseModel
	ContainerName string `json:"container_name" gorm:"size:60;comment:容器名;not null;default:''"`
	ServiceName   string `json:"service_name" gorm:"size:60;comment:服务名;not null;default:''"`
	IpAddress     string `json:"ip_address" gorm:"size:60;comment:IP地址;not null;default:''"`
	Image         string `json:"image_name" gorm:"size:60;comment:镜像;not null;default:''"`
	Message       string `json:"message" gorm:"comment:消息;default:''"`
	Status        string `json:"status" gorm:"size:20;comment:状态;not null;default:''"`
	InstallID     int64  `json:"install_id" gorm:"comment:安装ID;not null"`
}

func (*AppServiceStatus) TableName() string {
	return TableName("app_service_status")
}

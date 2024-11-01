package model

type AppInstalled struct {
	BaseModel
	Name          string `json:"name" gorm:"size:60;not null;default:''"`
	AppID         int64  `json:"app_id"`
	AppDetailID   int64  `json:"app_detail_id"`
	Key           string `json:"key" gorm:"size:60"`
	Version       string `json:"version" gorm:"size:40;not null;default:''"`
	Params        string `json:"params" gorm:"type:text"`
	Env           string `json:"env" gorm:"type:text"`
	DockerCompose string `json:"docker_compose" gorm:"type:text"`
	Status        string `json:"status" gorm:"size:20;not null;default:''"`
}

func (*AppInstalled) TableName() string {
	return TableName("app_installed")
}

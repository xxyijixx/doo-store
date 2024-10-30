package model

type AppInstalled struct {
	BaseModel
	Name          string `json:"name" gorm:"size:60;not null;default:''"`
	AppID         int64  `json:"appId"`
	AppDetailID   int64  `json:"appDetailId"`
	Key           string `json:"key" gorm:"size:60"`
	Version       string `json:"version" gorm:"size:40;not null;default:''"`
	Params        string `json:"params" gorm:"type:text"`
	Env           string `json:"env" gorm:"type:text"`
	DockerCompose string `json:"dockerCompose" gorm:"type:text"`
	Status        string `json:"status" gorm:"size:20;not null;default:''"`
}

func (*AppInstalled) TableName() string {
	return TableName("app_installed")
}

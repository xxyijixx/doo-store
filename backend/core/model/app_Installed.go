package model

type AppInstalled struct {
	BaseModel
	Name          string `json:"name" gorm:"size:60;not null;default:''"`
	AppID         int64  `json:"appId"`
	AppDetailID   int64  `json:"appDetailId"`
	Version       string `json:"version" gorm:"not null;default:''"`
	Params        string `json:"params"`
	Env           string `json:"env"`
	DockerCompose string `json:"dockerCompose"`
	Status        string `json:"status" gorm:"not null;default:''"`
}

func (*AppInstalled) TableName() string {
	return TableName("app_installed")
}

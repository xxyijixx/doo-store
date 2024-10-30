package model

type AppDetail struct {
	BaseModel
	AppID         int64  `json:"appId"`
	Version       string `json:"version" gorm:"size:60;not null;default:''"`
	Params        string `json:"-"`
	DockerCompose string `json:"dockerCompose"`
	Status        string `json:"status" gorm:"size:60;not null;default:''"`
}

func (*AppDetail) TableName() string {
	return TableName("app_details")
}

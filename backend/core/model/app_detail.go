package model

type AppDetail struct {
	BaseModel
	AppID          int64  `json:"appId"`
	Repo           string `json:"repo"`
	Version        string `json:"version" gorm:"size:40;not null;default:''"`
	DependsVersion string `json:"depends_version"`
	Params         string `json:"-" gorm:"type:text"`
	DockerCompose  string `json:"dockerCompose" gorm:"type:text"`
	Status         string `json:"status" gorm:"size:200;not null;default:''"`
}

func (*AppDetail) TableName() string {
	return TableName("app_details")
}

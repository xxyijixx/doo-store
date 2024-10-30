package model

type Tag struct {
	BaseModel
	Key  string `json:"key" gorm:"size:50;not null;default:''"`
	Name string `json:"name" gorm:"size:50;not null;default:''"`
	Sort int    `json:"sort" gorm:"default:99"`
}

func (*Tag) TableName() string {
	return TableName("tags")
}

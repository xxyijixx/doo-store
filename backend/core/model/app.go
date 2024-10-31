package model

type App struct {
	BaseModel
	Name        string `json:"name" gorm:"size:60;not null;default:''"`
	Key         string `json:"key" gorm:"size:60;not null;default:'';unique"`
	Icon        string `json:"icon"`
	Description string `json:"description" gorm:"size:255;not null;default:''"`
	Github      string `json:"github"`
	Type        string `json:"type" gorm:"size:60;not null;default:''"`
	Sort        int    `json:"sort" gorm:"default:999"`
	Status      string `json:"status" gorm:"size:20;not null;default:''"`
}

func (*App) TableName() string {
	return TableName("apps")
}

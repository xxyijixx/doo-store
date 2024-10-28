package model

type App struct {
	BaseModel
	Name string `json:"name"`
}

func (*App) TableName() string {
	return TableName("apps")
}

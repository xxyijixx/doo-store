package model

type AppTag struct {
	BaseModel
	AppID int64 `json:"appId"`
	TagID int64 `json:"tagId"`
}

func (*AppTag) TableName() string {
	return TableName("app_tags")
}

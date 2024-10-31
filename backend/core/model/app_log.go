package model

import "time"

type AppLog struct {
	ID             int64     `gorm:"primaryKey;not null;autoIncrement:true;comment:'id'" json:"id"`
	AppInstalledId int64     `gorm:"not null;comment:'app installed id'" json:"appInstalledId"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (*AppLog) TableName() string {
	return TableName("app_logs")
}

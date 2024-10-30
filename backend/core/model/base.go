package model

import (
	"doo-store/backend/config"
	"time"
)

type BaseModel struct {
	ID        int64     `gorm:"primaryKey;not null;autoIncrement:true;comment:'id'" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func TableName(name string) string {
	return config.EnvConfig.DB_PREFIX + name
}

package migrate

import (
	"doo-store/backend/config"
	"doo-store/backend/core/model"
	"fmt"

	"gorm.io/gorm"
)

func Migrate() {
	var err error
	db, err := gorm.Open(
		config.EnvConfig.GetGormDialector(),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	if err != nil {
		panic(fmt.Errorf("db connection failed: %v", err))
	}
	err = db.AutoMigrate(&model.App{}, &model.AppDetail{}, &model.AppInstalled{}, &model.AppTag{}, &model.Tag{}, &model.AppLog{})
	if err != nil {
		panic(fmt.Errorf("db migrate failed: %v", err))
	}
}

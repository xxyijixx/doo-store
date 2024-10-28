package migrate

import (
	"doo-store/backend/app/model"
	"doo-store/backend/config"
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
	err = db.AutoMigrate(&model.App{})
	if err != nil {
		panic(fmt.Errorf("db migrate failed: %v", err))
	}
}

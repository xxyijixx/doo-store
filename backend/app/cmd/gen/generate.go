package gen

import (
	"doo-store/backend/app/model"
	"doo-store/backend/config"

	"gorm.io/gen"
	"gorm.io/gorm"
)

func Generate() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./backend/app/repo",
		Mode:    gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface, // generate mode
	})

	gormdb, _ := gorm.Open(config.EnvConfig.GetGormDialector(), &gorm.Config{

		// 	NamingStrategy: schema.NamingStrategy{
		// 		TablePrefix: config.EnvConfig.DB_PREFIX,
		// 	},
	})

	// reuse your gorm db
	g.UseDB(gormdb)

	g.ApplyBasic(model.App{})

	// Generate the code
	g.Execute()
}

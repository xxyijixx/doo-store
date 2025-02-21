package gen

import (
	"doo-store/backend/core/model"

	"gorm.io/gen"
)

func Generate() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./backend/core/repo",
		Mode:    gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface, // generate mode
	})

	// gormdb, _ := gorm.Open(config.EnvConfig.GetGormDialector(), &gorm.Config{

	// 	// 	NamingStrategy: schema.NamingStrategy{
	// 	// 		TablePrefix: config.EnvConfig.DB_PREFIX,
	// 	// 	},
	// })

	// // reuse your gorm db
	// g.UseDB(gormdb)

	g.ApplyBasic(model.App{}, model.AppDetail{}, model.AppInstalled{}, model.AppServiceStatus{}, model.AppTag{}, model.Tag{}, model.AppLog{})

	// Generate the code
	g.Execute()
}

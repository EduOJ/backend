package database

import (
	"fmt"

	"github.com/EduOJ/backend/base"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupDatabaseForTest setups a in-memory db for testing purpose.
// Shouldn't be called out of test.
func SetupDatabaseForTest() func() {
	oldDB := base.DB
	x, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	viper.Set("database.dialect", "sqlite")
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	base.DB = x
	sqlDB, err := base.DB.DB()
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	sqlDB.SetMaxOpenConns(1)
	Migrate()
	return func() {
		base.DB = oldDB
	}
}

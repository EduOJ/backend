package database

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/base"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var testDatabaseLock = sync.Mutex{}

// SetupDatabaseForTest setups a in-memory db for testing purpose.
// Shouldn't be called out of test.
func SetupDatabaseForTest() func() {
	testDatabaseLock.Lock()
	oldDB := base.DB
	x, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
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
		testDatabaseLock.Unlock()
	}
}

package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base"
	"sync"
)

var testDatabaseLock = sync.Mutex{}

// SetupDatabaseForTest setups a in-memory db for testing purpose.
// Shouldn't be called out of test.
func SetupDatabaseForTest() func() {
	testDatabaseLock.Lock()
	oldDB := base.DB
	x, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	base.DB = x
	base.DB.DB().SetMaxOpenConns(1)
	base.DB.LogMode(false)
	Migrate()
	return func() {
		base.DB = oldDB
		testDatabaseLock.Unlock()
	}
}

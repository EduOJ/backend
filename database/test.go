package database

import (
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base"
)

// SetupDatabaseForTest setups a in-memory db for testing purpose.
// Shouldn't be called out of test.
func SetupDatabaseForTest() func() {
	oldDB := base.DB
	base.DB, _ = gorm.Open("sqlite3", ":memory:")
	Migrate()
	return func() {
		base.DB = oldDB
	}
}

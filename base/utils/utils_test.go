package utils

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	defer database.SetupDatabaseForTest()()

	PanicIfDBError(base.DB.AutoMigrate(&TestObject{}), "could not create table for test object")
	os.Exit(m.Run())
}

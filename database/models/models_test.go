package models

import (
	"github.com/leoleoasd/EduOJBackend/base/exit"
	"github.com/leoleoasd/EduOJBackend/database"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	defer database.SetupDatabaseForTest()()
	defer exit.SetupExitForTest()()
	os.Exit(m.Run())
}

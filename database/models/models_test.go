package models

import (
	"os"
	"testing"

	"github.com/EduOJ/backend/base/exit"
	"github.com/EduOJ/backend/database"
)

func TestMain(m *testing.M) {
	defer database.SetupDatabaseForTest()()
	defer exit.SetupExitForTest()()
	os.Exit(m.Run())
}

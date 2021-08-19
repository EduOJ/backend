package notification_test

import (
	"github.com/EduOJ/backend/base/notification"
	"github.com/EduOJ/backend/database"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	defer database.SetupDatabaseForTest()()
	notification.RegisteredPreferredNoticedMethod = []string{
		"test_send_message_registered_method",
		"register_existing_method",
	}
	os.Exit(m.Run())
}

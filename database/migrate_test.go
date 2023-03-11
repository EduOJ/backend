package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrate(t *testing.T) {
	defer SetupDatabaseForTest()()
	m := GetMigration()
	assert.NoError(t, m.RollbackTo("start"))
}

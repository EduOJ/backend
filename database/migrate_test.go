package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMigrate(t *testing.T) {
	defer SetupDatabaseForTest()()
	m := GetMigration()
	assert.NoError(t, m.RollbackTo("start"))
}

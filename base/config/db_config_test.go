package config

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDBConfig(t *testing.T) {
	t.Cleanup(database.SetupDatabaseForTest())
	assert.Equal(t, "", GetDBConfig("non_exiting"))
	count := 0
	base.DB.Model(&models.Config{}).Count(&count)
	config := models.Config{}
	base.DB.First(&config)
	assert.Equal(t, 1, count)
	assert.Equal(t, "non_exiting", config.Key)
	assert.Equal(t, "", *config.Value)
	assert.Equal(t, nil, SetDBConfig("non_exiting", "2333"))
	count = 0
	base.DB.Model(&models.Config{}).Count(&count)
	config = models.Config{}
	base.DB.First(&config)
	assert.Equal(t, 1, count)
	assert.Equal(t, "non_exiting", config.Key)
	assert.Equal(t, "2333", *config.Value)
}

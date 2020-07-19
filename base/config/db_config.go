package config

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
)

func GetDBConfig(key string) string {
	c := models.Config{
		Key: key,
	}
	base.DB.Where(&c).FirstOrCreate(&c)
	return *c.Value
}

func SetDBConfig(key string, value string) error {
	c := models.Config{}
	base.DB.FirstOrCreate(&c, models.Config{
		Key: key,
	})
	c.Value = &value
	if err := base.DB.Save(&c).Error; err != nil {
		return errors.Wrap(err, "could not save config to db")
	}
	return nil
}

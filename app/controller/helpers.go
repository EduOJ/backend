package controller

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func panicIfDBError(db *gorm.DB, message string) {
	if db.Error != nil {
		panic(errors.Wrap(db.Error, message))
	}
}

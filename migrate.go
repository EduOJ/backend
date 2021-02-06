package main

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"time"
)

func doMigrate() {

	if len(args) == 1 {
		readConfig()
		initGorm()
		log.Fatal("Migrate succeed!")
	} else {
		switch args[1] {
		case "help":
			fmt.Println(`Usage:
./backend migrate <command>
command: The command to run. default: migrate
Commands:
migrate: run migrations.
rollback: rollback all migrations.
rollback_last: rollback last migration.`)
		case "migrate":
			readConfig()
			initGorm()
			log.Fatal("Migrate succeed!")
		case "rollback":
			readConfig()
			initGorm(false)
			m := database.GetMigration()
			err := m.RollbackTo("start")
			if err != nil {
				log.Error(err)
			} else {
				log.Fatal("Migrate succeed!")
			}
		case "rollback_last":
			readConfig()
			initGorm(false)
			m := database.GetMigration()
			err := m.RollbackLast()
			if err != nil {
				log.Error(err)
			} else {
				log.Fatal("Migrate succeed!")
			}
		}
	}
	p := models.Problem{
		ID:                 0,
		Name:               "",
		Description:        "",
		AttachmentFileName: "",
		Public:             false,
		Privacy:            false,
		MemoryLimit:        0,
		TimeLimit:          0,
		LanguageAllowed:    nil,
		CompileEnvironment: "",
		CompareScriptID:    0,
		TestCases:          nil,
		CreatedAt:          time.Time{},
		UpdatedAt:          time.Time{},
		DeletedAt:          nil,
	}
	log.Debug(base.DB.Save(&p).Error)
	p1 := models.Problem{}
	log.Debug(base.DB.First(&p1, 1).Error)
	log.Debug(p1)

}

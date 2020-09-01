package main

import (
	"bufio"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
)

type testClass struct {
	ID   uint
	name string
}

func (t *testClass) GetID() uint {
	return t.ID
}
func (t *testClass) TypeName() string {
	return t.name
}

func editPermission() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			log.Fatal("Editing permission failed.")
		}
	}()

	readConfig()
	initGorm()
	initLog()

	if len(args) == 1 {
		quit := false
		log.Debug("Entered interactive mode, enter \"help\" to get usage help")
		for !quit {
			fmt.Print("\033[1mEdit Permission> \033[0m")
			input, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(errors.Wrap(err, "Error reading editing permission command"))
			}
			args = strings.Split(input[:len(input)-1], " ")
			quit = doEditPermission(args)
		}
	} else {
		doEditPermission(args[1:])
	}

}

func doEditPermission(args []string) (end bool) {
	var err error
	var operation string
	switch args[0] {
	case "help", "h":
		log.Info(`
Edit Permission

Usage:
  Single execution: $ EduOJ (edit-permission|edit-perm|ep) (operation) <args>...

  Enter interactive mode: $ EduOJ (edit-permission|edit-perm|ep)
  Command format in interactive mode:  (operation) <args>...

operations:
  (help|h)
  (create-role|cr) <name> [<target>]
  (grant-role|gr) <user_id|username> <role_id|role_name> [<target_id>]
  (add-permission|ap) <role_id|role_name> <permission>
  (quit|q)

Note:
  When the search value matches the name and ID at the same time, the system
  always selects the object that matches the ID.`)
	case "create-role", "cr":
		// edit_permission (create-role|cr) <name> [<target>]
		operation = "Creating role"
		err = validateArgumentsCount(len(args), 2, 3)
		if err != nil {
			break
		}
		r := models.Role{
			Name: args[1],
		}
		if len(args) == 3 {
			r.Target = &args[2]
		}
		err = base.DB.Create(&r).Error
	case "grant-role", "gr":
		// edit_permission (grant-role|gr) <user_id|username> <role_id|role_name> [<target_id>]
		operation = "Granting role"
		err = validateArgumentsCount(len(args), 3, 4)
		if err != nil {
			break
		}
		var user *models.User
		user, err = utils.FindUser(args[1])
		if err != nil {
			err = errors.Wrap(err, "find user")
			break
		}
		var role *models.Role
		role, err = findRole(args[2])
		if err != nil {
			break
		}
		if len(args) == 3 {
			user.GrantRole(*role)
		} else {
			var targetId uint64
			targetId, err = strconv.ParseUint(args[3], 10, 32)
			if err != nil {
				break
			}
			target := testClass{
				ID:   uint(targetId),
				name: *role.Target,
			}
			user.GrantRole(*role, &target)

		}
	case "add-permission", "ap":
		// edit_permission (add-permission|ap) <role_id|role_name> <permission>
		operation = "Adding permission"
		err = validateArgumentsCount(len(args), 3, 3)
		if err != nil {
			break
		}
		var role *models.Role
		role, err = findRole(args[1])
		if err != nil {
			break
		}
		role.AddPermission(args[2])
	case "quit", "q":
		log.Debug("Exited editing permission mode.")
		return true
	default:
		log.Debug("Unknown operation \"" + args[0] + "\".")
	}
	if operation != "" {
		if err == nil {
			log.Fatal(operation + " succeed!")
		} else {
			log.Error(err)
			log.Fatal(operation + " failed.")
		}
	}
	return false
}

func findRole(id string) (*models.Role, error) {
	role := models.Role{}
	err := base.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		err = base.DB.Where("name = ?", id).First(&role).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.New("role record not found")
			} else {
				return nil, errors.Wrap(err, "could not query role")
			}
		}
	}
	return &role, nil
}

func validateArgumentsCount(count int, min int, max int) (err error) {
	if count < min {
		err = errors.New("Too few command line parameters")
	} else if count > max {
		err = errors.New("Too many command line parameters")
	}
	return
}

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
		log.Debug("Entered editing permission mode, enter \"help\" to get usage help")
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
  The EP(Edit Permission) command can be executed as a single command in EP mode,
  or as program parameters of EduOJ.

  EP command:     (edit_permission|edit_perm|ep) (operation) <args>...

  Enter EP mode: 
    $ go run (path) (edit_permission|edit_perm|ep)
  Edit permission with out entering EP mode: 
    $ go run (path) (edit_permission|edit_perm|ep) (operation) <args>...

operations:
  edit_permission (help|h)
  edit_permission (createRole|cr) <name> [<target>]
  edit_permission (grantRole|gr) <user_id|username> <role_id|role_name> [<target_id>]
  edit_permission (addPermission|ap) <role_id|role_name> <permission>
  edit_permission (quit|q)

Note:
  When the search value matches the name and ID at the same time, the system
  always selects the object that matches the ID.`)
	case "createRole", "cr":
		// edit_permission (createRole|cr) <name> [<target>]
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
	case "grantRole", "gr":
		// edit_permission (grantRole|gr) <user_id|username> <role_id|role_name> [<target_id>]
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
			target := struct{}{}
			var targetId uint64
			targetId, err = strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				break
			}
			err = base.DB.Table(*role.Target).First(&target, targetId).Error
			if err != nil {
				break
			}
			log.Debug(target)
			//user.GrantRole(*role,target)
		}
	case "addPermission", "ap":
		// edit_permission (addPermission|ap) <role_id|role_name> <permission>
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

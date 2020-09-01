package main

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"github.com/xlab/treeprint"
	"os"
	"strconv"
	"strings"
)

// For role granting.
// Implements database/models/HasRole interface.
type dummyHasRole struct {
	ID   uint
	name string
}

func (t *dummyHasRole) GetID() uint {
	return t.ID
}
func (t *dummyHasRole) TypeName() string {
	return t.name
}

func permission() {
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
		log.Debug(`Entering interactive mode, enter "help" for help.`)
		for !quit {
			_, err := color.New(color.Bold).Print("EduOJ Permission> ")
			if err != nil {
				log.Error(errors.Wrap(err, "fail to print"))
				log.Fatal("Editing permission failed.")
				return
			}
			input, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(errors.Wrap(err, "Error reading command"))
				continue
			}
			args = strings.Split(input[:len(input)-1], " ")
			quit = doPermission(args)
		}
	} else {
		doPermission(args[1:])
	}
}

func doPermission(args []string) (end bool) {
	var err error
	var operation string
	switch args[0] {
	case "help", "h":
		log.Info(`
Edit Permission

Usage:
  One-line execution: $ EduOJ (permission|perm) (command) <args>...

  Enter interactive mode: $ EduOJ (permission|perm)
  Command format in interactive mode:  (command) <args>...

commands:
  (help|h)
  (list-roles|lr) [<role_id|role_name>]
  (create-role|cr) <name> [<target>]
  (grant-role|gr) <user_id|username> <role_id|role_name> [<target_id>]
  (delete-role|dr) <role_id|role_name>
  (add-permission|ap) <role_id|role_name> <permission>
  (quit|q)

Note:
  When the search value matches the name and ID at the same time, the system
  always selects the object that matches the ID.`)
	case "create-role", "cr":
		// (create-role|cr) <name> [<target>]
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
	case "list-roles", "lr":
		// (list-roles|lr) [<role_id|role_name>]
		operation = "Listing roles"
		err = validateArgumentsCount(len(args), 1, 2)
		tree := treeprint.New()
		if len(args) == 1 {
			var roles []models.Role
			err = base.DB.Set("gorm:auto_preload", true).Find(&roles).Error
			if err != nil {
				break
			}
			for _, role := range roles {
				listRole(tree, &role)
			}
		} else {
			var role *models.Role
			role, err = findRole(args[1])
			if err != nil {
				break
			}
			listRole(tree, role)
		}

		log.Info("\n" + tree.String())
	case "grant-role", "gr":
		// (grant-role|gr) <user_id|username> <role_id|role_name> [<target_id>]
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
			target := dummyHasRole{
				ID:   uint(targetId),
				name: *role.Target,
			}
			user.GrantRole(*role, &target)
		}
	case "delete-role", "dr":
		// (delete-role|dr) <role_id|role_name>
		operation = "Deleting role"
		err = validateArgumentsCount(len(args), 2, 2)
		var role *models.Role
		role, err := findRole(args[1])
		if err != nil {
			break
		}
		err = base.DB.Delete(&models.Permission{}, "role_id = ?", role.ID).Error
		if err != nil {
			break
		}
		err = base.DB.Delete(&role).Error
		if err != nil {
			break
		}
	case "add-permission", "ap":
		// (add-permission|ap) <role_id|role_name> <permission>
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
	err := base.DB.Set("gorm:auto_preload", true).Where("id = ?", id).First(&role).Error
	if err != nil {
		err = base.DB.Set("gorm:auto_preload", true).Where("name = ?", id).First(&role).Error
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

func listRole(root treeprint.Tree, role *models.Role) {
	roleString := role.Name
	if role.Target != nil {
		roleString += "(" + color.YellowString(*role.Target) + ")"
	}
	roleNode := root.AddMetaBranch(color.MagentaString("%d", role.ID), roleString)

	for _, perm := range role.Permissions {
		roleNode.AddMetaNode(color.MagentaString("%d", perm.ID), perm.Name)
	}
	return
}

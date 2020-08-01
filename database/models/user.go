package models

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/pkg/errors"
	"time"
)

type User struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5"`
	Nickname string `gorm:"index:nickname" json:"nickname"`
	Email    string `gorm:"unique_index" json:"email"`
	Password string `json:"-"`

	Roles      []UserHasRole `json:"roles"`
	RoleLoaded bool          `gorm:"-"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	//TODO: bio
}

func (u *User) GrantRole(role Role, target ...HasRole) {
	if len(target) == 0 {
		err := base.DB.Model(u).Association("Roles").Append(UserHasRole{
			Role: role,
		}).Error
		if err != nil {
			panic(err)
		}
	} else {
		err := base.DB.Model(u).Association("Roles").Append(UserHasRole{
			Role:     role,
			TargetID: target[0].GetID(),
		}).Error
		if err != nil {
			panic(err)
		}
	}
}

func (u *User) LoadRoles() {
	err := base.DB.Set("gorm:auto_preload", true).Model(u).Related(&u.Roles).Error
	if err != nil {
		panic(err)
	}
}

func (u *User) Can(permission string, target ...HasRole) bool {
	if len(target) > 1 {
		panic(errors.New("target length should be one!"))
	}
	if !u.RoleLoaded {
		u.LoadRoles()
	}
	if len(target) == 0 {
		// Generic permission
		for _, role := range u.Roles {
			if role.Role.Target == nil || *role.Role.Target == "" {
				for _, perm := range role.Role.Permissions {
					if perm.Name == permission || perm.Name == "all" {
						return true
					}
				}
			}
		}
	} else {
		// Specific permisison
		for _, role := range u.Roles {
			if role.Role.Target != nil && *role.Role.Target == target[0].TypeName() && role.TargetID == target[0].GetID() {
				for _, perm := range role.Role.Permissions {
					if perm.Name == permission || perm.Name == "all" {
						return true
					}
				}
			}
		}
	}
	return false
}

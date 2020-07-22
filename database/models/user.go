package models

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/pkg/errors"
	"time"
)

type User struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5"`
	Nickname string `json:"nickname"`
	Email    string `gorm:"unique_index" json:"email"`
	Password string `json:"-"`

	Roles      []UserHasRole `json:"roles"`
	RoleLoaded bool          `gorm:"-"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (u *User) LoadRoles() {
	base.DB.Set("gorm:auto_preload", true).Model(u).Related(&u.Roles)
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
				if role.Role.Name == "admin" {
					return true
				}
				for _, perm := range role.Role.Permissions {
					if perm.Name == permission {
						return true
					}
				}
			}
		}
	} else {
		// Specific permisison
		for _, role := range u.Roles {
			if role.Role.Target != nil && *role.Role.Target == target[0].TypeName() && role.TargetID == target[0].ID() {
				if role.Role.Name == "admin" {
					return true
				}
				for _, perm := range role.Role.Permissions {
					if perm.Name == permission {
						return true
					}
				}
			}
		}
	}
	return false
}

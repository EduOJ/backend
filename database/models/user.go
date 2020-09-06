package models

import (
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/pkg/errors"
	"time"
)

type User struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
	Nickname string `gorm:"index:nickname" json:"nickname"`
	Email    string `gorm:"unique_index" json:"email"`
	Password string `json:"-"`

	Roles      []UserHasRole `json:"roles"`
	RoleLoaded bool          `gorm:"-" json:"-"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	// TODO: bio
}

func (u *User) GrantRole(name string, target ...HasRole) {
	role := getRole(name, target...)
	if len(target) == 0 {
		if err := base.DB.Model(u).Association("Roles").Append(UserHasRole{
			Role: role,
		}).Error; err != nil {
			panic(err)
		}
	} else {
		if err := base.DB.Model(u).Association("Roles").Append(UserHasRole{
			Role:     role,
			TargetID: target[0].GetID(),
		}).Error; err != nil {
			panic(err)
		}
	}
}

func (u *User) DeleteRole(name string, target ...HasRole) {
	role := getRole(name, target...)
	userHasRole := UserHasRole{}
	if len(target) == 0 {
		if err := base.DB.Where("user_id = ? and role_id = ? and target_id = ?", u.ID, role.ID, 0).First(&userHasRole).Error; err == gorm.ErrRecordNotFound {
			return
		} else if err != nil {
			panic(err)
		}
	} else {
		if err := base.DB.Where("user_id = ? and role_id = ? and target_id = ?", u.ID, role.ID, target[0].GetID()).First(&userHasRole).Error; err == gorm.ErrRecordNotFound {
			return
		} else if err != nil {
			panic(err)
		}
	}
	if err := base.DB.Delete(&userHasRole).Error; err != nil {
		panic(err)
	}
}

func (u *User) HasRole(name string, target ...HasRole) bool {
	role := getRole(name, target...)
	userHasRole := UserHasRole{}
	if len(target) == 0 {
		if err := base.DB.Where("user_id = ? and role_id = ? and target_id = ?", u.ID, role.ID, 0).First(&userHasRole).Error; err == gorm.ErrRecordNotFound {
			return false
		} else if err != nil {
			panic(err)
		}
	} else {
		if err := base.DB.Where("user_id = ? and role_id = ? and target_id = ?", u.ID, role.ID, target[0].GetID()).First(&userHasRole).Error; err == gorm.ErrRecordNotFound {
			return false
		} else if err != nil {
			panic(err)
		}
	}
	return true
}

func getRole(name string, target ...HasRole) Role {
	role := Role{}
	if len(target) == 0 {
		if err := base.DB.Where("name = ? and target is null", name).FirstOrCreate(&role).Error; err != nil {
			panic(err)
		}
	} else {
		if err := base.DB.Where("name = ? and target = ?", name, target[0].TypeName()).FirstOrCreate(&role).Error; err != nil {
			panic(err)
		}
	}
	return role
}

func (u *User) LoadRoles() {
	err := base.DB.Set("gorm:auto_preload", true).Model(u).Related(&u.Roles).Error
	if err != nil {
		panic(err)
	}
	u.RoleLoaded = true
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

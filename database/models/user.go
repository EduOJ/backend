package models

import (
	"encoding/binary"
	"encoding/json"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
	Nickname string `gorm:"index:nickname" json:"nickname"`
	Email    string `gorm:"unique_index" json:"email"`
	Password string `json:"-"`

	Roles      []UserHasRole `json:"roles"`
	RoleLoaded bool          `gorm:"-" json:"-"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
	// TODO: bio

	Credentials []WebauthnCredential
}

func (u *User) WebAuthnID() (ret []byte) {
	ret = make([]byte, 8)
	binary.LittleEndian.PutUint64(ret, uint64(u.ID))
	return
}

func (u *User) WebAuthnName() string {
	return u.Username
}

func (u *User) WebAuthnDisplayName() string {
	return u.Nickname
}

func (u *User) WebAuthnIcon() string {
	return ""
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	err := base.DB.Model(u).Association("Credentials").Find(&u.Credentials)
	if err != nil {
		panic(errors.Wrap(err, "could not query user credentials"))
	}
	ret := make([]webauthn.Credential, len(u.Credentials))
	for i, v := range u.Credentials {
		err := json.Unmarshal([]byte(v.Content), &ret[i])
		if err != nil {
			panic(errors.Wrap(err, "wrong json in user's credential"))
		}
	}
	return ret
}

func (u *User) GrantRole(name string, target ...HasRole) {
	role := getRole(name, target...)
	if len(target) == 0 {
		if err := base.DB.Model(u).Association("Roles").Append(&UserHasRole{
			Role: role,
		}); err != nil {
			panic(err)
		}
	} else {
		if err := base.DB.Model(u).Association("Roles").Append(&UserHasRole{
			Role:     role,
			TargetID: target[0].GetID(),
		}); err != nil {
			panic(err)
		}
	}
}

func (u *User) DeleteRole(name string, target ...HasRole) {
	role := getRole(name, target...)
	userHasRole := UserHasRole{}
	if len(target) == 0 {
		if err := base.DB.Where("user_id = ? and role_id = ? and target_id = ?", u.ID, role.ID, 0).First(&userHasRole).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			return
		} else if err != nil {
			panic(err)
		}
	} else {
		if err := base.DB.Where("user_id = ? and role_id = ? and target_id = ?", u.ID, role.ID, target[0].GetID()).First(&userHasRole).Error; errors.Is(err, gorm.ErrRecordNotFound) {
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
		if err := base.DB.Where("user_id = ? and role_id = ? and target_id = ?", u.ID, role.ID, 0).First(&userHasRole).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		} else if err != nil {
			panic(err)
		}
	} else {
		if err := base.DB.Where("user_id = ? and role_id = ? and target_id = ?", u.ID, role.ID, target[0].GetID()).First(&userHasRole).Error; errors.Is(err, gorm.ErrRecordNotFound) {
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
	err := base.DB.Preload("Role.Permissions").Model(u).Association("Roles").Find(&u.Roles)
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

package models

import "github.com/leoleoasd/EduOJBackend/base"

type HasRole interface {
	GetID() uint
	TypeName() string
}

type UserHasRole struct {
	ID       uint `gorm:"primary_key" json:"id"`
	UserID   uint `json:"user_id"`
	RoleID   uint `json:"role_id"`
	Role     Role `json:"role"`
	TargetID uint `json:"target_id"`
}

type Role struct {
	ID          uint    `gorm:"primary_key" json:"id"`
	Name        string  `json:"name"`
	Target      *string `json:"target"`
	Permissions []Permission
}

type Permission struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	RoleID uint   `json:"role_id"`
	Name   string `json:"name"`
}

func (r *Role) AddPermission(name string) {
	p := Permission{
		RoleID: r.ID,
		Name:   name,
	}
	base.DB.Model(r).Association("Permissions").Append(p)
}

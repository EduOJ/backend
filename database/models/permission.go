package models

type HasRole interface {
	ID() uint
	TypeName() string
}

type UserHasRole struct {
	ID       uint `gorm:"primary_key" json:"id"`
	UserID   uint `json:"user_id"`
	RoleID   uint `json:"role_id"`
	Role     `json:"role"`
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

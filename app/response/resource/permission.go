package resource

import "github.com/leoleoasd/EduOJBackend/database/models"

type PermissionProfile struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}

type RoleProfile struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Target      *string             `json:"target"`
	Permissions []PermissionProfile `json:"permissions"`
	TargetID    uint                `json:"target_id"`
}

func (p *PermissionProfile) Convert(perm *models.Permission) {
	p.ID = perm.ID
	p.Name = perm.Name
}

func (p *RoleProfile) Convert(userHasRole *models.UserHasRole) {
	p.Name = userHasRole.Role.Name
	p.Target = userHasRole.Role.Target
	p.TargetID = userHasRole.TargetID
	p.Permissions = make([]PermissionProfile, len(userHasRole.Role.Permissions))
	for i, perm := range userHasRole.Role.Permissions {
		p.Permissions[i].Convert(&perm)
	}
}

func GetPermissionProfile(perm *models.Permission) *PermissionProfile {
	p := PermissionProfile{}
	p.Convert(perm)
	return &p
}

func GetRoleProfile(role *models.UserHasRole) *RoleProfile {
	p := RoleProfile{}
	p.Convert(role)
	return &p
}

func GetRoleProfileSlice(roles []models.UserHasRole) (profiles []RoleProfile) {
	profiles = make([]RoleProfile, len(roles))
	for i, role := range roles {
		profiles[i].Convert(&role)
	}
	return
}

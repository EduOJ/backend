package resource

import "github.com/leoleoasd/EduOJBackend/database/models"

type UserProfile struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
	Nickname string `gorm:"index:nickname" json:"nickname"`
	Email    string `gorm:"unique_index" json:"email"`
}

type UserProfileForAdmin struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
	Nickname string `gorm:"index:nickname" json:"nickname"`
	Email    string `gorm:"unique_index" json:"email"`

	Roles []RoleProfile `json:"roles"`
}

type UserProfileForMe struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
	Nickname string `gorm:"index:nickname" json:"nickname"`
	Email    string `gorm:"unique_index" json:"email"`

	Roles []RoleProfile `json:"roles"`
}

type RoleProfile struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Target      *string             `json:"target"`
	Permissions []PermissionProfile `json:"permissions"`
	TargetID    uint                `json:"target_id"`
}

type PermissionProfile struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}

func (p *PermissionProfile) Convert(perm *models.Permission) {
	p.ID = perm.ID
	p.Name = perm.Name
}

func GetPermissionProfile(perm *models.Permission) *PermissionProfile {
	p := PermissionProfile{}
	p.Convert(perm)
	return &p
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

func (p *UserProfile) Convert(user *models.User) {
	if user == nil {
		return
	}
	p.ID = user.ID
	p.Username = user.Username
	p.Nickname = user.Nickname
	p.Email = user.Email
}

func (p *UserProfileForAdmin) Convert(user *models.User) {
	if user == nil {
		return
	}
	p.ID = user.ID
	p.Username = user.Username
	p.Nickname = user.Nickname
	p.Email = user.Email
	p.Roles = GetRoleProfileSlice(user.Roles)
}

func (p *UserProfileForMe) Convert(user *models.User) {
	if user == nil {
		return
	}
	p.ID = user.ID
	p.Username = user.Username
	p.Nickname = user.Nickname
	p.Email = user.Email
	p.Roles = GetRoleProfileSlice(user.Roles)
}

func GetUserProfile(user *models.User) *UserProfile {
	p := UserProfile{}
	p.Convert(user)
	return &p
}

func GetUserProfileForAdmin(user *models.User) *UserProfileForAdmin {
	p := UserProfileForAdmin{}
	p.Convert(user)
	return &p
}

func GetUserProfileForMe(user *models.User) *UserProfileForMe {
	p := UserProfileForMe{}
	p.Convert(user)
	return &p
}

func GetUserProfileSlice(users []models.User) (profiles []UserProfile) {
	profiles = make([]UserProfile, len(users))
	for i, user := range users {
		profiles[i] = *GetUserProfile(&user)
	}
	return
}

// TODO: test

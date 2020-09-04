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

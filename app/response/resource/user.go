package resource

import "github.com/leoleoasd/EduOJBackend/database/models"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
	Nickname string `gorm:"index:nickname" json:"nickname"`
	Email    string `gorm:"unique_index" json:"email"`
}

type UserForAdmin struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
	Nickname string `gorm:"index:nickname" json:"nickname"`
	Email    string `gorm:"unique_index" json:"email"`

	Roles  []Role  `json:"roles"`
	Scores []Grade `json:"scores"`
}

func (p *User) convert(user *models.User) {
	if user == nil {
		return
	}
	p.ID = user.ID
	p.Username = user.Username
	p.Nickname = user.Nickname
	p.Email = user.Email
}

func (p *UserForAdmin) convert(user *models.User) {
	if user == nil {
		return
	}
	p.ID = user.ID
	p.Username = user.Username
	p.Nickname = user.Nickname
	p.Email = user.Email
	p.Roles = GetRoleSlice(user.Roles)
	p.Scores = GetGradeSlice(user.Scores)
}

func GetUser(user *models.User) *User {
	p := User{}
	p.convert(user)
	return &p
}

func GetUserForAdmin(user *models.User) *UserForAdmin {
	p := UserForAdmin{}
	p.convert(user)
	return &p
}

func GetUserSlice(users []*models.User) (profiles []User) {
	profiles = make([]User, len(users))
	for i, user := range users {
		profiles[i] = *GetUser(user)
	}
	return
}

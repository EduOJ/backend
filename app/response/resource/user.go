package resource

import "github.com/EduOJ/backend/database/models"

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// @description UserForAdmin is a user with additional, credential data, only accessible by people has permission,
// @description e.g. admin can access to all user's credential data, and a user can access to his/her credential data.
type UserForAdmin struct {
	// ID is the user's id.
	ID uint `json:"id"`
	// Username is the user's username, usually the student ID if used in schools.
	Username string `json:"username"`
	// Nickname is the user's nickname, usually the student name if used in schools.
	Nickname string `json:"nickname"`
	// Email is the user's email.
	Email string `json:"email"`

	// Role is the user's role, and is used to obtain the permissions of a user.
	Roles []Role `json:"roles"`
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

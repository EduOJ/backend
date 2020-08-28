package response

type UserProfile struct {
	//TODO: add this to responses
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
	Nickname string `gorm:"index:nickname" json:"nickname"`
	Email    string `gorm:"unique_index" json:"email"`
}

type UserProfileForAdmin struct {
}

type RoleProfile struct {
}

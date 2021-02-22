package models

type Class struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name" gorm:"size:255;default:'';not null"`
	CourseName  string `json:"course_name" gorm:"size:255;default:'';not null"`
	Description string `json:"description"`
	InviteCode  string `json:"invite_code" gorm:"size:255;default:'';not null"`
	Managers    []User `json:"managers" gorm:"many2many:user_manage_classes"`
	Students    []User `json:"students" gorm:"many2many:user_in_classes"`
}

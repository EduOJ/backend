package models

import (
	"github.com/leoleoasd/EduOJBackend/database"
	"time"
)

type Language struct {
	Name             string               `gorm:"primaryKey"`
	ExtensionAllowed database.StringArray `gorm:"type:string"`
	BuildScriptName  string
	BuildScript      *Script `gorm:"foreignKey:BuildScriptName"`
	RunScriptName    string
	RunScript        *Script `gorm:"foreignKey:RunScriptName"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

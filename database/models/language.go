package models

import (
	"github.com/leoleoasd/EduOJBackend/database"
	"time"
)

type Language struct {
	Name             string               `gorm:"primaryKey" json:"name"`
	ExtensionAllowed database.StringArray `gorm:"type:string" json:"extension_allowed"`
	BuildScriptName  string               `json:"build_script_name"`
	BuildScript      *Script              `gorm:"foreignKey:BuildScriptName" json:"-"`
	RunScriptName    string               `json:"run_script_name"`
	RunScript        *Script              `gorm:"foreignKey:RunScriptName" json:"-"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

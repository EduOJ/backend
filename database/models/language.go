package models

import "time"

type Language struct {
	Name            string `gorm:"primaryKey"`
	BuildScriptName string
	BuildScript     *Script `gorm:"foreignKey:BuildScriptName"`
	RunScriptName   string
	RunScript       *Script `gorm:"foreignKey:RunScriptName"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

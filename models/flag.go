package models

import "gorm.io/gorm"

type Flag struct {
	gorm.Model
	Name        string
	Key         string
	Type        string
	IsEnabled   bool
	Variations  []Variation `gorm:"foreignKey:FlagID"`
	WorkspaceID int64
	Workspace   Workspace
}

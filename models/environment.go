package models

import "gorm.io/gorm"

type Environment struct {
	gorm.Model
	Name        string
	WorkspaceID int64
	Workspace   Workspace
	Targets     []Target `gorm:"foreignKey:EnvironmentID"`
}

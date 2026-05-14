package models

import "gorm.io/gorm"

type Flag struct {
	gorm.Model
	Name        string
	Key         string
	Type        string
	WorkspaceID int64
	Workspace   Workspace
}

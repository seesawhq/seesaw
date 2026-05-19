package models

import "gorm.io/gorm"

type Workspace struct {
	gorm.Model
	ID           int64 `gorm:"primarykey"`
	Name         string
	Flags        []Flag        `gorm:"foreignKey:WorkspaceID"`
	Environments []Environment `gorm:"foreignKey:WorkspaceID"`
	Members      []Member      `gorm:"foreignKey:WorkspaceID"`
}

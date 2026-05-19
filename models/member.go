package models

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	WorkspaceID int64
	Workspace   Workspace
	UserID      int64
	User        User
	Role        string
}

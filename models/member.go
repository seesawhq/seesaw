package models

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	ID     int64 `gorm:"primarykey"`
	TeamID int64
	Team   Workspace
	UserID int64
	User   User
	Role   string
}

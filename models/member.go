package models

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	ID     int64
	TeamID int64
	Team   Team
	UserID int64
	User   User
	Role   string
}

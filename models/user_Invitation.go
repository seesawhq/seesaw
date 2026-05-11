package models

import (
	"gorm.io/gorm"
)

type UserInvitation struct {
	gorm.Model
	ID        int64  `gorm:"primarykey"`
	Email     string `gorm:"unique;not null"`
	Token     string `gorm:"unique;not null"`
	ValidTill int64
}

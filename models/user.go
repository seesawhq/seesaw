package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        int64 `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
}

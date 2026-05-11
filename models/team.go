package models

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	ID   int64 `gorm:"primarykey"`
	Name string
}

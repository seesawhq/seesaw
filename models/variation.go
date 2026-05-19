package models

import "gorm.io/gorm"

type Variation struct {
	gorm.Model
	Name             string
	Value            string
	IsDisableDefault bool
	FlagID           int64
	Flag             Flag
}

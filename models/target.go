package models

import "gorm.io/gorm"

type Target struct {
	gorm.Model
	Name          string
	EnvironmentID int64
	Environment   Environment
	FlagID        int64
	Flag          Flag
	Type          string
	ServeType     string
	Rollouts      []Rollout `gorm:"foreignKey:TargetID"`
}

type Rollout struct {
	gorm.Model
	TargetID    int64
	Target      Target
	Percentage  int64
	VariationID int64
	Variation   Variation
	IsControl   bool
}

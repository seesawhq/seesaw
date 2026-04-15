package models

type FeatureFlag struct {
	Name       string
	TargetType string
	TeamID     int64
	Team       Team
}

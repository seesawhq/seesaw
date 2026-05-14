package models

type Environment struct {
	Name        string
	TargetType  string
	WorkspaceID int64
	Workspace   Workspace
}

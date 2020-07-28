package models

import (
	"time"
)

// Assignment represents a policy being applied/attached to a role
type Assignment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	ProjectID string    `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Assignment) GetEncryptable() bool {
	return false
}

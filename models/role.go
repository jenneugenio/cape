package models

import (
	"time"
)

const (
	// AdminRole is the label of the admin role
	AdminRole = Label("admin")

	// UserRole is the label of the global role
	UserRole = Label("global")

	ProjectOwnerRole       = Label("project-owner")
	ProjectContributorRole = Label("project-contributor")
	ProjectReaderRole      = Label("project-reader")
)

var OrgRoles = []Label{AdminRole, UserRole}
var ProjectRoles = []Label{ProjectOwnerRole, ProjectContributorRole, ProjectReaderRole}

func ValidOrgRole(role Label) bool {
	for _, r := range OrgRoles {
		if role == r {
			return true
		}
	}

	return false
}

func ValidProjectRole(role Label) bool {
	for _, r := range ProjectRoles {
		if role == r {
			return true
		}
	}

	return false
}

// Role in a role in the system (e.g. Admin, user, etc)
type Role struct {
	ID        string    `json:"id"`
	Version   uint8     `json:"version"`
	Label     Label     `json:"label"`
	System    bool      `json:"system"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewRole returns a mutable role struct
func NewRole(label Label, system bool) Role {
	return Role{
		ID:        NewID(),
		Version:   modelVersion,
		Label:     label,
		System:    system,
		CreatedAt: now(),
	}
}

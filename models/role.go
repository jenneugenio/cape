package models

import (
	"time"
)

type Permission uint64

const (
	WritePolicy Permission = 1 << iota
	CreateProject

	// Tokens
	CreateOwnToken
	CreateAnyToken
	RemoveOwnToken
	RemoveAnyToken
	ListOwnTokens
	ListAnyTokens

	// Projects
	ArchiveProject
	UnarchiveProject
	DeleteOwnedProject
	DeleteAnyProject

	AddUser
	DeleteUser
	UpdateProject
	SuggestPolicy
	AcceptPolicy
	ReadPolicy

	// Roles
	ChangeRole
)

type Principal interface{}

const (
	// AdminRole is the label of the admin role
	AdminRole = Label("admin")

	// UserRole is the label of the global role
	UserRole = Label("user")

	ProjectOwnerRole       = Label("project-owner")
	ProjectContributorRole = Label("project-contributor")
	ProjectReaderRole      = Label("project-reader")
)

var (
	adminRules = withRules(
		WritePolicy, CreateProject, AddUser, DeleteUser,

		CreateOwnToken, RemoveOwnToken, ListOwnTokens,
		CreateAnyToken, RemoveAnyToken, ListAnyTokens,

		DeleteAnyProject,

		ChangeRole,
	)

	userRules = withRules(
		WritePolicy, CreateProject,

		CreateOwnToken, RemoveOwnToken, ListOwnTokens,
	)

	projectReaderRules = withRules(
		ReadPolicy,
	)

	projectContributorRules = withRules(
		projectReaderRules, UpdateProject, SuggestPolicy,
	)

	projectOwnerRules = withRules(
		projectContributorRules, AcceptPolicy, UnarchiveProject, DeleteOwnedProject,
	)

	DefaultPermissions = map[Label]Permission{
		AdminRole:              adminRules,
		UserRole:               userRules,
		ProjectOwnerRole:       projectOwnerRules,
		ProjectContributorRole: projectContributorRules,
		ProjectReaderRole:      projectReaderRules,
	}
)

func withRules(perms ...Permission) Permission {
	var p Permission
	for _, perm := range perms {
		p |= perm
	}

	return p
}

var OrgRoles = []Label{AdminRole, UserRole}
var ProjectRoles = []Label{ProjectOwnerRole, ProjectContributorRole, ProjectReaderRole}

var SystemRoles = []Label{AdminRole, UserRole, ProjectReaderRole, ProjectContributorRole, ProjectOwnerRole}

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

type ProjectRolesMap map[Label]Role

func (p ProjectRolesMap) Get(l Label) Role {
	return p[l]
}

type UserRoles struct {
	Global   Role
	Projects ProjectRolesMap
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

func (r *Role) Can(action Permission) bool {
	return DefaultPermissions[r.Label]&action != 0
}

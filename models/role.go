package models

import (
	"fmt"
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
	RejectPolicy
	ReadPolicy
	ListPolicySuggestions

	// Roles
	ChangeRole
	ChangeProjectRole
)

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
		projectReaderRules, UpdateProject, SuggestPolicy, ListPolicySuggestions, RejectPolicy,
	)

	projectOwnerRules = withRules(
		projectContributorRules, AcceptPolicy, UnarchiveProject, DeleteOwnedProject, ChangeProjectRole,
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

// OrgRoles are roles that can occur outside the scope of a project. There are currently only
// admin and user roles.
var OrgRoles = []Label{AdminRole, UserRole}

// ProjectRoles are roles that are only related to projects. Currently there is a project
// owner, a contributor and a reader.
var ProjectRoles = []Label{ProjectOwnerRole, ProjectContributorRole, ProjectReaderRole}

// SystemRoles are all builtin roles
var SystemRoles = append(OrgRoles, ProjectRoles...)

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

func (p ProjectRolesMap) Get(l Label) (*Role, error) {
	r, ok := p[l]
	if !ok {
		return nil, fmt.Errorf("unable to find membership in requested project")
	}

	return &r, nil
}

// UserRoles represents the roles assigned to a user. A user
// can only have one global role and then one project role per project
// that they are a member of.
type UserRoles struct {
	// Global is the global role assigned to a user
	Global Role

	// Projects is a map between a projects Label and the role they have
	// in that project.
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

// Can checks to see if a role can do an action
func (r *Role) Can(action Permission) bool {
	return DefaultPermissions[r.Label]&action != 0
}

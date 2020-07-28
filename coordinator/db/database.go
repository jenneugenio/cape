package db

import (
	"context"
	"errors"

	"github.com/capeprivacy/cape/models"
)

type Interface interface {
	Policies() PolicyDB
	Roles() RoleDB
	Users() UserDB
	Projects() ProjectsDB
	Contributors() ContributorDB
	Config() ConfigDB
}

// Interfaces

type PolicyDB interface {
	Create(context.Context, models.Policy) error
	Delete(context.Context, models.Label) (DeleteStatus, error)
	Get(context.Context, models.Label) (*models.Policy, error)
	GetByID(context.Context, string) (*models.Policy, error)
	List(ctx context.Context, opts *ListPolicyOptions) ([]models.Policy, error)
}

type UserDB interface {
	Create(context.Context, models.User) error
	Update(context.Context, string, models.User) error
	Delete(context.Context, models.Email) (DeleteStatus, error)
	Get(context.Context, models.Email) (*models.User, error)
	GetByID(context.Context, string) (*models.User, error)
	List(context.Context, *ListUserOptions) ([]models.User, error)
}

type RoleDB interface {
	Get(context.Context, models.Label) (*models.Role, error)
	GetByID(context.Context, string) (*models.Role, error)
	List(context.Context, *ListRoleOptions) ([]*models.Role, error)

	GetAll(context.Context, string) (*models.UserRoles, error)

	SetOrgRole(context.Context, models.Email, models.Label) (*models.Assignment, error)
	GetOrgRole(context.Context, models.Email) (*models.Role, error)
	SetProjectRole(context.Context, models.Email, models.Label, models.Label) (*models.Assignment, error)
	GetProjectRole(context.Context, models.Email, string) (*models.Role, error)

	CreateSystemRoles(context.Context) error
}

type ConfigDB interface {
	Create(context.Context, models.Config) error
	Get(context.Context) (*models.Config, error)
}

type ContributorDB interface {
	Add(context.Context, models.Label, models.Email) (*models.Contributor, error)
	Get(context.Context, models.Label, models.Email) (*models.Contributor, error)
	List(context.Context, models.Label) ([]models.Contributor, error)
	Delete(context.Context, models.Label, models.Email) (*models.Contributor, error)
}

type ProjectsDB interface {
	Get(context.Context, models.Label) (*models.Project, error)
	GetByID(context.Context, string) (*models.Project, error)
	Create(context.Context, models.Project) error
	Update(context.Context, models.Project) error

	List(context.Context) ([]models.Project, error)
	ListByStatus(context.Context, models.ProjectStatus) ([]models.Project, error)

	CreateProjectSpec(context.Context, models.ProjectSpec) error
	GetProjectSpec(context.Context, string) (*models.ProjectSpec, error)
}

// Options

type ListPolicyOptions struct {
	Options *struct {
		Offset uint64
		Limit  uint64
	}

	FilterIDs []string
}

type ListRoleOptions struct {
	Offset uint64
	Limit  uint64
}

type ListUserOptions struct {
	Options *struct {
		Offset uint64
		Limit  uint64
	}

	FilterIDs []string
}

// Statuses

type DeleteStatus string

const (
	DeleteStatusDeleted      DeleteStatus = "deleted"
	DeleteStatusDoesNotExist DeleteStatus = "does_not_exist"
	DeleteStatusError        DeleteStatus = "error"
)

// Errors

var ErrDuplicateKey = errors.New("duplicate key")
var ErrNoRows = errors.New("no rows")

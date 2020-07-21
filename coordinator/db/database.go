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
	RBAC() RBACDB
}

// Interfaces

type PolicyDB interface {
	Create(context.Context, models.Policy) error
	Delete(context.Context, models.Label) (DeleteStatus, error)
	Get(context.Context, models.Label) (*models.Policy, error)
	GetByID(context.Context, string) (*models.Policy, error)
	List(ctx context.Context, opts *ListPolicyOptions) ([]models.Policy, error)
}

type RBACDB interface {
	Create(context.Context, models.RBACPolicy) error
	List(ctx context.Context, opts *ListRBACOptions) ([]models.RBACPolicy, error)
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
	Create(context.Context, *models.Role) error
	Delete(context.Context, models.Label) (DeleteStatus, error)
	Get(context.Context, models.Label) (*models.Role, error)
	GetByID(context.Context, string) (*models.Role, error)
	List(context.Context, *ListRoleOptions) ([]*models.Role, error)

	AttachPolicy(context.Context, models.Label) error
	DetachPolicy(context.Context, models.Label) error
}

type ContributorDB interface {
	Add(context.Context, models.Label, models.Email, models.Label) (*models.Contributor, error)
	Get(context.Context, models.Label, models.Email) (*models.Contributor, error)
	List(context.Context, models.Label) ([]models.Contributor, error)
	Delete(context.Context, models.Label, models.Email) (*models.Contributor, error)
}

type ProjectsDB interface {
	Get(context.Context, models.Label) (*models.Project, error)
	GetByID(context.Context, string) (*models.Project, error)
}

// Options
type ListPolicyOptions struct {
	Options *struct {
		Offset uint64
		Limit  uint64
	}

	FilterIDs []string
}

type ListRBACOptions struct {
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

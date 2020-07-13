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
}

type PolicyDB interface {
	Create(context.Context, models.Policy) error
	Delete(context.Context, models.Label) error
	Get(context.Context, models.Label) (*models.Policy, error)
	GetByID(context.Context, string) (*models.Policy, error)
	List(ctx context.Context, opts *ListPolicyOptions) ([]models.Policy, error)
}

type UserDB interface {
	Create(context.Context, models.User) error
	Update(context.Context, string, models.User) error
	Delete(context.Context, models.Email) error
	Get(context.Context, models.Email) (*models.User, error)
	GetByID(context.Context, string) (*models.User, error)
	List(context.Context, *ListUserOptions) ([]models.User, error)
}

type RoleDB interface {
	Create(context.Context, *models.Role) error
	Delete(context.Context, models.Label) error
	Get(context.Context, models.Label) (*models.Role, error)
	List(context.Context, *ListRoleOptions) ([]*models.Role, error)

	AttachPolicy(context.Context, models.Label) error
	DetachPolicy(context.Context, models.Label) error
}

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

var ErrDuplicateKey = errors.New("duplicate key")
var ErrNoRows = errors.New("no rows")

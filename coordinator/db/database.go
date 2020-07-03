package db

import (
	"context"

	"github.com/capeprivacy/cape/models"
)

type Interface interface {
	Policies() PolicyDB
	Roles() RoleDB
}

type PolicyDB interface {
	Create(context.Context, *models.Policy) error
	Delete(context.Context, models.Label) error
	Get(context.Context, models.Label) (*models.Policy, error)
	List(ctx context.Context, opts *ListPolicyOptions) ([]*models.Policy, error)
}

type RoleDB interface {
	Create(context.Context, *models.Policy) (*models.Policy, error)
	Delete(context.Context, models.Label) error
	Get(context.Context, models.Label) (*models.Policy, error)
	List(context.Context, *ListPolicyOptions) ([]*models.Policy, error)

	AttachPolicy(context.Context, models.Label) error
	DetachPolicy(context.Context, models.Label) error
}

type ListPolicyOptions struct {
	Offset uint64
	Limit uint64
}

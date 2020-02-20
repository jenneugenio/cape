package database

import (
	"context"

	"github.com/dropoutlabs/privacyai/primitives"
)

// Querier represents various backend queries you can perform
type Querier interface {
	Create(context.Context, *primitives.Primitive) error
	Get(context.Context, primitives.ID) (*primitives.Primitive, error)
	Delete(context.Context, *primitives.Primitive) error
	Update(context.Context, *primitives.Primitive) error
}

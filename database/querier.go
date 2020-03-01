package database

import (
	"context"

	"github.com/dropoutlabs/privacyai/primitives"
)

// Querier represents various backend queries you can perform
type Querier interface {
	Create(context.Context, ...primitives.Entity) error
	Get(context.Context, primitives.ID, primitives.Entity) error
	Delete(context.Context, primitives.ID) error
	Update(context.Context, primitives.Entity) error
	QueryOne(context.Context, primitives.Entity, Filter) error
	Query(context.Context, interface{}, Filter) error
}

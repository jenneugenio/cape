package database

import (
	"context"

	"github.com/capeprivacy/cape/coordinator/database/types"
)

// Querier represents various backend queries you can perform
type Querier interface {
	Create(context.Context, ...Entity) error
	Get(context.Context, ID, Entity) error
	Delete(ctx context.Context, typ types.Type, id ID) error
	Update(context.Context, Entity) error
	QueryOne(context.Context, Entity, Filter) error
	Query(context.Context, interface{}, Filter) error
}

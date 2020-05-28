package database

import (
	"context"

	"github.com/capeprivacy/cape/coordinator/database/types"
)

// Querier represents various backend queries you can perform
type Querier interface {
	Create(context.Context, ...Entity) error
	Get(context.Context, ID, Entity) error
	Delete(context.Context, types.Type, ID) error
	Upsert(context.Context, Entity) error
	Update(context.Context, Entity) error
	SubQueryOne(context.Context, Entity, *Select, Filter) error
	QueryOne(context.Context, Entity, Filter) error
	Query(context.Context, interface{}, Filter) error
}

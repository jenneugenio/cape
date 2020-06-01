package auth

import (
	"context"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// A Canner verifies whether the owner of a c can perform a specific Action on a given Type
type Canner interface {
	Can(primitives.Action, types.Type) error
}

// Enforcer enforces authorization for accessing primitive types tables.
// This required by the graphql resolvers so that they can check to see
// if the requesting user can access the primitive tables.
//
// Example usage:
//
// func (r *Resolver) resolver(ctx context, id database.ID) SomeData {
//    c := framework.Session(ctx)
//
//    enforcer := NewEnforcer(c, r.Backend)
//
//    return enforcer.Get(ctx, id)
// }
//
type Enforcer struct {
	c  Canner
	db database.Querier
}

// NewEnforcer creates a new enforcer
func NewEnforcer(c Canner, db database.Querier) *Enforcer {
	return &Enforcer{
		c:  c,
		db: db,
	}
}

// Create calls down to the underlying db function as long as the contained policies
// can create the given entities.
func (e *Enforcer) Create(ctx context.Context, entity ...database.Entity) error {
	if len(entity) == 0 {
		return errors.New(errors.InvalidArgumentCause, "cannot create 0 entities")
	}

	err := e.c.Can(primitives.Create, entity[0].GetType())
	if err != nil {
		return err
	}

	err = e.db.Create(ctx, entity...)
	if err != nil {
		return err
	}

	return nil
}

// Get calls down to the underlying db function as long as the contained policies
// can query the given entities
func (e *Enforcer) Get(ctx context.Context, id database.ID, entity database.Entity) error {
	if entity == nil {
		return errors.New(errors.InvalidArgumentCause, "entity cannot be nil")
	}

	err := e.c.Can(primitives.Read, entity.GetType())
	if err != nil {
		return err
	}

	err = e.db.Get(ctx, id, entity)
	if err != nil {
		return err
	}

	return nil
}

// Delete calls down to the underlying db function as long as the contained policies
// can delete the given entity
func (e *Enforcer) Delete(ctx context.Context, typ types.Type, id database.ID) error {
	err := e.c.Can(primitives.Delete, typ)
	if err != nil {
		return err
	}

	err = e.db.Delete(ctx, typ, id)
	if err != nil {
		return err
	}

	return nil
}

// Upsert calls down to the underlying db function as long as the contained policies
// can update AND create the given entity
func (e *Enforcer) Upsert(ctx context.Context, entity database.Entity) error {
	if entity == nil {
		return errors.New(errors.InvalidArgumentCause, "entity cannot be nil")
	}

	err := e.c.Can(primitives.Update, entity.GetType())
	if err != nil {
		return err
	}

	err = e.c.Can(primitives.Create, entity.GetType())
	if err != nil {
		return err
	}

	err = e.db.Upsert(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

// Update calls down to the underlying db function as long as the contained policies
// can update the given entity
func (e *Enforcer) Update(ctx context.Context, entity database.Entity) error {
	if entity == nil {
		return errors.New(errors.InvalidArgumentCause, "entity cannot be nil")
	}

	err := e.c.Can(primitives.Update, entity.GetType())
	if err != nil {
		return err
	}

	err = e.db.Update(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

// QueryOne calls down to the underlying db function as long as the contained policies
// can query the given entity
func (e *Enforcer) QueryOne(ctx context.Context, entity database.Entity, filter database.Filter) error {
	if entity == nil {
		return errors.New(errors.InvalidArgumentCause, "entity cannot be nil")
	}

	err := e.c.Can(primitives.Read, entity.GetType())
	if err != nil {
		return err
	}

	err = e.db.QueryOne(ctx, entity, filter)
	if err != nil {
		return err
	}

	return nil
}

// Query calls down to the underlying db function as long as the contained policies
// can query the given entities
func (e *Enforcer) Query(ctx context.Context, i interface{}, filter database.Filter) error {
	typ := database.EntityTypeFromPtrSlice(i)
	err := e.c.Can(primitives.Read, typ)
	if err != nil {
		return err
	}

	err = e.db.Query(ctx, i, filter)
	if err != nil {
		return err
	}

	return nil
}

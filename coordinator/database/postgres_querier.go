package database

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"

	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type pgConn interface {
	Query(ctx context.Context, q string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, q string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, q string, args ...interface{}) (pgconn.CommandTag, error)
}

type postgresQuerier struct {
	conn pgConn
}

// Create an entity inside the database
func (q *postgresQuerier) Create(ctx context.Context, entities ...Entity) error {
	if len(entities) == 0 {
		return nil
	}

	t := entities[0].GetType()
	inserts, values := buildInsert(entities, t)

	sql := fmt.Sprintf(`INSERT INTO %s (data) VALUES %s`, t.String(), inserts)
	_, err := q.conn.Exec(ctx, sql, values...)
	switch e := err.(type) {
	case *pgconn.PgError:
		if e.Code == "23505" {
			return errors.New(DuplicateCause, "entity already exists")
		}

		return err
	default:
		return err
	}
}

// Get an entity from the database
func (q *postgresQuerier) Get(ctx context.Context, id ID, e Entity) error {
	t, err := id.Type()
	if err != nil {
		return err
	}

	// We return an error here instead of panic as it's possible the ID came
	// from outside the system. However, this really isn't where this should be
	// checked! This should be handled _prior_ to getting to this part of the
	// code.
	if t != e.GetType() {
		return errors.New(TypeMismatchCause, "id and entity do not match: %s - %s",
			t.String(), e.GetType().String())
	}

	sql := fmt.Sprintf(`SELECT data FROM %s WHERE id = $1 LIMIT 1`, t.String())
	r := q.conn.QueryRow(ctx, sql, id.String())

	switch err := r.Scan(e); err {
	case pgx.ErrNoRows:
		return errors.New(NotFoundCause, "could not find %s: %s", t.String(), id)
	default:
		return err
	}
}

// QueryOne uses a query to return a single entity from the database
func (q *postgresQuerier) QueryOne(ctx context.Context, e Entity, f Filter) error {
	if f.Page != nil {
		panic("Pagination cannot be performed via a QueryOne")
	}

	f.Page = &Page{Limit: 1}

	t := e.GetType()
	where, params, err := buildFilter(f)
	if err != nil {
		if err == ErrEmptyIn {
			return errors.New(NotFoundCause, "could not find %s", t.String())
		}

		return err
	}

	sql := fmt.Sprintf(`SELECT data from %s %s`, t.String(), where)
	r := q.conn.QueryRow(ctx, sql, params...)

	switch err := r.Scan(e); err {
	case pgx.ErrNoRows:
		return errors.New(NotFoundCause, "could not find %s", t.String())
	default:
		return err
	}
}

// Query retrieves entities of a single type from the database
func (q *postgresQuerier) Query(ctx context.Context, arr interface{}, f Filter) error {
	// We want to use reflect to check whether or not the underlying type of
	// arr is a pointer to a slice or not.
	arrPtr := reflect.ValueOf(arr)
	if arrPtr.Kind() != reflect.Ptr {
		panic("Expected arr to be a pointer to a slice")
	}

	arrValue := arrPtr.Elem()
	if arrValue.Kind() != reflect.Slice {
		panic("Expected arr to be a pointer to a slice")
	}

	// Now we need to figure out the underlying concrete type and ensure that
	// it satisfies the primitive.Entity interface. If it doesn't then we've
	// encountered a developer error!
	//
	// To do this, we create an instance of the underlying type of this slice.
	// This can then be used later on to determine the actual table to query.
	entityType := reflect.TypeOf((*Entity)(nil)).Elem()
	itemType := reflect.New(arrValue.Type().Elem())
	if !itemType.Type().Implements(entityType) && itemType.Kind() == reflect.Ptr {
		itemType = reflect.New(arrValue.Type().Elem().Elem())
	}
	e := itemType.Interface().(Entity)

	where, params, err := buildFilter(f)
	if err != nil {
		// If we got back an empty in then we will treat this as a no-operation
		// as the query is not asking for any data back to be returned.
		if err == ErrEmptyIn {
			return nil
		}

		return err
	}

	sql := fmt.Sprintf(`SELECT data FROM %s %s`, e.GetType().String(), where)
	rows, err := q.conn.Query(ctx, sql, params...)
	if err != nil {
		// XXX Come back and figure out what errors can be returned here and
		// then mutate them appropriately.
		return err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		// Grow the slice
		item := getItem(arrValue, i)
		err = rows.Scan(item)
		if err != nil {
			return err
		}
	}

	return rows.Err()
}

// Delete an entity from the database
func (q *postgresQuerier) Delete(ctx context.Context, typ types.Type, id ID) error {
	t, err := id.Type()
	if err != nil {
		return err
	}

	if t != typ {
		return errors.New(TypeMismatchCause, "Type of ID (%s) does not match specified type (%s)", t, typ)
	}

	sql := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, t)
	ct, err := q.conn.Exec(ctx, sql, id.String())
	if err != nil {
		return err
	}

	if ct.RowsAffected() != 1 {
		return errors.New(NotFoundCause, "could not find %s: %s", t.String(), id)
	}

	return nil
}

// Update an entity inside the database
func (q *postgresQuerier) Update(ctx context.Context, e Entity) error {
	t := e.GetType()
	if !t.Mutable() {
		panic("Cannot update an immutable entity")
	}

	err := e.SetUpdatedAt(time.Now().UTC())
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`UPDATE %s SET data = $1 WHERE id = $2`, t.String())
	ct, err := q.conn.Exec(ctx, sql, e, e.GetID().String())
	if err != nil {
		return err
	}

	if ct.RowsAffected() != 1 {
		return errors.New(NotFoundCause, "could not find %s: %s", t.String(), e.GetID())
	}

	return nil
}

// getItem returns a primitive.Entity from the given slice thats passed as a
// reflect.Value. We then use this value to determine if/how we should grow the
// slice.
func getItem(v reflect.Value, pos int) interface{} {
	if v.Type().Kind() != reflect.Slice {
		panic("expected a slice")
	}

	num := pos + 1 // pos is a 0-indexed position in slice
	if num >= v.Cap() {
		cap := v.Cap() + v.Cap()/2
		if cap < 4 {
			cap = 4
		}

		newV := reflect.MakeSlice(v.Type(), v.Len(), cap)
		reflect.Copy(newV, v)
		v.Set(newV)
	}

	if num >= v.Len() {
		v.SetLen(num)
	}

	return v.Index(pos).Addr().Interface()
}

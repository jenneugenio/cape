package database

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
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

	// codec can be nil. A nil codec means no encryption will be performed
	codec crypto.EncryptionCodec
}

// Create an entity inside the database
func (q *postgresQuerier) Create(ctx context.Context, entities ...Entity) error {
	if len(entities) == 0 {
		return nil
	}

	t := entities[0].GetType()
	inserts, values := buildInsert(entities, t)

	if q.codec == nil && entities[0].GetEncryptable() {
		return errors.New(NoEncryptionCodec,
			"No encryption codec was found but encountered encrytable primitives %s", entities[0].GetType())
	} else if entities[0].GetEncryptable() {
		v, err := handleEncrypt(ctx, q.codec, entities...)
		if err != nil {
			return err
		}
		values = v
	}

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

	var bytes []byte
	err = r.Scan(&bytes)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return errors.New(NotFoundCause, "could not find %s: %s", t.String(), id)
		default:
			return err
		}
	}

	if q.codec == nil && e.GetEncryptable() {
		return errors.New(NoEncryptionCodec,
			"No encryption codec was found but encountered encrytable primitives %s", e.GetType())
	} else if e.GetEncryptable() {
		encryptable := e.(interface{}).(crypto.Encryptable)
		return handleDecrypt(ctx, q.codec, bytes, encryptable)
	}

	return json.Unmarshal(bytes, e)
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
			return errors.New(NotFoundCause, "could not find %s entity with filter %s", t.String(), f.Where)
		}

		return err
	}

	sql := fmt.Sprintf(`SELECT data from %s %s`, t.String(), where)
	r := q.conn.QueryRow(ctx, sql, params...)

	var bytes []byte
	err = r.Scan(&bytes)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return errors.New(NotFoundCause, "could not find %s entity with filter %s", t.String(), f.Where)
		default:
			return err
		}
	}

	if q.codec == nil && e.GetEncryptable() {
		return errors.New(NoEncryptionCodec,
			"No encryption codec was found but encountered encrytable primitives %s", e.GetType())
	} else if e.GetEncryptable() {
		encryptable := e.(interface{}).(crypto.Encryptable)
		return handleDecrypt(ctx, q.codec, bytes, encryptable)
	}

	return json.Unmarshal(bytes, e)
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
		var bytes []byte
		err = rows.Scan(&bytes)
		if err != nil {
			return err
		}

		// Grow the slice
		item := getItem(arrValue, i)

		if q.codec == nil && e.GetEncryptable() {
			return errors.New(NoEncryptionCodec,
				"No encryption codec was found but encountered encrytable primitives %s", e.GetType())
		} else if e.GetEncryptable() {
			encryptable := item.(interface{}).(crypto.Encryptable)

			err := handleDecrypt(ctx, q.codec, bytes, encryptable)
			if err != nil {
				return err
			}
			continue
		}

		err := json.Unmarshal(bytes, item)
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

type updateMode string

const (
	Update updateMode = "UPDATE"
	Upsert updateMode = "UPSERT"
)

func (q *postgresQuerier) update(ctx context.Context, e Entity, mode updateMode) error {
	t := e.GetType()
	if !t.Mutable() {
		panic("Cannot update an immutable entity")
	}

	err := e.SetUpdatedAt(time.Now().UTC())
	if err != nil {
		return err
	}

	var value interface{} = e
	if q.codec == nil && e.GetEncryptable() {
		return errors.New(NoEncryptionCodec,
			"No encryption codec was found but encountered encrytable primitives %s", e.GetType())
	} else if e.GetEncryptable() {
		v, err := handleEncrypt(ctx, q.codec, e)
		if err != nil {
			return err
		}
		value = v[0]
	}

	sql := fmt.Sprintf(`%s %s SET data = $1 WHERE id = $2`, mode, t.String())
	ct, err := q.conn.Exec(ctx, sql, value, e.GetID().String())
	if err != nil {
		return err
	}

	if ct.RowsAffected() != 1 {
		return errors.New(NotFoundCause, "could not find %s: %s", t.String(), e.GetID())
	}

	return nil
}

// Upsert (update if it exists, insert if not) an entity inside the database
func (q *postgresQuerier) Upsert(ctx context.Context, e Entity) error {
	return q.update(ctx, e, Upsert)
}

// Update an entity inside the database
func (q *postgresQuerier) Update(ctx context.Context, e Entity) error {
	return q.update(ctx, e, Update)
}

// getItem returns a primitive.Entity from the given slice thats passed as a
// reflect.Value. We then use this value to determine if/how we should grow the
// slice.
func getItem(v reflect.Value, pos int) Entity {
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

	// The underlying data in the array can either
	// be a pointer to a primitive or just a primitive.
	// If its a pointer we need to allocate a pointer before
	// returning the interface
	if v.Index(pos).Kind() == reflect.Ptr {
		v.Index(pos).Set(reflect.New(v.Index(pos).Type().Elem()))
		return v.Index(pos).Interface().(Entity)
	}
	return v.Index(pos).Addr().Interface().(Entity)
}

func EntityTypeFromPtrSlice(arr interface{}) types.Type {
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

	return itemType.Interface().(Entity).GetType()
}

// handleEncrypt requires any caller to make sure that entities are Encryptable.
// Will panic if not.
func handleEncrypt(ctx context.Context, codec crypto.EncryptionCodec, entities ...Entity) ([]interface{}, error) {
	encrypted := make([]interface{}, len(entities))
	for i, e := range entities {
		encryptable := e.(interface{}).(crypto.Encryptable)
		by, err := encryptable.Encrypt(ctx, codec)
		if err != nil {
			return nil, err
		}

		// pgx can understand json input as byte array
		// or structs with json tags
		encrypted[i] = by
	}

	return encrypted, nil
}

func handleDecrypt(ctx context.Context, codec crypto.EncryptionCodec, input []byte,
	encryptable crypto.Encryptable) error {
	err := encryptable.Decrypt(ctx, codec, input)
	if err != nil {
		return err
	}

	return nil
}

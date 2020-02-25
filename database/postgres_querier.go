package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"

	errors "github.com/dropoutlabs/privacyai/partyerrors"
	"github.com/dropoutlabs/privacyai/primitives"
)

type pgConn interface {
	Query(ctx context.Context, q string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, q string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, q string, args ...interface{}) (pgconn.CommandTag, error)
	Close()
}

type postgresQuerier struct {
	conn pgConn
}

// Create an entity inside the database
func (q *postgresQuerier) Create(ctx context.Context, e primitives.Entity) error {
	sql := fmt.Sprintf(`INSERT INTO %s (data) VALUES ($1)`, e.GetType().String())

	_, err := q.conn.Exec(ctx, sql, e)
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
func (q *postgresQuerier) Get(ctx context.Context, id primitives.ID, e primitives.Entity) error {
	t, err := id.Type()
	if err != nil {
		return err
	}

	if t != e.GetType() {
		return errors.New(TypeMismatchCause, "id and entity do not match: %s - %s",
			t.String(), e.GetType().String())
	}

	sql := fmt.Sprintf(`SELECT data FROM %s WHERE id = $1 LIMIT 1`, t.String())
	r := q.conn.QueryRow(ctx, sql, id.String())

	switch err := r.Scan(e); err {
	case pgx.ErrNoRows:
		return errors.New(NotFoundCause, "could not find entity: %s", id)
	default:
		return err
	}
}

// Delete an entity from the database
func (q *postgresQuerier) Delete(ctx context.Context, id primitives.ID) error {
	t, err := id.Type()
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, t.String())

	ct, err := q.conn.Exec(ctx, sql, id.String())
	if err != nil {
		return err
	}

	if ct.RowsAffected() != 1 {
		return errors.New(NotFoundCause, "could not find entity: %s", id)
	}

	return nil
}

// Update an entity inside the database
func (q *postgresQuerier) Update(ctx context.Context, e primitives.Entity) error {
	t := e.GetType()
	if !t.Mutable() {
		return errors.New(NotMutableCause, "cannot update an immutable entity: %s", t.String())
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
		return errors.New(NotFoundCause, "could not find entity: %s", e.GetID())
	}

	return nil
}

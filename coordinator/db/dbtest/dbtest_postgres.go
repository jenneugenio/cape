package dbtest

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	pgxpool "github.com/jackc/pgx/v4/pgxpool"

	errors "github.com/capeprivacy/cape/partyerrors"
)

// ConnectionError occurs when something can't be done because the database is
// not connected
var ConnectionError = errors.NewCause(errors.BadRequestCategory, "connection_error")

// TestPostgres implements the TestDatabase providing functionality for setting
// up a test environment for integration testing.
type TestPostgres struct {
	rootURL *url.URL
	dbURL   *url.URL // database this instance manages
	dbName  string
	db      *pgxpool.Pool
}

// NewTestPostgres returns an instance of a TestPostgres struct
func NewTestPostgres(rootURL string) (TestDatabase, error) {
	r, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}

	dbName, err := GenerateName()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(r.String())
	if err != nil {
		return nil, err
	}
	u.Path = dbName

	return &TestPostgres{
		rootURL: r,
		dbURL:   u,
		dbName:  dbName,
	}, nil
}

// Setup creates a randomly named database for testing
func (t *TestPostgres) Setup(ctx context.Context) error {
	db, err := pgxpool.Connect(ctx, t.rootURL.String())
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", t.dbName))
	if err != nil {
		return err
	}

	// okay now we can create our long lasting connection
	db, err = pgxpool.Connect(ctx, t.dbURL.String())
	if err != nil {
		return err
	}

	t.db = db

	return nil
}

// Teardown destroys the test database and closes any connection to the
// database.
func (t *TestPostgres) Teardown(ctx context.Context) error {
	if t.db == nil {
		return nil
	}

	t.db.Close()
	t.db = nil

	db, err := pgxpool.Connect(ctx, t.rootURL.String())
	if err != nil {
		return err
	}
	defer db.Close() //nolint: errcheck

	_, err = db.Exec(ctx, "DROP DATABASE "+t.dbName)
	if err != nil {
		return err
	}

	return nil
}

// Truncate truncates all tables inside the database resetting them back to
// empty while retaining any triggers or tables.
func (t *TestPostgres) Truncate(ctx context.Context) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint: errcheck

	// TODO: Check that the migrations table is actually called 'migrations'
	rows, err := tx.Query(ctx, `
		SELECT
			tablename
		FROM
			pg_catalog.pg_tables
		WHERE
			schemaname = 'public'
		AND tablename != 'migrations'
	`) // XXX: Should we create a schema for our tables?
	if err != nil {
		return err
	}
	defer rows.Close() //nolint: errcheck

	tables := []string{}
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return err
		}

		tables = append(tables, tableName)
	}

	for _, name := range tables {
		_, err = tx.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %v CASCADE", name))
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// Query exposes a method for querying the test database directly
func (t *TestPostgres) Query(ctx context.Context, q string, args ...interface{}) (pgx.Rows, error) {
	if t.db == nil {
		return nil, errors.New(ConnectionError, "must setup the database to query")
	}

	return t.db.Query(ctx, q, args...)
}

// Exec exposes a method for executing a query against the database
func (t *TestPostgres) Exec(ctx context.Context, q string, args ...interface{}) (pgconn.CommandTag, error) {
	if t.db == nil {
		return nil, errors.New(ConnectionError, "must setup the database to exec a statement")
	}

	return t.db.Exec(ctx, q, args...)
}

// URL returns the connection string for the underlying test database.
func (t *TestPostgres) URL() *url.URL {
	var n *url.URL = &url.URL{}
	*n = *(t.dbURL)

	return n
}

package dbtest

import (
	"context"
	"database/sql"
	"net/url"

	_ "github.com/lib/pq"

	errors "github.com/dropoutlabs/privacyai/partyerrors"
)

// ConnectionError occurs when something can't be done because the database is
// not connected
var ConnectionError = errors.NewCause(errors.BadRequestCategory, "connection_error")

// TestPostgres implements the TestDatabase providing functionality for setting
// up a test enviornment for integration testing.
type TestPostgres struct {
	rootURL *url.URL // the root database (e.g. template1)
	dbURL   *url.URL // database this instance manages
	dbName  string
	db      *sql.DB
}

// NewTestPostgres returns an instance of a TestPostgres struct
func NewTestPostgres(dbURL string) (TestDatabase, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, err
	}

	// TODO: Generate a random name for the test database
	return &TestPostgres{
		dbURL:   u,
		rootURL: u,
		dbName:  "test",
	}, nil
}

// Setup creates a database and migrates it to the appropriate state.
func (t *TestPostgres) Setup(ctx context.Context) error {
	db, err := sql.Open("postgres", t.dbURL.String())
	if err != nil {
		return err
	}

	// TODO: Create a database we can use for running tests :)
	t.db = db
	return nil
}

// Teardown destroys the test database and closes any connection to the
// database.
func (t *TestPostgres) Teardown(ctx context.Context) error {
	// TODO: Delete the database we created!
	return t.db.Close()
}

// Truncate truncates all tables inside the database resetting them back to
// empty while retaining any triggers or tables.
func (t *TestPostgres) Truncate(ctx context.Context) error {
	// TODO: Make sure not to truncate the migrations table
	return nil
}

// RawQuery exposes a method for directly querying the postgres backend for
// testing purposes.
func (t *TestPostgres) RawQuery(ctx context.Context, q string, args ...interface{}) (sql.Result, error) {
	if t.db == nil {
		return nil, errors.New(ConnectionError, "must setup the database to issue a query")
	}

	return t.db.ExecContext(ctx, q)
}

// URL returns the connection string for the underlying test database.
func (t *TestPostgres) URL() string {
	return t.dbURL.String()
}

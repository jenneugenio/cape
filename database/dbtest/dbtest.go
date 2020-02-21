// dbtest contains functionality for writing tests
package dbtest

import (
	"context"
	"net/url"

	errors "github.com/dropoutlabs/privacyai/partyerrors"
)

// TestDatabase represents a test sotrage backend (e.g. Postgres, MySQL,
// SQLite, etc). Implementations are responsible for providing functionality
// for setting up and tearing down test environments for integration testing.
type TestDatabase interface {
	Setup(context.Context) error
	Teardown(context.Context) error
	Truncate(context.Context) error
	URL() string
}

// NewTestDatabaseFunc represents a constructor of a TestDatabase
type NewTestDatabaseFunc func(string) (TestDatabase, error)

var validDBs = map[string]NewTestDatabaseFunc{
	"postgres": NewTestPostgres,
}

// Registry for storing global test databases that are shared across tests
var existingDBs = map[string]TestDatabase{}

// New returns an instance of a testing database for use in integration tests.
func New(dbURL string) (TestDatabase, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, err
	}

	ctor, ok := validDBs[u.Scheme]
	if !ok {
		return nil, errors.New(errors.UnsupportedErrorCause, "test database type not supported")
	}

	if testDB, ok := existingDBs[u.String()]; ok {
		return testDB, nil
	}

	testDB, err := ctor(dbURL)
	if err != nil {
		return nil, err
	}

	testDB = Wrap(testDB) // Wrap it up!
	existingDBs[u.String()] = testDB
	return testDB, nil
}

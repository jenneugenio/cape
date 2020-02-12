package database

import (
	"context"
	"net/url"

	errors "github.com/dropoutlabs/privacyai/partyerrors"
)

// Backend represents a storage backend (e.g. Postgres, MySQL, etc).
// See a concrete implementation of this interface (e.g. PostgresBackend) for
// more details.
type Backend interface {
	Open(context.Context) error
	Close() error
	Transaction() (*Transaction, error)
}

var validDbs = map[string]func(*url.URL) Backend{
	"postgres": NewPostgresBackend,
}

// New returns a new backend.
func New(dbURL *url.URL) (Backend, error) {
	ctor, ok := validDbs[dbURL.Scheme]
	if !ok {
		return nil, errors.New(NotImplementedDBCause, "database not supported")
	}

	return ctor(dbURL), nil
}

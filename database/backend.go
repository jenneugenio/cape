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
	Querier
	Open(context.Context) error
	Close() error
	Transaction(context.Context) (Transaction, error)
}

// NewBackendFunc represents a constructor of a Backend implementation
type NewBackendFunc func(*url.URL, string) (Backend, error)

var validDBs = map[string]NewBackendFunc{
	"postgres": NewPostgresBackend,
}

// New returns a new backend for the given application name
func New(dbURL *url.URL, appName string) (Backend, error) {
	ctor, ok := validDBs[dbURL.Scheme]
	if !ok {
		return nil, errors.New(NotImplementedDBCause, "database not supported")
	}

	return ctor(dbURL, appName)
}

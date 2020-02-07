package database

import (
	"context"
	"net/url"
)

// Backend represents a storage backend (e.g. Postgres, MySQL, etc).
// See a concrete implementation of this interface (e.g. PostgresBackend) for
// more details.
type Backend interface {
	Open(context.Context) error
	Close() error
	Transaction() (*Transaction, error)
}

var validDbs = map[string]func(string) Backend{
	"postgres": NewPostgresBackend,
}

// New returns a new backend.
func New(dbURL string) (Backend, error) {
	url, err := url.Parse(dbURL)
	if err != nil {
		return nil, err
	}

	ctor, ok := validDbs[url.Scheme]
	if !ok {
		return nil, newUnsupportedBackendError(url.Scheme)
	}

	return ctor(url.Host), nil
}

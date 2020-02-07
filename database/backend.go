package database

import (
	"context"
	"os"
)

// Backend represents a storage backend (e.g. Postgres, MySQL, etc).
// See a concrete implementation of this interface (e.g. PostgresBackend) for
// more details.
type Backend interface {
	Open(context.Context, *Config) error
	Close() error
	Transaction() (*Transaction, error)
}

var validDbs = map[string]func(string) Backend{
	"postgres": NewPostgresBackend,
}

// NewBackend returns a new backend.
func NewBackend() (Backend, error) {
	backend := os.Getenv("DB_BACKEND")
	if len(backend) == 0 {
		return nil, newUnspecifiedBackendError()
	}

	ctor, ok := validDbs[backend]
	if !ok {
		return nil, newUnsupportedBackendError(backend)
	}

	addr := os.Getenv("DB_ADDR")
	return ctor(addr), nil
}

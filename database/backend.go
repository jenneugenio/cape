package database

import (
	"context"
)

// Backend represents a storage backend (e.g. Postgres, MySQL, etc).
// See a concrete implementation of this interface (e.g. PostgresBackend) for
// more details.
type Backend interface {
	Open(context.Context, *Config) error
	Close() error
	Transaction() (*Transaction, error)
}

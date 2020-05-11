package database

import (
	"context"
	"net/url"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Backend represents a storage backend (e.g. Postgres, MySQL, etc).
// See a concrete implementation of this interface (e.g. PostgresBackend) for
// more details.
type Backend interface {
	Querier
	Open(context.Context) error
	Close() error
	Transaction(context.Context) (Transaction, error)
	URL() *url.URL
}

// NewBackendFunc represents a constructor of a Backend implementation
type NewBackendFunc func(crypto.EncryptionCodec, *url.URL, string) (Backend, error)

var validDBs = map[string]NewBackendFunc{
	"postgres": NewPostgresBackend,
}

// New returns a new backend for the given application name
func New(codec crypto.EncryptionCodec, dbURL *url.URL, appName string) (Backend, error) {
	ctor, ok := validDBs[dbURL.Scheme]
	if !ok {
		return nil, errors.New(NotImplementedCause, "database not supported")
	}

	return ctor(codec, dbURL, appName)
}

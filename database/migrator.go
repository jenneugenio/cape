package database

import (
	"context"
	errors "github.com/dropoutlabs/privacyai/partyerrors"
	"net/url"
)

// Migrator represents backend/database migrations
// See a concrete implementation of this interface (e.g. PostgresMigrator) for more details.
type Migrator interface {
	Up (context.Context) error
	Down (context.Context) error
}

// NewMigratorFunc represents a constructor of a Migrator implementation
type NewMigratorFunc func(dbURL *url.URL, migrations ...string) (Migrator, error)

var validMigrators = map[string]NewMigratorFunc{
	"postgres": NewPostgresMigrator,
}

// NewMigrator returns a new migrator for the given db type (based off of the provided URL)
func NewMigrator (dbURL *url.URL, migrations ...string) (Migrator, error) {
	ctor, ok := validMigrators[dbURL.Scheme]
	if !ok {
		return nil, errors.New(NotImplementedCause, "migrator not supported")
	}

	return ctor(dbURL, migrations...)
}


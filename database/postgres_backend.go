package database

import (
	"context"
	"database/sql"
	"net/url"

	// postgres driver
	_ "github.com/lib/pq"
)

// PostgresBackend implements the backend interface for a pg database
type PostgresBackend struct {
	dbURL *url.URL
	db    *sql.DB
}

// Open the database
func (p *PostgresBackend) Open(ctx context.Context) error {
	// XXX: We should look into the pgx driver
	db, err := sql.Open("postgres", p.dbURL.String())
	if err != nil {
		return err
	}

	p.db = db

	return nil
}

// Close the database
func (p *PostgresBackend) Close() error {
	return nil
}

// Transaction starts a new transaction
func (p *PostgresBackend) Transaction() (*Transaction, error) {
	return nil, nil
}

// NewPostgresBackend returns a new postgres backend instance
func NewPostgresBackend(dbURL *url.URL) Backend {
	return &PostgresBackend{
		dbURL: dbURL,
	}
}

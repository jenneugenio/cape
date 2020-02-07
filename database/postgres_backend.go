package database

import "context"

// PostgresBackend implements the backend interface for a pg database
type PostgresBackend struct {
	connString string
}

// Open the database
func (p *PostgresBackend) Open(ctx context.Context) error {
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
func NewPostgresBackend(connString string) Backend {
	return &PostgresBackend{
		connString: connString,
	}
}

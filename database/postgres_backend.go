package database

import (
	"context"
	"net/url"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresBackend implements the backend interface for a pg database
type PostgresBackend struct {
	*postgresQuerier
	dbURL *url.URL
	cfg   *pgxpool.Config
}

// Open the database
func (p *PostgresBackend) Open(ctx context.Context) error {
	c, err := pgxpool.ConnectConfig(ctx, p.cfg)
	if err != nil {
		return err
	}

	p.conn = c // inherited from postgresQuerier

	return nil
}

// Close the database
func (p *PostgresBackend) Close() error {
	p.conn.Close()
	p.conn = nil

	return nil
}

// Transaction starts a new transaction
func (p *PostgresBackend) Transaction() (*Transaction, error) {
	return nil, nil
}

// NewPostgresBackend returns a new postgres backend instance
func NewPostgresBackend(dbURL *url.URL, name string) (Backend, error) {
	cfg, err := pgxpool.ParseConfig(dbURL.String())
	if err != nil {
		return nil, err
	}

	// Set the application name which can be used for identifying which service
	// is connecting to postgres
	cfg.ConnConfig.RuntimeParams = map[string]string{
		"application_name": name,
	}

	return &PostgresBackend{
		postgresQuerier: &postgresQuerier{},
		dbURL:           dbURL,
		cfg:             cfg,
	}, nil
}

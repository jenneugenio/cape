package database

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"net/url"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
)

// PostgresBackend implements the backend interface for a pg database
type PostgresBackend struct {
	*postgresQuerier
	dbURL *url.URL
	cfg   *pgxpool.Config
	pool  *pgxpool.Pool
}

// Open the database
func (p *PostgresBackend) Open(ctx context.Context) error {
	c, err := pgxpool.ConnectConfig(ctx, p.cfg)
	if err != nil {
		return err
	}

	// We need to separate out the `conn` from the `pool` as the
	// `postgresQuerier` and `Querier` interface do not implement any
	// transaction related methods. This is because Querier is a common
	// interface over both a pgconn.Conn and pgxpool.Pool
	p.conn = c // inherited from postgresQuerier
	p.pool = c

	return nil
}

// Close the database
func (p *PostgresBackend) Close() error {
	if p.pool == nil && p.conn == nil {
		return nil
	}

	p.pool.Close()
	p.conn = nil
	p.pool = nil

	return nil
}

// Transaction starts a new transaction.
//
// This method returns a dedication connection which any users of must manage
// themselves. Please ensure to call Rollback() or Commit() so the connection
// is returned to the pool once you're done.
func (p *PostgresBackend) Transaction(ctx context.Context) (Transaction, error) {
	pgtx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// We need to set both the `conn` and `tx` as `conn` satisfies the Querier
	// interface while `tx` is a straight pgx.Tx type
	return &PostgresTransaction{
		postgresQuerier: &postgresQuerier{
			conn:  pgtx,
			codec: p.codec,
		},
		tx: pgtx,
	}, nil
}

// URL returns the underlying database URL
func (p *PostgresBackend) URL() *url.URL {
	var c *url.URL = &url.URL{}
	*c = *(p.dbURL)

	return c
}

func (p *PostgresBackend) SetEncryptionCodec(codec crypto.EncryptionCodec) {
	p.codec = codec
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

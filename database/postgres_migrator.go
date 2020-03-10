package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"net/url"
)

// PostgresMigrator implements the Migrator interface for a postgres backend.
type PostgresMigrator struct {
	dbURL      *url.URL
	migrations []string
}

func (p *PostgresMigrator) getMigrator(ctx context.Context, conn *pgx.Conn) (*migrate.Migrator, error) {
	m, err := migrate.NewMigrator(ctx, conn, "migrations")
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(p.migrations); i++ {
		err = m.LoadMigrations(p.migrations[i])
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

// Up runs all of the provided migrations
func (p *PostgresMigrator) Up(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, p.dbURL.String())
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	m, err := p.getMigrator(ctx, conn)
	if err != nil {
		return err
	}

	err = m.Migrate(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Down rolls the db back to initial state (ie undoes all of the migrations)
func (p *PostgresMigrator) Down(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, p.dbURL.String())
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	m, err := p.getMigrator(ctx, conn)
	if err != nil {
		return err
	}

	return m.MigrateTo(ctx, 0)
}

// NewPostgresMigrator returns a postgres migrator that adheres to the Migrator interface
func NewPostgresMigrator(dbURL *url.URL, migrations ...string) (Migrator, error) {
	return &PostgresMigrator{dbURL: dbURL, migrations: migrations}, nil
}

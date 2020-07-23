package capepg

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/jackc/pgx/v4"
)

type pgConfig struct {
	pool    Pool
	timeout time.Duration
}

var _ db.ConfigDB = &pgConfig{}

func (c *pgConfig) Create(ctx context.Context, config models.Config) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.Get(ctx)
	if err == nil {
		return fmt.Errorf("error config already exists")
	}

	s, args, err := sq.Insert("config").
		PlaceholderFormat(sq.Dollar).
		Columns("data").
		Values(config).
		ToSql()

	if err != nil {
		return fmt.Errorf("error generating query: %w", err)
	}

	_, err = c.pool.Exec(ctx, s, args...)
	if err != nil {
		return fmt.Errorf("error creating config: %w", err)
	}

	return nil
}

func (c *pgConfig) Get(ctx context.Context) (*models.Config, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("config").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating query: %w", err)
	}

	row := c.pool.QueryRow(ctx, s, args...)

	var config models.Config
	err = row.Scan(&config)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, db.ErrNoRows
		}
		return nil, fmt.Errorf("error getting config: %w", err)
	}

	return &config, nil
}

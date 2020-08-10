package capepg

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
)

type pgSecret struct {
	pool    Pool
	timeout time.Duration
}

var _ db.SecretDB = &pgSecret{}

func (p *pgSecret) Create(ctx context.Context, secret models.SecretArg) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Insert("secrets").
		PlaceholderFormat(sq.Dollar).
		Columns("data").
		Values(secret).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building secret create query: %w", err)
	}

	_, err = p.pool.Exec(ctx, s, args...)
	if err != nil {
		return fmt.Errorf("error creating secret: %w", err)
	}

	return nil
}

func (p *pgSecret) Delete(ctx context.Context, e string) (db.DeleteStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Delete("secrets").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"data->>'email'": e}).
		ToSql()
	if err != nil {
		return db.DeleteStatusError, fmt.Errorf("error building secret update query: %w", err)
	}

	tag, err := p.pool.Exec(ctx, s, args...)
	if err != nil {
		return db.DeleteStatusError, fmt.Errorf("error deleting secret: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return db.DeleteStatusDoesNotExist, nil
	}

	return db.DeleteStatusDeleted, nil
}

func (p *pgSecret) Get(ctx context.Context, name string) (*models.SecretArg, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("secrets").
		Where(sq.Eq{"data->>'name'": name}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := p.pool.QueryRow(ctx, s, args...)
	secret := &models.SecretArg{}
	err = row.Scan(secret)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, db.ErrCannotFindSecret
		}
		return nil, fmt.Errorf("error retrieving secret: %w", err)
	}

	return secret, nil
}

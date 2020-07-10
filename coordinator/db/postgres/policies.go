package capepg

import (
	"context"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
)

type pgPolicy struct {
	pool    Pool
	timeout time.Duration
}

var _ db.PolicyDB = &pgPolicy{}

func (p *pgPolicy) Create(ctx context.Context, policy models.Policy) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Insert("policies").
		PlaceholderFormat(sq.Dollar).
		Columns("data").
		Values(policy).
		ToSql()

	if err != nil {
		return fmt.Errorf("error generating query: %w", err)
	}

	_, err = p.pool.Exec(ctx, s, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return db.ErrDuplicateKey
		}
		return fmt.Errorf("error creating policy: %w", err)
	}

	return nil
}

func (p *pgPolicy) Delete(ctx context.Context, l models.Label) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Delete("policies").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"data->>'label'": l}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = p.pool.Exec(ctx, s, args...)

	if err != nil {
		return fmt.Errorf("error deleting policy: %w", err)
	}
	return nil
}

func (p *pgPolicy) Get(ctx context.Context, l models.Label) (*models.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("policies").
		Where(sq.Eq{"data->>'label'": string(l)}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating query: %w", err)
	}

	row := p.pool.QueryRow(ctx, s, args...)

	var policy models.Policy
	err = row.Scan(&policy)
	if err != nil {
		return nil, fmt.Errorf("error getting policy: %w", err)
	}

	return &policy, nil
}

func (p *pgPolicy) List(ctx context.Context, opts *db.ListPolicyOptions) ([]models.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("policies").
		OrderBy("data->>'created_at'")

	if opts != nil {
		query = query.Limit(opts.Limit).Offset(opts.Offset)
	}

	s, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, s, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing policies: %w", err)
	}
	defer rows.Close()

	var policies []models.Policy
	for rows.Next() {
		var policy models.Policy
		if err := rows.Scan(&policy); err != nil {
			return nil, fmt.Errorf("TODO: be more graceful when a policy errors like %w", err)
		}
		policies = append(policies, policy)
	}

	return policies, rows.Err()
}

package capepg

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"time"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pgPolicy struct {
	pool    *pgxpool.Pool
	timeout time.Duration
}

var _ db.PolicyDB = &pgPolicy{}

func (p *pgPolicy) Create(ctx context.Context, policy *models.Policy) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Insert("policies").
		PlaceholderFormat(sq.Dollar).
		Columns("data").
		Values(policy).
		ToSql()

	_, err = p.pool.Exec(ctx, s, args...)
	if err != nil {
		return fmt.Errorf("error creating policy: %w", err)
	}

	return nil
}

func (p *pgPolicy) Delete(ctx context.Context, l models.Label) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Delete("policies").
		Where(sq.Eq{"data->>'label'": l}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = p.pool.Exec(ctx, s, args)

	if err != nil {
		return fmt.Errorf("error deleting policy: %w", err)
	}
	return nil
}

func (p *pgPolicy) Get(ctx context.Context, l models.Label) (*models.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		From("policies").
		Where(sq.Eq{"data->>'label'": l}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := p.pool.QueryRow(ctx, s, args)
	var policyBlob []byte
	err = row.Scan(policyBlob)
	if err != nil {
		return nil, err
	}

	var policy models.Policy
	err = json.Unmarshal(policyBlob, &policy)

	return &policy, err
}

func (p *pgPolicy) List(ctx context.Context, opts *db.ListPolicyOptions) ([]*models.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		From("policies").
		OrderBy("data->>'created_at'").
		Limit(opts.Limit).
		Offset(opts.Offset).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, s, args)
	defer rows.Close()

	var policies []*models.Policy
	for rows.Next() {
		var policyString []byte
		if err := rows.Scan(&policyString); err != nil {
			return nil, fmt.Errorf("TODO: be more graceful when a policy errors like %w", err)
		}
		policy, err := models.ParsePolicy([]byte(policyString))
		if err != nil {
			return nil, fmt.Errorf("malformed policy encountered: %w", err)
		}
		policies = append(policies, policy)
	}

	return policies, nil
}

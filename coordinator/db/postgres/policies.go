package capepg

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"time"
)

type pgPolicy struct {
	db      *sql.DB
	timeout time.Duration
}

func (p *pgPolicy) Create(ctx context.Context, policy *models.Policy) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	policyBlob, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	s, args, err := sq.Insert("policies").
		Columns("data").
		Values(policyBlob).
		ToSql()

	_, err = p.db.ExecContext(ctx, s, args)
	if err != nil {
		return fmt.Errorf("error creating policy: %w", err)
	}
	return nil
}

func (p *pgPolicy) Delete(ctx context.Context, label models.Label) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Delete("policies").
		Where(sq.Eq{"data->>'label'": label}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, s, args)
	if err != nil {
		return fmt.Errorf("error deleting policy: %w", err)
	}
	return nil
}

func (p *pgPolicy) Get(ctx context.Context, label models.Label) (*models.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		From("policies").
		Where(sq.Eq{"data->>'label'": label}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := p.db.QueryRowContext(ctx, s, args)
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

	rows, err := p.db.QueryContext(ctx,s, args)
	if err != nil {
		return nil, fmt.Errorf("error retrieving policies: %w", err)
	}
	defer rows.Close()

	var policies []*models.Policy
	for rows.Next() {
		var policyString string
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
package capepg

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type pgPolicy struct {
	db      *sql.DB
	timeout time.Duration
}

var _ db.PolicyDB = &pgPolicy{}

func (p *pgPolicyDB) Create(ctx context.Context, policy models.Policy) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	err := p.db.Exec("INSERT INTO policies (label, spec) VALUES ($1, $2)", policy.Label, policy.String())
	if err != nil {
		return fmt.Errorf("error creating policy: %w", err)
	}
	return nil
}

func (p *pgPolicyDB) Delete(ctx context.Context, l models.Label) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	err := p.db.ExecContext("DELETE FROM policies WHERE label=$1", l)
	if err != nil {
		return fmt.Errorf("error deleting policy: %w", err)
	}
	return nil
}

func (p *pgPolicyDB) Get(context.Context, models.Label) (models.Policy, error) {
}

func (p *pgPolicyDB) List(ctx context.Context, opt ListPolicyOptions) ([]models.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	rows, err := p.db.QueryContext("SELECT spec FROM policies ORDER BY created_at LIMIT $1 OFFSET $2", opt.Limit, opt.Offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving policies: %w", err)
	}
	defer rows.Close()

	var p []models.Policy
	for rows.Next() {
		var policyString string
		if err := rows.Scan(&policyString); err != nil {
			return fmt.Errorf("TODO: be more graceful when a policy errors like %w", err)
		}
		policy, err := models.ParsePolicy(policyString)
		if err != nil {
			return fmt.Errorf("malformed policy encountered: %w", err)
		}
		p = append(p, policy)
	}

	return p, nil
}

type ListPolicyOptions struct {
	Limit  int
	Offset int
}

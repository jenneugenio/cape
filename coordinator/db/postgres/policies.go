package capepg

import (
	"context"
	"fmt"
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

func (p *pgPolicy) Create(ctx context.Context, policy models.Policy) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := p.pool.Exec(ctx, "INSERT INTO policies (data) VALUES ($1)", policy)
	if err != nil {
		return fmt.Errorf("error creating policy: %w", err)
	}

	return nil
}

func (p *pgPolicy) Delete(ctx context.Context, l models.Label) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := p.pool.Exec(ctx, "DELETE FROM policies WHERE label=$1", l)
	if err != nil {
		return fmt.Errorf("error deleting policy: %w", err)
	}
	return nil
}

func (p *pgPolicy) Get(ctx context.Context, l models.Label) (*models.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_ = p.pool.QueryRow(ctx, "SELECT data FROM policies WHERE label=$1", l)

	// This needs to be populated from the returned row
	var policy *models.Policy

	return policy, nil
}

func (p *pgPolicy) List(ctx context.Context, opt db.ListPolicyOptions) ([]models.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	rows, err := p.pool.Query(ctx, "SELECT data FROM policies ORDER BY created_at LIMIT $1 OFFSET $2", opt.Limit, opt.Offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving policies: %w", err)
	}
	defer rows.Close()

	var policies []models.Policy
	for rows.Next() {
		var policyString []byte
		if err := rows.Scan(&policyString); err != nil {
			return nil, fmt.Errorf("TODO: be more graceful when a policy errors like %w", err)
		}
		policy, err := models.ParsePolicy(policyString)
		if err != nil {
			return nil, fmt.Errorf("malformed policy encountered: %w", err)
		}
		policies = append(policies, *policy)
	}

	return policies, nil
}

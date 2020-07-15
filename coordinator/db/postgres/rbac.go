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

type pgRBAC struct {
	pool    Pool
	timeout time.Duration
}

var _ db.RBACDB = &pgRBAC{}

func (p *pgRBAC) Create(ctx context.Context, rbac models.RBACPolicy) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Insert("policies").
		PlaceholderFormat(sq.Dollar).
		Columns("data", "type").
		Values(rbac, PolicyTypeRBAC).
		ToSql()

	if err != nil {
		return fmt.Errorf("error generating query: %w", err)
	}

	_, err = p.pool.Exec(ctx, s, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return db.ErrDuplicateKey
		}
		return fmt.Errorf("error creating rbac: %w", err)
	}

	return nil
}

func (p *pgRBAC) List(ctx context.Context, opts *db.ListRBACOptions) ([]models.RBACPolicy, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("policies").
		OrderBy("data->>'created_at'").
		Where(sq.Eq{"type": PolicyTypeRBAC})

	if opts != nil {
		if len(opts.FilterIDs) > 0 {
			query = query.Where(sq.Eq{"id": opts.FilterIDs, "type": PolicyTypeRBAC})
		}
	}

	s, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, s, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing rbacs: %w", err)
	}
	defer rows.Close()

	var rbacs []models.RBACPolicy
	for rows.Next() {
		var rbac models.RBACPolicy
		if err := rows.Scan(&rbac); err != nil {
			return nil, fmt.Errorf("TODO: be more graceful when a rbac policy errors like %w", err)
		}

		rbacs = append(rbacs, rbac)
	}

	return rbacs, rows.Err()
}

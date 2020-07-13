package capepg

import (
	"bytes"
	"context"
	"encoding/json"
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

func (p *pgRBAC) Create(ctx context.Context, rbac models.RBAC) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Insert("policies").
		PlaceholderFormat(sq.Dollar).
		Columns("data").
		Values(rbac).
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

func (p *pgRBAC) List(ctx context.Context, opts *db.ListRBACOptions) ([]models.RBAC, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("policies").
		OrderBy("data->>'created_at'")

	if opts != nil {
		if opts.FilterIDs != nil && len(opts.FilterIDs) > 0 {
			query = query.Where(sq.Eq{"id": opts.FilterIDs})
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

	var rbacs []models.RBAC
	for rows.Next() {
		var rbacBytes []byte
		if err := rows.Scan(&rbacBytes); err != nil {
			continue
		}

		dec := json.NewDecoder(bytes.NewBuffer(rbacBytes))
		dec.DisallowUnknownFields()

		var rbac models.RBAC
		err := dec.Decode(&rbac)
		if err != nil {
			continue
		}

		rbacs = append(rbacs, rbac)
	}

	return rbacs, rows.Err()
}

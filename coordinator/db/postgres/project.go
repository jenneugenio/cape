package capepg

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/jackc/pgx/v4"
	"time"
)

type pgProject struct {
	pool    Pool
	timeout time.Duration
}

func (p *pgProject) GetByID(ctx context.Context, ID string) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("projects").
		Where(sq.Eq{"id": ID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := p.pool.QueryRow(ctx, s, args...)
	project := &models.Project{}
	err = row.Scan(project)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, db.ErrNoRows
		}
		return nil, err
	}

	return project, nil
}

func (p *pgProject) Get(ctx context.Context, label models.Label) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("projects").
		Where(sq.Eq{"data->>'label'": label}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := p.pool.QueryRow(ctx, s, args...)
	project := &models.Project{}
	err = row.Scan(project)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, db.ErrNoRows
		}
		return nil, err
	}

	return project, nil
}

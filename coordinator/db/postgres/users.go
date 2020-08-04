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

type pgUser struct {
	pool    Pool
	timeout time.Duration
}

var _ db.UserDB = &pgUser{}

func (p *pgUser) Create(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("data").
		Values(user).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building user create query: %w", err)
	}

	_, err = p.pool.Exec(ctx, s, args...)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (p *pgUser) Update(ctx context.Context, id string, user models.User) error {
	s, args, err := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("data", user).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building user update query: %w", err)
	}

	_, err = p.pool.Exec(ctx, s, args...)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (p *pgUser) Delete(ctx context.Context, e models.Email) (db.DeleteStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"data->>'email'": e}).
		ToSql()
	if err != nil {
		return db.DeleteStatusError, fmt.Errorf("error building user update query: %w", err)
	}

	tag, err := p.pool.Exec(ctx, s, args...)
	if err != nil {
		return db.DeleteStatusError, fmt.Errorf("error deleting user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return db.DeleteStatusDoesNotExist, nil
	}

	return db.DeleteStatusDeleted, nil
}

func (p *pgUser) Get(ctx context.Context, e models.Email) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"data->>'email'": e.String()}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := p.pool.QueryRow(ctx, s, args...)
	user := &models.User{}
	err = row.Scan(user)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, db.ErrCannotFindUser
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return user, nil
}

func (p *pgUser) GetByID(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := p.pool.QueryRow(ctx, s, args...)
	user := &models.User{}
	err = row.Scan(user)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, db.ErrCannotFindUser
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return user, nil
}

func (p *pgUser) List(ctx context.Context, opts *db.ListUserOptions) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("users").
		OrderBy("data->>'created_at'")

	if opts != nil {
		if opts.Options != nil {
			query = query.Limit(opts.Options.Limit).Offset(opts.Options.Offset)
		}

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
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user); err != nil {
			return nil, fmt.Errorf("TODO: be more graceful when a user errors like %w", err)
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

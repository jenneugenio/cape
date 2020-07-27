package capepg

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"time"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
)

type pgRole struct {
	pool    Pool
	timeout time.Duration
}

var _ db.RoleDB = &pgRole{}

func (r *pgRole) Create(context.Context, *models.Role) error {
	return errors.New("not implemented")
}

func (r *pgRole) Delete(context.Context, models.Label) (db.DeleteStatus, error) {
	return db.DeleteStatusError, errors.New("not implemented")
}

func (r *pgRole) Get(ctx context.Context, label models.Label) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("roles").
		Where(sq.Eq{"data->>'label'": label}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := r.pool.QueryRow(ctx, s, args...)
	role := &models.Role{}
	err = row.Scan(role)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, db.ErrNoRows
		}
		return nil, err
	}

	return role, nil
}

func (r *pgRole) GetByID(ctx context.Context, ID string) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	s, args, err := sq.Select("data").
		PlaceholderFormat(sq.Dollar).
		From("roles").
		Where(sq.Eq{"id": ID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := r.pool.QueryRow(ctx, s, args...)
	role := &models.Role{}
	err = row.Scan(role)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, db.ErrNoRows
		}
		return nil, err
	}

	return role, nil
}

func (r *pgRole) List(context.Context, *db.ListRoleOptions) ([]*models.Role, error) {
	return nil, errors.New("not implemented")
}

func (r *pgRole) AttachPolicy(context.Context, models.Label) error {
	return errors.New("not implemented")
}

func (r *pgRole) DetachPolicy(context.Context, models.Label) error {
	return errors.New("not implemented")
}

func (r *pgRole) GetByUserID(ctx context.Context, userId string) ([]models.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// get all global roles

	s := `select roles.data 
			from roles, assignments 
			where roles.id = assignments.role_id and assignments.user_id = $1;`

	rows, err := r.pool.Query(ctx, s, userId)
	if err != nil {
		return nil, err
	}

	var roles []models.Role
	for rows.Next() {
		var r models.Role
		rows.Scan(&r)

		roles = append(roles, r)
	}

	return roles, nil
}


package capepg

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"strings"
	"time"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
)

type pgRole struct {
	pool    Pool
	timeout time.Duration
}

var _ db.RoleDB = &pgRole{}


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

type assignmentIDs struct {
	UserID    string `json:"user_id"`
	RoleID    string `json:"user_id"`
	ProjectID string `json:"project_id"`
}

func (r *pgRole) SetOrgRole(ctx context.Context, email models.Email, label models.Label) (*models.Assignment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var ids assignmentIDs
	IDsQuery := `select json_build_object(
		'user_id', users.data->>'id',
		'role_id', roles.data->>'id')

		from users, roles
		where 
			users.data->>'email' = $1 AND
			roles.data->>'label' = $2;`

	row := r.pool.QueryRow(ctx, IDsQuery, email, label)
	err := row.Scan(&ids)
	if err != nil {
		return nil, err
	}

	assignment := models.Assignment{
		ID:        models.NewID(),
		UserID:    ids.UserID,
		RoleID:    ids.RoleID,
		ProjectID: "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s := `insert into assignments (data) VALUES ($1)`
	_, err = r.pool.Exec(ctx, s, assignment)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, db.ErrDuplicateKey
		}

		return nil, fmt.Errorf("error creating contributor: %w", err)
	}

	return &assignment, nil
}

func (r *pgRole) SetProjectRole(ctx context.Context, email models.Email, label models.Label, label2 models.Label) (*models.Assignment, error) {
	panic("implement me")
}



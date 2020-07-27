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

// Return all of the roles (global & project) that a user belongs to
func (r *pgRole) GetAll(ctx context.Context, userID string) (*models.UserRoles, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// get the users global role
	s := `select roles.data 
			from roles, assignments 
			where roles.id = assignments.role_id and assignments.user_id = $1;`

	userRoles := models.UserRoles{}
	row := r.pool.QueryRow(ctx, s, userID)
	err := row.Scan(&userRoles.Global)
	if err != nil {
		return nil, err
	}

	// get the project roles
	s = `select json_build_object(projects.data->>'label', roles.data)
		from roles, assignments, projects
		where roles.id = assignments.role_id and assignments.user_id = $1 and projects.id = assignments.data->>'project_id';`

	row = r.pool.QueryRow(ctx, s, userID)
	err = row.Scan(&userRoles.Projects)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return &userRoles, nil
		}

		return nil, err
	}

	return &userRoles, nil
}

type assignmentIDs struct {
	UserID    string `json:"user_id"`
	RoleID    string `json:"role_id"`
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

func (r *pgRole) SetProjectRole(ctx context.Context, email models.Email, project models.Label, role models.Label) (*models.Assignment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var ids assignmentIDs
	IDsQuery := `select json_build_object(
		'user_id', users.data->>'id',
		'project_id', projects.data->>'id',
		'role_id', roles.data->>'id')

		from users, roles, projects
		where 
			users.data->>'email' = $1 AND
			projects.data->>'label' = $2 AND
			roles.data->>'label' = $3;
`

	row := r.pool.QueryRow(ctx, IDsQuery, email, project, role)
	err := row.Scan(&ids)
	if err != nil {
		return nil, err
	}

	assignment := models.Assignment{
		ID:        models.NewID(),
		UserID:    ids.UserID,
		RoleID:    ids.RoleID,
		ProjectID: ids.ProjectID,
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

func (r *pgRole) CreateSystemRoles(ctx context.Context) error {
	roles := make([]models.Role, len(models.SystemRoles))
	insert := sq.Insert("roles").
		PlaceholderFormat(sq.Dollar).
		Columns("data")

	for i, r := range models.SystemRoles {
		role := models.Role{
			ID:        models.NewID(),
			Version:   1,
			Label:     r,
			System:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		roles[i] = role

		insert = insert.Values(role)
	}

	s, args, err := insert.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, s, args...)
	return err
}

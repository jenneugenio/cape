package capepg

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

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
			return nil, db.ErrCannotFindRole
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
			return nil, db.ErrCannotFindRole
		}
		return nil, err
	}

	return role, nil
}

// List is not yet implemented
func (r *pgRole) List(context.Context, *db.ListRoleOptions) ([]*models.Role, error) {
	return nil, errors.New("not implemented")
}

// GetAll returns all of the roles (global & project) that a user belongs to
func (r *pgRole) GetAll(ctx context.Context, userID string) (*models.UserRoles, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// get the users global role
	s := `select roles.data
			from roles, assignments
			where roles.id = assignments.role_id and assignments.user_id = $1 and assignments.data->>'project_id' = '';`

	userRoles := models.UserRoles{
		Projects: models.ProjectRolesMap{},
	}
	row := r.pool.QueryRow(ctx, s, userID)
	err := row.Scan(&userRoles.Global)
	if err != nil {
		return nil, err
	}

	s = `select projects.data->>'label' as project_label, roles.data as role
		from
			roles, assignments, projects
		where
			assignments.data->>'role_id'=roles.data->>'id' and
			projects.data->>'id' = assignments.data->>'project_id' and
			assignments.data->>'user_id' = $1;`

	rows, err := r.pool.Query(ctx, s, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		resp := struct {
			Project string
			Role    models.Role
		}{}

		err = rows.Scan(&resp.Project, &resp.Role)
		if err != nil {
			return nil, err
		}

		userRoles.Projects[models.Label(resp.Project)] = resp.Role
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

	s := `insert into assignments (data) VALUES ($1)
			on conflict on constraint unique_assignment
			do update set data = assignments.data || jsonb_build_object('role_id', $2::text), role_id = $2;`

	_, err = r.pool.Exec(ctx, s, assignment, ids.RoleID) // , ids.UserID)
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

	s := `insert into assignments (data) VALUES ($1)
			on conflict on constraint unique_assignment
			do update set data = assignments.data || jsonb_build_object('role_id', $2::text), role_id = $2;`
	_, err = r.pool.Exec(ctx, s, assignment, ids.RoleID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, db.ErrDuplicateKey
		}

		return nil, fmt.Errorf("error creating contributor: %w", err)
	}

	return &assignment, nil
}

func (r *pgRole) GetOrgRole(ctx context.Context, email models.Email) (*models.Role, error) {
	s := `select roles.data from roles, assignments, users
		where roles.data->>'id' = assignments.data->>'role_id' and
		assignments.data->>'user_id' = users.data->>'id' and
		users.data->>'email' = $1 and
		assignments.data->>'project_id' = '';`

	var role models.Role
	row := r.pool.QueryRow(ctx, s, email)
	err := row.Scan(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *pgRole) GetProjectRole(ctx context.Context, email models.Email, project string) (*models.Role, error) {
	s := `select roles.data from roles, assignments, users
		where roles.data->>'id' = assignments.data->>'role_id' and
		assignments.data->>'user_id' = users.data->>'id' and
		users.data->>'email' = $1 and
		assignments.data->>'project_id' = $2;`

	var role models.Role
	row := r.pool.QueryRow(ctx, s, email, project)
	err := row.Scan(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
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

package capepg

import (
	"context"
	"fmt"
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

var _ db.ProjectsDB = &pgProject{}

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

func (p *pgProject) Create(ctx context.Context, project models.Project) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `insert into projects (data) values ($1)`
	_, err := p.pool.Exec(ctx, s, project)
	if err != nil {
		return fmt.Errorf("entity already exists")
	}
	return err
}

func (p *pgProject) Update(ctx context.Context, project models.Project) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `update projects set data = $1 where id = $2`
	_, err := p.pool.Exec(ctx, s, project, project.ID)
	return err
}

func (p *pgProject) CreateProjectSpec(ctx context.Context, spec models.Policy) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `insert into project_specs (data) values ($1)`
	_, err := p.pool.Exec(ctx, s, spec)
	return err
}

func (p *pgProject) GetProjectSpec(ctx context.Context, id string) (*models.Policy, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `select data from project_specs where id = $1`
	row := p.pool.QueryRow(ctx, s, id)
	var spec models.Policy
	err := row.Scan(&spec)

	return &spec, err
}

func (p *pgProject) CreateSuggestion(ctx context.Context, suggestion models.Suggestion) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `insert into suggestions (data) values ($1)`
	_, err := p.pool.Exec(ctx, s, suggestion)

	return err
}

func (p *pgProject) GetSuggestions(ctx context.Context, projectLabel models.Label) ([]models.Suggestion, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `select data from suggestions where project_id = (select id from projects where projects.data->>'label' = $1)`
	rows, err := p.pool.Query(ctx, s, projectLabel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []models.Suggestion
	for rows.Next() {
		var s models.Suggestion
		err = rows.Scan(&s)
		if err != nil {
			return nil, err
		}

		suggestions = append(suggestions, s)
	}

	return suggestions, err
}

func (p *pgProject) GetSuggestion(ctx context.Context, id string) (*models.Suggestion, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `select data from suggestions where id = $1`
	row := p.pool.QueryRow(ctx, s, id)
	var suggestion models.Suggestion

	err := row.Scan(&suggestion)
	if err != nil {
		return nil, err
	}

	return &suggestion, err
}

func (p *pgProject) UpdateSuggestion(ctx context.Context, suggestion models.Suggestion) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `update suggestions set data = $1 where id = $2`
	_, err := p.pool.Exec(ctx, s, suggestion, suggestion.ID)
	return err
}

func (p *pgProject) List(ctx context.Context) ([]models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `select data from projects`
	rows, err := p.pool.Query(ctx, s)
	if err != nil {
		return nil, err
	}

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(&p)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}

func (p *pgProject) ListByStatus(ctx context.Context, status models.ProjectStatus) ([]models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `select data from projects where data->>'status' = $1`
	rows, err := p.pool.Query(ctx, s, status)
	if err != nil {
		return nil, err
	}

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(&p)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}

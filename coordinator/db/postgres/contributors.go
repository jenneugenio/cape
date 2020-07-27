package capepg

import (
	"context"
	"fmt"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"strings"
	"time"
)

type pgContributor struct {
	pool    Pool
	timeout time.Duration
}

var _ db.ContributorDB = &pgContributor{}

type contributorAddIDs struct {
	UserID    string `json:"user_id"`
	RoleID    string `json:"role_id"`
	ProjectID string `json:"project_id"`
}

func (p *pgContributor) Add(ctx context.Context, project models.Label, email models.Email) (*models.Contributor, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// Get the user, project, & role id so we can build the contributor object
	var ids contributorAddIDs

	IDsQuery := `select json_build_object(
		'user_id', users.data->>'id',
		'project_id', projects.data->>'id') from users, projects where
		users.data->>'email' = $1 and
		projects.data->>'label' = $2;`

	row := p.pool.QueryRow(ctx, IDsQuery, email, project)
	err := row.Scan(&ids)
	if err != nil {
		return nil, err
	}

	contributor := models.Contributor{
		ID:        models.NewID(),
		UserID:    ids.UserID,
		ProjectID: ids.ProjectID,
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s := `INSERT INTO contributors (data) VALUES ($1)`

	_, err = p.pool.Exec(ctx, s, contributor)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, db.ErrDuplicateKey
		}

		return nil, fmt.Errorf("error creating contributor: %w", err)
	}

	return &contributor, nil
}

func (p *pgContributor) Get(ctx context.Context, project models.Label, email models.Email) (*models.Contributor, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `select data from contributors where
  		 data->>'user_id' = (select id from users where data->>'email' = $1) AND
         data->>'project_id' = (select id from projects where data->>'label' = $2);
	`

	var c models.Contributor
	row := p.pool.QueryRow(ctx, s, email, project)
	err := row.Scan(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (p *pgContributor) List(ctx context.Context, project models.Label) ([]models.Contributor, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `select data from contributors where
          data->>'project_id' = (select id from projects where data->>'label' = $1);
	`

	rows, err := p.pool.Query(ctx, s, project)
	if err != nil {
		return nil, fmt.Errorf("error creating contributor: %w", err)
	}
	defer rows.Close()

	var contributors []models.Contributor
	for rows.Next() {
		var c models.Contributor
		if err := rows.Scan(&c); err != nil {
			return nil, fmt.Errorf("error fetching contributors")
		}
		contributors = append(contributors, c)
	}

	return contributors, nil
}

func (p *pgContributor) Delete(ctx context.Context, projectLabel models.Label, userEmail models.Email) (*models.Contributor, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := `delete from contributors where 
		  data->>'user_id' = (
    	  	select id from users where data->>'email' = $1
          ) AND
          data->>'project_id' = (
			select id from projects where data->>'label' = $2
		  )
          RETURNING data; 
	`

	rows, err := p.pool.Query(ctx, s, userEmail, projectLabel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var c models.Contributor
	rows.Next()
	err = rows.Scan(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

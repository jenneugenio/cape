package capepg

import (
	"context"
	"fmt"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"strings"
	"time"
)

type pgAssignments struct {
	pool    Pool
	timeout time.Duration
}

var _ db.AssignmentDB = &pgAssignments{}

type assignmentIDs struct {
	UserID    string `json:"user_id"`
	RoleID    string `json:"user_id"`
	ProjectID string `json:"project_id"`
}

func (a pgAssignments) SetOrg(ctx context.Context, email models.Email, label models.Label) (*models.Assignment, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	var ids assignmentIDs
	IDsQuery := `select json_build_object(
		'user_id', users.data->>'id',
		'role_id', roles.data->>'id')

		from users, roles
		where 
			users.data->>'email' = $1 AND
			roles.data->>'label' = $2;`

	row := a.pool.QueryRow(ctx, IDsQuery, email, label)
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
	_, err = a.pool.Exec(ctx, s, assignment)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, db.ErrDuplicateKey
		}

		return nil, fmt.Errorf("error creating contributor: %w", err)
	}

	return &assignment, nil
}

func (a pgAssignments) SetProject(ctx context.Context, email models.Email, label models.Label, label2 models.Label) (*models.Assignment, error) {
	panic("implement me")
}

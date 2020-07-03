package capepg

import (
	"context"
	"errors"
	"time"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pgRole struct {
	pool    *pgxpool.Pool
	timeout time.Duration
}

var _ db.RoleDB = &pgRole{}

func (r *pgRole) Create(context.Context, models.Role) error {
	return errors.New("not implemented")
}

func (r *pgRole) Delete(context.Context, models.Label) error {
	return errors.New("not implemented")
}

func (r *pgRole) Get(context.Context, models.Label) (models.Role, error) {
	return models.Role{}, errors.New("not implemented")
}

func (r *pgRole) List(context.Context, db.ListRoleOptions) ([]models.Role, error) {
	return nil, errors.New("not implemented")
}

func (r *pgRole) AttachPolicy(context.Context, models.Label) error {
	return errors.New("not implemented")
}

func (r *pgRole) DetachPolicy(context.Context, models.Label) error {
	return errors.New("not implemented")
}

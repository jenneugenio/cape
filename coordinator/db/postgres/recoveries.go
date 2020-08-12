package capepg

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"time"
)

type pgRecovery struct {
	pool    Pool
	timeout time.Duration
}

var _ db.RecoveryDB = &pgRecovery{}

type dbRecovery struct {
	*models.Recovery
	Credentials *models.Credentials `json:"credentials"`
}

func (p *pgRecovery) Get(ctx context.Context, ID string) (*models.Recovery, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var recovery dbRecovery
	s := "select data from recoveries where id = $1;"
	row := p.pool.QueryRow(ctx, s, ID)
	err := row.Scan(&recovery)
	if err != nil {
		return nil, err
	}

	return &models.Recovery{
		ID:          recovery.ID,
		UserID:      recovery.UserID,
		Credentials: recovery.Credentials,
		ExpiresAt:   recovery.ExpiresAt,
	}, nil
}

func (p *pgRecovery) Create(ctx context.Context, recovery models.Recovery) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	r := dbRecovery{
		Recovery:    &recovery,
		Credentials: recovery.Credentials,
	}

	s := "insert into recoveries (data) values ($1);"
	_, err := p.pool.Exec(ctx, s, r)
	return err
}

func (p *pgRecovery) Delete(ctx context.Context, ID string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := "delete from recoveries where id = $1;"
	_, err := p.pool.Exec(ctx, s, ID)
	return err
}

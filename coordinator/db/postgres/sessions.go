package capepg

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"time"
)

type pgSession struct {
	pool    Pool
	timeout time.Duration
}

var _ db.SessionDB = &pgSession{}

func (p *pgSession) Get(ctx context.Context, ID string) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var session models.Session
	s := "select data from sessions where id = $1;"
	row := p.pool.QueryRow(ctx, s, ID)
	err := row.Scan(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (p *pgSession) Create(ctx context.Context, session models.Session) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := "insert into sessions (data) values ($1);"
	_, err := p.pool.Exec(ctx, s, session)
	return err
}

func (p *pgSession) Delete(ctx context.Context, ID string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := "delete from sessions where ID = $1;"
	_, err := p.pool.Exec(ctx, s, ID)
	return err
}

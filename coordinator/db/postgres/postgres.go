package capepg

import (
	"context"
	"time"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// CapePg is a postgresql implementation of the cape database interface
type CapePg struct {
	pool    Pool
	timeout time.Duration
}

var _ db.Interface = &CapePg{}

func New(pool Pool) *CapePg {
	return &CapePg{
		pool:    pool,
		timeout: 5 * time.Second,
	}
}

func (c *CapePg) Roles() db.RoleDB               { return &pgRole{c.pool, c.timeout} }
func (c *CapePg) Contributors() db.ContributorDB { return &pgContributor{c.pool, c.timeout} }
func (c *CapePg) Projects() db.ProjectsDB        { return &pgProject{c.pool, c.timeout} }
func (c *CapePg) Users() db.UserDB               { return &pgUser{c.pool, c.timeout} }
func (c *CapePg) Config() db.ConfigDB            { return &pgConfig{c.pool, c.timeout} }
func (c *CapePg) Secrets() db.SecretDB           { return &pgSecret{c.pool, c.timeout} }
func (c *CapePg) Tokens() db.TokensDB            { return &pgToken{c.pool, c.timeout} }
func (c *CapePg) Session() db.SessionDB          { return &pgSession{c.pool, c.timeout} }
func (c *CapePg) Recoveries() db.RecoveryDB      { return &pgRecovery{c.pool, c.timeout} }

type Pool interface {
	Exec(_ context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(_ context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(_ context.Context, sql string, args ...interface{}) pgx.Row
}

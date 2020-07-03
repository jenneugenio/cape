package capepg

import (
	"time"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

// CapePg is a postgresql implementation of the cape database interface
type CapePg struct {
	pool    *pgxpool.Pool
	timeout time.Duration
}


var _ db.Interface = &CapePg{}

func New(pool *pgxpool.Pool) *CapePg {
	return &CapePg{
		pool: pool,
		timeout: 5 * time.Second,
	}
}

func (c *CapePg) Policies() db.PolicyDB { return &pgPolicy{c.pool, c.timeout} }
func (c *CapePg) Roles() db.RoleDB      { return &pgRole{c.pool, c.timeout} }

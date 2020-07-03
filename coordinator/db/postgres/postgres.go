package capepg

import (
	"database/sql"
	"github.com/capeprivacy/cape/coordinator/db"
)

// CapePg is a postgresql implementation of the cape database interface
type CapePg struct {
	db *sql.DB
}

func New(db *sql.DB) *CapePg {
	return &CapePg{
		db: db,
	}
}

func (c *CapePg) Policies() db.PolicyDB { return &pgPolicy{c.db} }
func (c *CapePg) Roles() db.RoleDB      { return &pgRole{c.db} }

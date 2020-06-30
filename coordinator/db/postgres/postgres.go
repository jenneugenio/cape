package capepg

import "database/sql"

// CapePg is a postgresql implementation of the cape database interface
type CapePg struct {
	db *sql.DB
}

var _ db.Interface = &CapePg{}

func New(db *sql.DB) *CapePg {
	return &CapePg{
		db: db,
	}
}

func (c *CapePg) Policies() db.PolicyDB { return &pgPolicy{c.db} }
func (c *CapePg) Roles() db.RoleDB      { return &pgRole{c.db} }

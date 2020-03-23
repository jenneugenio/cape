package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// Role in a role in the system (e.g. Admin, user, etc)
type Role struct {
	*database.Primitive
	Label string `json:"label"`
}

// GetType returns the type of this entity
func (r *Role) GetType() types.Type {
	return RoleType
}

// NewRole returns a mutable role struct
func NewRole(label string) (*Role, error) {
	p, err := database.NewPrimitive(RoleType)
	if err != nil {
		return nil, err
	}

	return &Role{
		Primitive: p,
		Label:     label,
	}, nil
}

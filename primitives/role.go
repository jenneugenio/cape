package primitives

import (
	"github.com/dropoutlabs/privacyai/primitives/types"
)

// Role in a role in the system (e.g. Admin, user, etc)
type Role struct {
	*Primitive
}

// GetType returns the type of this entity
func (r *Role) GetType() types.Type {
	return RoleType
}

// NewRole returns a mutable role struct
func NewRole() (*Role, error) {
	p, err := newPrimitive(RoleType)
	if err != nil {
		return nil, err
	}

	return &Role{
		Primitive: p,
	}, nil
}

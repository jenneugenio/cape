package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

const (
	// AdminRole is the label of the admin role
	AdminRole = Label("admin")

	// GlobalRole is the label of the global role
	GlobalRole = Label("global")

	// DataConnectorRole is the label of the data connector roles
	DataConnectorRole = Label("data-connector")
)

var SystemRoles = []Label{AdminRole, GlobalRole, DataConnectorRole}

// Role in a role in the system (e.g. Admin, user, etc)
type Role struct {
	*database.Primitive
	Label  Label `json:"label"`
	System bool  `json:"system"`
}

// GetType returns the type of this entity
func (r *Role) GetType() types.Type {
	return RoleType
}

// NewRole returns a mutable role struct
func NewRole(label Label, system bool) (*Role, error) {
	p, err := database.NewPrimitive(RoleType)
	if err != nil {
		return nil, err
	}

	return &Role{
		Primitive: p,
		Label:     label,
		System:    system,
	}, nil
}

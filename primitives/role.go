package primitives

import (
	"fmt"

	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/oklog/ulid"
)

const (
	// AdminRole is the label of the admin role
	AdminRole = Label("admin")

	// GlobalRole is the label of the global role
	GlobalRole = Label("global")
)

var SystemRoles = []Label{AdminRole, GlobalRole}

// Role in a role in the system (e.g. Admin, user, etc)
type Role struct {
	ID     ulid.ULID
	Label  Label `json:"label"`
	System bool  `json:"system"`
}

func (r *Role) Validate() error {
	if err := r.Label.Validate(); err != nil {
		return errors.Wrap(InvalidRoleCause, err)
	}

	return nil
}

// GetType returns the type of this entity
func (r *Role) GetType() types.Type {
	return RoleType
}

// NewRole returns a mutable role struct
func NewRole(id string, label Label, system bool) (*Role, error) {
	fmt.Println("HI", id, "HELLO")
	u, err := ulid.Parse(id)
	if err != nil {
		fmt.Println(id, err)
		return nil, err
	}
	role := &Role{
		ID:     u,
		Label:  label,
		System: system,
	}

	return role, role.Validate()
}

func (r *Role) GetEncryptable() bool {
	return false
}

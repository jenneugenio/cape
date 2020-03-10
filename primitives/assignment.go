package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// Assignment represents a policy being applied/attached to a role
type Assignment struct {
	*database.Primitive
	IdentityID database.ID `json:"identity_id"`
	RoleID     database.ID `json:"role_id"`
}

// GetType returns the type for this entity
func (a *Assignment) GetType() types.Type {
	return AssignmentType
}

// NewAssignment returns a new Assignment
func NewAssignment(identityID, roleID database.ID) (*Assignment, error) {
	p, err := database.NewPrimitive(AssignmentType)
	if err != nil {
		return nil, err
	}

	// TODO: Pass in the values of the Assignment!
	//
	// An Assignment is considered an immutable type in our object system (as
	// defined by the type)
	a := &Assignment{
		Primitive:  p,
		IdentityID: identityID,
		RoleID:     roleID,
	}

	ID, err := database.DeriveID(a)
	if err != nil {
		return nil, err
	}

	a.ID = ID
	return a, nil
}

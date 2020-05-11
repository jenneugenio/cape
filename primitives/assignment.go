package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Assignment represents a policy being applied/attached to a role
type Assignment struct {
	*database.Primitive
	IdentityID database.ID `json:"identity_id"`
	RoleID     database.ID `json:"role_id"`
}

func (a *Assignment) Validate() error {
	if err := a.Primitive.Validate(); err != nil {
		return errors.Wrap(InvalidAssignmentCause, err)
	}

	if err := a.IdentityID.Validate(); err != nil {
		return errors.New(InvalidAssignmentCause, "Assignment identity id must be valid")
	}

	typ, err := a.IdentityID.Type()
	if err != nil {
		return errors.New(InvalidAssignmentCause, "Invalid Identity ID provided")
	}

	if typ != UserType && typ != ServicePrimitiveType {
		return errors.New(InvalidAssignmentCause, "Invalid Identity ID provided")
	}

	if err := a.RoleID.Validate(); err != nil {
		return errors.New(InvalidAssignmentCause, "Assignment role id must be valid")
	}

	typ, err = a.RoleID.Type()
	if err != nil {
		return errors.New(InvalidAssignmentCause, "Invalid Role ID provider")
	}

	if typ != RoleType {
		return errors.New(InvalidAssignmentCause, "Invalid Role ID provider")
	}

	return nil
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
	return a, a.Validate()
}

func (a *Assignment) GetEncryptable() bool {
	return false
}

package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Attachment represents a policy being applied/attached to a role
type Attachment struct {
	*database.Primitive
	PolicyID string      `json:"policy_id"`
	RoleID   database.ID `json:"role_id"`
}

func (a *Attachment) Validate() error {
	if err := a.Primitive.Validate(); err != nil {
		return errors.Wrap(InvalidAttachmentCause, err)
	}

	if err := a.RoleID.Validate(); err != nil {
		return errors.New(InvalidAttachmentCause, "Attachment role ID must be valid")
	}

	typ, err := a.RoleID.Type()
	if err != nil {
		return errors.New(InvalidAttachmentCause, "Invalid Role ID provided")
	}

	if typ != RoleType {
		return errors.New(InvalidAttachmentCause, "Invalid Role ID provided")
	}

	return nil
}

// GetType returns the type for this entity
func (a *Attachment) GetType() types.Type {
	return AttachmentType
}

// NewAttachment returns a new attachment
func NewAttachment(policyID string, roleID database.ID) (*Attachment, error) {
	p, err := database.NewPrimitive(AttachmentType)
	if err != nil {
		return nil, err
	}

	// An attachment is considered an immutable type in our object system (as
	// defined by the type)
	a := &Attachment{
		Primitive: p,
		PolicyID:  policyID,
		RoleID:    roleID,
	}

	ID, err := database.DeriveID(a)
	if err != nil {
		return nil, err
	}

	a.ID = ID
	return a, a.Validate()
}

func (a *Attachment) GetEncryptable() bool {
	return false
}

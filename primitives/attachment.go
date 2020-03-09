package primitives

import (
	"github.com/dropoutlabs/privacyai/database"
	"github.com/dropoutlabs/privacyai/database/types"
)

// Attachment represents a policy being applied/attached to a role
type Attachment struct {
	*database.Primitive
	PolicyID database.ID `json:"policy_id"`
	RoleID   database.ID `json:"role_id"`
}

// GetType returns the type for this entity
func (a *Attachment) GetType() types.Type {
	return AttachmentType
}

// NewAttachment returns a new attachment
func NewAttachment(policyID, roleID database.ID) (*Attachment, error) {
	p, err := database.NewPrimitive(AttachmentType)
	if err != nil {
		return nil, err
	}

	// TODO: Pass in the values of the attachment!
	//
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
	return a, nil
}

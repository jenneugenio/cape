package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/oklog/ulid"
)

// Attachment represents a policy being applied/attached to a role
type Attachment struct {
	*database.Primitive
	PolicyID database.ID `json:"policy_id"`
	RoleID   ulid.ULID   `json:"role_id"`
}

func (a *Attachment) Validate() error {
	if err := a.Primitive.Validate(); err != nil {
		return errors.Wrap(InvalidAttachmentCause, err)
	}

	if err := a.PolicyID.Validate(); err != nil {
		return errors.New(InvalidAttachmentCause, "Attachment policy id must be valid")
	}

	typ, err := a.PolicyID.Type()
	if err != nil {
		return errors.New(InvalidAttachmentCause, "Invalid Policy ID provider")
	}

	if typ != PolicyType {
		return errors.New(InvalidAttachmentCause, "Invalid Policy ID provider")
	}

	return nil
}

// GetType returns the type for this entity
func (a *Attachment) GetType() types.Type {
	return AttachmentType
}

// NewAttachment returns a new attachment
func NewAttachment(policyID database.ID, roleID string) (*Attachment, error) {
	p, err := database.NewPrimitive(AttachmentType)
	if err != nil {
		return nil, err
	}

	u, err := ulid.Parse(roleID)
	if err != nil {
		return nil, err
	}

	// An attachment is considered an immutable type in our object system (as
	// defined by the type)
	a := &Attachment{
		Primitive: p,
		PolicyID:  policyID,
		RoleID:    u,
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

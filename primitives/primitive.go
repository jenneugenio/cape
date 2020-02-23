package primitives

import (
	"time"

	"github.com/dropoutlabs/privacyai/primitives/types"
)

// Primitive a low level object that can be inherited from for getting access
// to common values and methods.
type Primitive struct {
	ID        ID         `json:"id"`
	Type      types.Type `json:"type"`
	Version   uint8      `json:"version"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// GetID satisifes the Entity interface to return an ID
func (p *Primitive) GetID() ID {
	return p.ID
}

// GetType satisfies the Entity interface to return the type of the struct
func (p *Primitive) GetType() types.Type {
	return p.Type
}

// GetVersion satisfies the Entity interface to return a Version.
//
// If a primitive has a non-1 version it should implement this version itself
func (p *Primitive) GetVersion() uint8 {
	return 1
}

// GetCreatedAt returns the CreatedAt value for this struct
func (p *Primitive) GetCreatedAt() time.Time {
	return p.CreatedAt
}

// GetUpdatedAt returns the UpdateAt value for this struct
func (p *Primitive) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

// newPrimitive returns a new primitive entity object
//
// If the type is mutable then this function will also generate an ID. For
// immutable types the caller must manually lock the struct by deriving an ID.
func newPrimitive(t types.Type) (*Primitive, error) {
	p := &Primitive{
		Version:   1,
		Type:      t,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if t.Mutable() {
		ID, err := GenerateID(p)
		if err != nil {
			return nil, err
		}

		p.ID = ID
	}

	return p, nil
}

var (
	// The following types represent the core primitives in the system
	UserType       types.Type = 0x000
	ServiceType    types.Type = 0x001
	TokenType      types.Type = 0x002
	RoleType       types.Type = 0x003
	PolicyType     types.Type = 0x004
	AttachmentType types.Type = 0x005
)

func init() {
	types.Register(UserType, "user", false)
	types.Register(ServiceType, "service", false)
	types.Register(TokenType, "token", true)
	types.Register(RoleType, "role", false)
	types.Register(PolicyType, "policy", false)
	types.Register(AttachmentType, "attachment", true)
}

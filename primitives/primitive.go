package primitives

import (
	"time"

	"github.com/dropoutlabs/privacyai/primitives/types"
)

// Entity represents any primitive data structure stored inside the Controller.
// All primitives must satisfy this interface to be stored in the ata layer.
type Entity interface {
	GetID() ID
	GetType() types.Type
	GetVersion() int
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

// Primitive a low level object that can be inherited from for getting access
// to common values and methods.
//
// Primitives _must_ implement their own `GetType` value.
type Primitive struct {
	ID        ID         `json:"id"`
	Type      types.Type `json:"type"`
	Version   int        `json:"version"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// GetID satisifes the Entity interface to return an ID
func (p *Primitive) GetID() ID {
	return p.ID
}

// GetVersion satisfies the Entity interface to return a Version
func (p *Primitive) GetVersion() int {
	return p.Version
}

// GetCreatedAt returns the CreatedAt value for this struct
func (p *Primitive) GetCreatedAt() time.Time {
	return p.CreatedAt
}

// GetUpdatedAt returns the UpdateAt value for this struct
func (p *Primitive) GetUpatedAt() time.Time {
	return p.UpdatedAt
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
	types.Register(RoleType, "role", true)
	types.Register(PolicyType, "policy", true)
	types.Register(AttachmentType, "attachment", true)
}

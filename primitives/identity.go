package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// Identity represents an identity type such as user or service
type Identity interface {
	GetID() database.ID
	GetType() types.Type
	GetCredentials() *Credentials
	GetEmail() string
}

// IdentityImpl implements the above interface and is a common
// container for identity data
type IdentityImpl struct {
	*database.Primitive
	Credentials *Credentials `json:"credentials"`
	Email       string       `json:"email"`
}

// GetCredentials implements Identity interface
func (i *IdentityImpl) GetCredentials() *Credentials {
	return i.Credentials
}

// GetEmail implements Identity interface
func (i *IdentityImpl) GetEmail() string {
	return i.Email
}

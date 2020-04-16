package primitives

import (
	"github.com/capeprivacy/cape/database"
)

// Identity represents an identity type such as user or service
type Identity interface {
	database.Entity
	GetCredentials() *Credentials
	GetEmail() Email
}

// IdentityImpl implements the above interface and is a common
// container for identity data. Its required to deserialize the
// common data in services and users.
type IdentityImpl struct {
	*database.Primitive
	Credentials *Credentials `json:"credentials"`
	Email       Email        `json:"email"`
}

// GetCredentials implements Identity interface
func (i *IdentityImpl) GetCredentials() *Credentials {
	return i.Credentials
}

// GetEmail implements Identity interface
func (i *IdentityImpl) GetEmail() Email {
	return i.Email
}

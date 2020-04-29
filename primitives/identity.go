package primitives

import (
	"github.com/capeprivacy/cape/database"
)

// Identity represents an identity type such as user or service
type Identity interface {
	database.Entity
	GetEmail() Email
}

// IdentityImpl implements the above interface and is a common
// container for identity data. Its required to deserialize the
// common data in services and users.
type IdentityImpl struct {
	*database.Primitive
	Email Email `json:"email"`
}

// GetEmail implements Identity interface
func (i *IdentityImpl) GetEmail() Email {
	return i.Email
}

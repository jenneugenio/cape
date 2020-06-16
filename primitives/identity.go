package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
)

// Identity represents an identity type such as user or service
type Identity interface {
	database.Entity
	GetEmail() Email
	GetName() Name
}

// IdentityImpl implements the above interface and is a common
// container for identity data. Its required to deserialize the
// common data in services and users.
type IdentityImpl struct {
	*database.Primitive
	Email Email `json:"email"`
	Name  Name  `json:"name"`
}

// GetEmail implements Identity interface
func (i *IdentityImpl) GetEmail() Email {
	return i.Email
}

func (i *IdentityImpl) GetName() Name {
	return i.Name
}

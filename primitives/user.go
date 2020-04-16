package primitives

import (
	"github.com/capeprivacy/cape/database"
	"github.com/capeprivacy/cape/database/types"
)

// User represents a user of the system
type User struct {
	*IdentityImpl
	Name Name `json:"name"`
}

// GetType returns the type for this entity
func (u *User) GetType() types.Type {
	return UserType
}

// NewUser returns a new User struct
func NewUser(name Name, email Email, creds *Credentials) (*User, error) {
	p, err := database.NewPrimitive(UserType)
	if err != nil {
		return nil, err
	}

	return &User{
		Name: name, // TODO: Figure out what to do about validation
		IdentityImpl: &IdentityImpl{
			Primitive:   p,
			Email:       email,
			Credentials: creds,
		},
	}, nil
}

// GetCredentials satisfies Identity interface
func (u *User) GetCredentials() *Credentials {
	return u.Credentials
}

// GetEmail satisfies the Identity interface
func (u *User) GetEmail() Email {
	return u.Email
}

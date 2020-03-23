package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// User represents a user of the system
type User struct {
	*IdentityImpl
	Name string `json:"name"`
}

// GetType returns the type for this entity
func (u *User) GetType() types.Type {
	return UserType
}

// NewUser returns a new User struct
func NewUser(name string, email string, creds *Credentials) (*User, error) {
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
func (u *User) GetEmail() string {
	return u.Email
}

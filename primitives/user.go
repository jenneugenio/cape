package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

// User represents a user of the system
type User struct {
	*database.Primitive
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	Credentials *Credentials `json:"credentials"`
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
		Primitive:   p,
		Name:        name, // TODO: Figure out what to do about validation
		Email:       email,
		Credentials: creds,
	}, nil
}

// GetCredentials satisfies Identity interface
func (u *User) GetCredentials() *Credentials {
	return u.Credentials
}

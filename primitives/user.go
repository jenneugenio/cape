package primitives

import (
	"github.com/dropoutlabs/privacyai/primitives/types"
)

// User represents a user of the system
type User struct {
	*Primitive
	Name string `json:"name"`
}

// GetType returns the type for this entity
func (u *User) GetType() types.Type {
	return UserType
}

// NewUser returns a new User struct
func NewUser(name string) (*User, error) {
	p, err := newPrimitive(UserType)
	if err != nil {
		return nil, err
	}

	return &User{
		Primitive: p,
		Name:      name, // TODO: Figure out what to do about validation
	}, nil
}

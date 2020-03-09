package primitives

import (
	"github.com/dropoutlabs/privacyai/database"
	"github.com/dropoutlabs/privacyai/database/types"
)

// User represents a user of the system
type User struct {
	*database.Primitive
	Name string `json:"name"`
}

// GetType returns the type for this entity
func (u *User) GetType() types.Type {
	return UserType
}

// NewUser returns a new User struct
func NewUser(name string) (*User, error) {
	p, err := database.NewPrimitive(UserType)
	if err != nil {
		return nil, err
	}

	return &User{
		Primitive: p,
		Name:      name, // TODO: Figure out what to do about validation
	}, nil
}

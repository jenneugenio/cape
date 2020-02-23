package primitives

import (
	"github.com/dropoutlabs/privacyai/primitives/types"
)

// User represents a user of the system
type User struct {
	Primitive
	Name string
}

// GetType returns the type for this primitive
func (u *User) GetType() types.Type {
	return UserType
}

package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// User represents a user of the system
type User struct {
	*IdentityImpl
	Name Name `json:"name"`
}

func (u *User) Validate() error {
	if err := u.Primitive.Validate(); err != nil {
		return errors.Wrap(InvalidUserCause, err)
	}

	if err := u.Email.Validate(); err != nil {
		return errors.Wrap(InvalidUserCause, err)
	}

	if err := u.Name.Validate(); err != nil {
		return errors.Wrap(InvalidUserCause, err)
	}

	if err := u.Credentials.Validate(); err != nil {
		return errors.Wrap(InvalidUserCause, err)
	}

	return nil
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

	user := &User{
		Name: name, // TODO: Figure out what to do about validation
		IdentityImpl: &IdentityImpl{
			Primitive:   p,
			Email:       email,
			Credentials: creds,
		},
	}

	return user, user.Validate()
}

// GetCredentials satisfies Identity interface
func (u *User) GetCredentials() *Credentials {
	return u.Credentials
}

// GetEmail satisfies the Identity interface
func (u *User) GetEmail() Email {
	return u.Email
}

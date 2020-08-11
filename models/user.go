package models

import (
	"time"
)

func (u *User) GetCredentials() (*Credentials, error) {
	return &Credentials{
		Secret: u.Credentials.Secret,
		Salt:   u.Credentials.Salt,
		Alg:    u.Credentials.Alg,
	}, nil
}

func (u *User) GetUserID() string {
	return u.ID
}

func (u *User) GetStringID() string {
	return u.ID
}

// User represents a user of the system
type User struct {
	ID        string    `json:"id"`
	Version   uint8     `json:"version"`
	Email     Email     `json:"email"`
	Name      Name      `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// We never want to send Credentials over the wire!
	Credentials Credentials `json:"credentials" gqlgen:"-"`
}

// NewUser returns a new User struct
func NewUser(name Name, email Email, creds Credentials) User {
	user := User{
		ID:          NewID(),
		Credentials: creds,
		Email:       email,
		Name:        name,
		CreatedAt:   now(),
	}

	return user
}

// GenerateUser returns an instantiated user for use in unit testing
//
// This function _should only ever_ be used inside of a test.
func GenerateUser(name, email string) (Password, User) {
	password := GeneratePassword()

	n := Name(name)
	e := Email(email)

	c := GenerateCredentials()

	user := NewUser(n, e, Credentials{
		Secret: c.Secret,
		Salt:   c.Salt,
		Alg:    c.Alg,
	})
	return password, user
}

package models

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/capeprivacy/cape/primitives"
	"github.com/manifoldco/go-base64"
)

// CredentialsAlgType enum holding the supported crypto algorithms
type CredentialsAlgType string

var (
	// UnknownAlg represents the empty state of the CredentialsAlgType
	UnknownAlg CredentialsAlgType = ""

	// EDDSA is used for generating asymmetric keypairs for signing tokens and
	// other items within the cape ecosystem.
	EDDSA CredentialsAlgType = "eddsa"

	// Argon2ID exists for production usage, it's the most recent winner of the
	// Password Hashing Competition and is tuned inside of the auth package.
	Argon2ID CredentialsAlgType = "argon2id"

	// SHA256 only exists for internal testing, it should never be used in any
	// production scenario.
	//
	// SHA256 is used as a password hashing algorithm
	SHA256 CredentialsAlgType = "sha256"
)

// String returns the string represented by the enum value
func (c *CredentialsAlgType) String() string {
	return string(*c)
}

// UnmarshalGQL unmarshals a string in the CredentialsAlgType enum
func (c *CredentialsAlgType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New("cannot unmarshal CredentialsAlgType")
	}

	*c = CredentialsAlgType(str)

	return nil
}

// MarshalGQL marshals a CredentailsAlgType enum to string
func (c CredentialsAlgType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

type Credentials struct {
	// Secret is the "scrypt'ed" secret which we store in the database.
	Secret *base64.Value      `json:"secret"`
	Salt   *base64.Value      `json:"salt"`
	Alg    CredentialsAlgType `json:"alg"`
}

func (u *User) GetCredentials() (*primitives.Credentials, error) {
	return &primitives.Credentials{
		Secret: u.Credentials.Secret,
		Salt:   u.Credentials.Salt,
		Alg:    primitives.CredentialsAlgType(u.Credentials.Alg),
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
func GenerateUser(name, email string) (primitives.Password, User) {
	password := primitives.GeneratePassword()

	n := Name(name)
	e := Email(email)

	c := primitives.GenerateCredentials()

	user := NewUser(n, e, Credentials{
		Secret: c.Secret,
		Salt:   c.Salt,
		Alg:    CredentialsAlgType(c.Alg),
	})
	return password, user
}

package primitives

import (
	"context"
	"encoding/json"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// User represents a user of the system
type User struct {
	*IdentityImpl

	// We never want to send Credentials over the wire!
	Credentials *Credentials `json:"-" gqlgen:"-"`
}

type encryptedUser struct {
	*User
	Credentials *base64.Value `json:"credentials"`
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
		Credentials: creds,
		IdentityImpl: &IdentityImpl{
			Primitive: p,
			Email:     email,
			Name:      name,
		},
	}

	return user, user.Validate()
}

func (u *User) GetIdentityID() database.ID {
	return u.ID
}

// GetCredentials satisfies Identity interface
func (u *User) GetCredentials() (*Credentials, error) {
	return u.Credentials, nil
}

// GetEmail satisfies the Identity interface
func (u *User) GetEmail() Email {
	return u.Email
}

func (u *User) Encrypt(ctx context.Context, codec crypto.EncryptionCodec) ([]byte, error) {
	creds, err := json.Marshal(u.Credentials)
	if err != nil {
		return nil, err
	}

	data, err := codec.Encrypt(ctx, base64.New(creds))
	if err != nil {
		return nil, err
	}

	return json.Marshal(encryptedUser{
		User:        u,
		Credentials: data,
	})
}

func (u *User) Decrypt(ctx context.Context, codec crypto.EncryptionCodec, data []byte) error {
	in := &encryptedUser{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	unencrypted, err := codec.Decrypt(ctx, in.Credentials)
	if err != nil {
		return err
	}

	creds := &Credentials{}
	err = json.Unmarshal([]byte(*unencrypted), creds)
	if err != nil {
		return err
	}

	u.IdentityImpl = in.IdentityImpl
	u.Credentials = creds
	return nil
}

func (u *User) GetEncryptable() bool {
	return true
}

// GenerateUser returns an instantiated user for use in unit testing
//
// This function _should only ever_ be used inside of a test.
func GenerateUser(name, email string) (Password, *User, error) {
	password, err := GeneratePassword()
	if err != nil {
		return EmptyPassword, nil, err
	}

	n, err := NewName(name)
	if err != nil {
		return EmptyPassword, nil, err
	}

	e, err := NewEmail(email)
	if err != nil {
		return EmptyPassword, nil, err
	}

	c, err := GenerateCredentials()
	if err != nil {
		return EmptyPassword, nil, err
	}

	user, err := NewUser(n, e, c)
	return password, user, err
}

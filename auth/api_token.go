package auth

import (
	"fmt"
	"github.com/capeprivacy/cape/models"
	"strings"

	"github.com/manifoldco/go-base64"

	errors "github.com/capeprivacy/cape/partyerrors"
)

const (
	tokenVersion = 0x01
	secretBytes  = 24
)

// Secret represents a Secret stored inside of an APIToken
type Secret []byte

// Validate returns whether or not the underlying secret is valid
func (s Secret) Validate() error {
	if len([]byte(s)) != secretBytes {
		return errors.New(BadSecretLength, "APIToken secret length is not the correct length")
	}

	return nil
}

func (s Secret) String() string {
	return string(s)
}

func (s Secret) Password() models.Password {
	return models.Password(base64.New(s).String())
}

func FromPassword(password models.Password) (Secret, error) {
	value, err := base64.NewFromString(password.String())
	if err != nil {
		return Secret([]byte{}), err
	}

	s := Secret(*value)

	return s, s.Validate()
}

// APIToken represents a token that is used by a user
// to authenticate with a coordinator.
type APIToken struct {
	TokenID string
	Version byte
	Secret  Secret
}

// NewAPIToken returns a new api token from email and url
func NewAPIToken(secret Secret, tokenCredentialID string) (*APIToken, error) {
	return &APIToken{
		TokenID: tokenCredentialID,
		Version: tokenVersion,
		Secret:  secret,
	}, nil
}

// Validate returns an error if the underlying APIToken has invalid contents in
// its fields.
func (a *APIToken) Validate() error {
	if a.Version != tokenVersion {
		return errors.New(BadAPITokenVersion, "Expected version: %x", tokenVersion)
	}

	if err := a.Secret.Validate(); err != nil {
		return err
	}

	return nil
}

// Marshal marshals the api token into a string.
// Format of the output is {email},{version}|{secret}|{url}
// {version}|{secret}|{url} are bytes concatenated together and
// encoded as base64
func (a *APIToken) Marshal() (string, error) {
	bytes := make([]byte, 1+secretBytes)
	bytes[0] = a.Version
	copy(bytes[1:], []byte(a.Secret))
	copy(bytes[1:secretBytes+1], []byte(a.Secret))

	val := base64.New(bytes)
	tokenStr := fmt.Sprintf("%s,%s", a.TokenID, val.String())

	return tokenStr, nil
}

// Unmarshal unmarshals the string into the APIToken struct
func (a *APIToken) Unmarshal(token string) error {
	strs := strings.Split(token, ",")
	if len(strs) != 2 {
		return errors.New(BadTokenFormat, "Invalid API Token provided")
	}

	ID := strs[0]

	val, err := base64.NewFromString(strs[1])
	if err != nil {
		return err
	}

	tokenBytes := []byte(*val)

	a.TokenID = ID
	a.Version = tokenBytes[0]
	a.Secret = tokenBytes[1:]

	return a.Validate()
}

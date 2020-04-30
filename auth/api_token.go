package auth

import (
	"crypto/rand"
	"fmt"
	"github.com/capeprivacy/cape/database"
	"strings"

	"github.com/capeprivacy/cape/primitives"

	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/manifoldco/go-base64"
)

const (
	tokenVersion = 0x01
	secretBytes  = 16
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

// APIToken represents a token that is used by a service or user
// to authenticate with a coordinator. Currently we're using the email
// so that we can use the normal user login flow but in the future
// the APIToken will be tied to a token (token_id will replace email)
// that is tied to an identity (user or service)
type APIToken struct {
	TokenCredentialID database.ID
	URL               *primitives.URL
	Version           byte
	Secret            Secret
}

// NewAPIToken returns a new api token from email and url
func NewAPIToken(secret Secret, tokenCredentialID database.ID, u *primitives.URL) (*APIToken, error) {
	secretBytes := make([]byte, secretBytes)
	_, err := rand.Read(secretBytes)
	if err != nil {
		return nil, err
	}

	return &APIToken{
		TokenCredentialID: tokenCredentialID,
		URL:               u,
		Version:           tokenVersion,
		Secret:            secret,
	}, nil
}

// Validate returns an error if the underlying APIToken has invalid contents in
// its fields.
func (a *APIToken) Validate() error {
	if err := a.URL.Validate(); err != nil {
		return err
	}

	if a.Version != tokenVersion {
		return errors.New(BadAPITokenVersion, "Expected version: %x", tokenVersion)
	}

	if err := a.Secret.Validate(); err != nil {
		return err
	}

	return nil
}

// Credentials returns Credentials with the secret stored on
// the auth token.
func (a *APIToken) Credentials() (*Credentials, error) {
	creds, err := NewCredentials(a.Secret, nil)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

// Marshal marshals the api token into a string.
// Format of the output is {email},{version}|{secret}|{url}
// {version}|{secret}|{url} are bytes concatenated together and
// encoded as base64
func (a *APIToken) Marshal() (string, error) {
	urlBytes := []byte(a.URL.String())

	bytes := make([]byte, len(urlBytes)+1+secretBytes)
	bytes[0] = a.Version

	copy(bytes[1:secretBytes+1], []byte(a.Secret))
	copy(bytes[secretBytes+1:], urlBytes)

	val := base64.New(bytes)

	tokenStr := fmt.Sprintf("%s,%s", a.TokenCredentialID, val.String())

	return tokenStr, nil
}

// Unmarshal unmarshals the string into the APIToken struct
func (a *APIToken) Unmarshal(token string) error {
	strs := strings.Split(token, ",")
	if len(strs) != 2 {
		return errors.New(BadTokenFormat, "Invalid API Token provided")
	}

	tokenCredentialID, err := database.DecodeFromString(strs[0])
	if err != nil {
		return err
	}

	a.TokenCredentialID = tokenCredentialID

	val, err := base64.NewFromString(strs[1])
	if err != nil {
		return err
	}

	tokenBytes := []byte(*val)

	a.Version = tokenBytes[0]
	a.Secret = Secret(tokenBytes[1 : secretBytes+1])

	u, err := primitives.NewURL(string(tokenBytes[secretBytes+1:]))
	if err != nil {
		return err
	}

	a.URL = u
	return a.Validate()
}

// Parse returns an APIToken from a given string and validates the underlying
// APIToken is sensical.
func ParseAPIToken(in string) (*APIToken, error) {
	token := &APIToken{}
	err := token.Unmarshal(in)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Random secret will return a random password
func RandomSecret() (Secret, error) {
	secret := make([]byte, secretBytes)
	_, err := rand.Read(secret)
	return secret, err
}

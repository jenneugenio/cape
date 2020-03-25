package auth

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"strings"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/manifoldco/go-base64"
)

const (
	tokenVersion = 0x01
	secretBytes  = 16
)

// APIToken represents a token that is used by a service or user
// to authenticate with a controller. Currently we're using the email
// so that we can use the normal user login flow but in the future
// the APIToken will be tied to a token (token_id will replace email)
// that is tied to an identity (user or service)
type APIToken struct {
	Email   string
	URL     *url.URL
	Version byte
	Secret  []byte
}

// NewAPIToken returns a new api token from email and url
func NewAPIToken(email string, url *url.URL) (*APIToken, error) {
	secretBytes := make([]byte, secretBytes)
	_, err := rand.Read(secretBytes)
	if err != nil {
		return nil, err
	}

	return &APIToken{
		Email:   email,
		URL:     url,
		Version: tokenVersion,
		Secret:  secretBytes,
	}, nil
}

// Marshal marshals the api token into a string.
// Format of the output is {email},{version}|{secret}|{url}
// {version}|{secret}|{url} are bytes concatenated together and
// encoded as base64
func (a *APIToken) Marshal() (string, error) {
	urlBytes := []byte(a.URL.String())

	bytes := make([]byte, len(urlBytes)+1+secretBytes)
	bytes[0] = a.Version

	copy(bytes[1:secretBytes+1], a.Secret)
	copy(bytes[secretBytes+1:], urlBytes)

	val := base64.New(bytes)

	tokenStr := fmt.Sprintf("%s,%s", a.Email, val.String())

	return tokenStr, nil
}

// Unmarshal unmarshals the string into the APIToken struct
func (a *APIToken) Unmarshal(token string) error {
	strs := strings.Split(token, ",")
	if len(strs) > 2 {
		return errors.New(BadTokenFormat, "API Token only has two parts, email and base64 encoded data")
	}

	a.Email = strs[0]

	val, err := base64.NewFromString(strs[1])
	if err != nil {
		return err
	}

	tokenBytes := []byte(*val)

	a.Version = tokenBytes[0]
	a.Secret = tokenBytes[1 : secretBytes+1]

	u, err := url.Parse(string(tokenBytes[secretBytes+1:]))
	if err != nil {
		return err
	}

	a.URL = u

	return nil
}

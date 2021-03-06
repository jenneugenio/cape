package crypto

import (
	"crypto/rand"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/manifoldco/go-base64"

	errors "github.com/capeprivacy/cape/partyerrors"
)

type KeyURLType string

const (
	Base64Key KeyURLType = "base64key"
	AzureKey  KeyURLType = "azurekeyvault"
)

func (k KeyURLType) Validate() error {
	switch k {
	case Base64Key:
		return nil
	case AzureKey:
		return nil
	default:
		return errors.New(InvalidKeyURLCause, "Invalid scheme got %s", k)
	}
}

// NewKeyURL parses the given string and returns a key url.
func NewKeyURL(in string) (*KeyURL, error) {
	u, err := url.Parse(in)
	if err != nil {
		return nil, err
	}

	d := &KeyURL{URL: u}
	return d, d.Validate()
}

// KeyURLFromURL returns a KeyURL from a net/url.URL
func KeyURLFromURL(u *url.URL) (*KeyURL, error) {
	d := &KeyURL{URL: u}

	return d, d.Validate()
}

// KeyURL contains a url for a key
// Used for data encryption and mostly intended to be used with
// gocloud secrets and related utilities
type KeyURL struct {
	*url.URL
}

// Validate returns an error if the uri is not a valid key uri
func (d *KeyURL) Validate() error {
	if d.URL == nil {
		return errors.New(InvalidKeyURLCause, "Missing db url")
	}

	if d.URL.Host == "" {
		return errors.New(InvalidKeyURLCause, "A host must be provided")
	}

	typ := KeyURLType(d.Scheme)

	return typ.Validate()
}

// ToURL returns the underlying url.URL
func (d *KeyURL) ToURL() *url.URL {
	return d.URL
}

// MarshalJSON implements the JSON.Marshaller interface
func (d *KeyURL) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(d.String())), nil
}

// UnmarshalJSON implements the JSON.Unmarshaller interface
func (d *KeyURL) UnmarshalJSON(b []byte) error {
	var keyURL string
	err := json.Unmarshal(b, &keyURL)
	if err != nil {
		return err
	}

	u, err := url.Parse(keyURL)
	if err != nil {
		return err
	}

	d.URL = u
	return d.Validate()
}

func (d *KeyURL) Type() KeyURLType {
	return KeyURLType(d.Scheme)
}

// NewBase64KeyURL generates a KeyURL for you if key is nil. Otherwise
// key must be KeyLength long.
func NewBase64KeyURL(key []byte) (*KeyURL, error) {
	if key == nil {
		key = make([]byte, KeyLength)

		_, err := rand.Read(key)
		if err != nil {
			return nil, err
		}
	}

	if len(key) != KeyLength {
		return nil, errors.New(InvalidKeyURLCause, "Key must be %d bytes long", KeyLength)
	}

	return NewKeyURL("base64key://" + base64.New(key).String())
}

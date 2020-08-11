package models

import (
	"crypto/rand"
	"fmt"
	"github.com/manifoldco/go-base64"
	"io"
	"strconv"
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
		return fmt.Errorf("cannot unmarshal CredentialsAlgType")
	}

	*c = CredentialsAlgType(str)

	return nil
}

// MarshalGQL marshals a CredentailsAlgType enum to string
func (c CredentialsAlgType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

const (
	SecretLength = 32
	SaltLength   = 16
)

type Credentials struct {
	// Secret is the "scrypt'ed" secret which we store in the database.
	Secret *base64.Value      `json:"secret"`
	Salt   *base64.Value      `json:"salt"`
	Alg    CredentialsAlgType `json:"alg"`
}

// GenerateCredentials returns an instantiated Credentials for use in unit testing.
//
// This function _should only ever_ be used inside of a test.
func GenerateCredentials() *Credentials {
	secret := make([]byte, SecretLength)
	salt := make([]byte, SaltLength)

	_, err := rand.Read(secret)
	if err != nil {
		panic("Unable to read from rand")
	}

	_, err = rand.Read(salt)
	if err != nil {
		panic("Unable to read from rand")
	}

	return &Credentials{
		Secret: base64.New(secret),
		Salt:   base64.New(salt),
		Alg:    SHA256,
	}
}

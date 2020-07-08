package primitives

import (
	"fmt"
	"io"
	"strconv"

	errors "github.com/capeprivacy/cape/partyerrors"
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

// Validate checks to see if the algorithm type is valid
func (c *CredentialsAlgType) Validate() error {
	switch *c {
	case Argon2ID:
		return nil
	case SHA256:
		return nil
	case EDDSA:
		return nil
	default:
		return errors.New(InvalidAlgType, "%s is not a valid CredentialsAlgType", *c)
	}
}

// UnmarshalGQL unmarshals a string in the CredentialsAlgType enum
func (c *CredentialsAlgType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New(InvalidAlgType, "Cannot unmarshal CredentialsAlgType")
	}

	*c = CredentialsAlgType(str)

	return c.Validate()
}

// MarshalGQL marshals a CredentailsAlgType enum to string
func (c CredentialsAlgType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

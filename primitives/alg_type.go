package primitives

import (
	"fmt"
	"io"
	"strconv"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

// CredentialsAlgType enum holding the supported crypto algorithms
type CredentialsAlgType string

var (
	// EDDSA algorithm type
	EDDSA CredentialsAlgType = "EDDSA"
)

// String returns the string represented by the enum value
func (c *CredentialsAlgType) String() string {
	return string(*c)
}

// UnmarshalGQL unmarshals a string in the CredentialsAlgType enum
func (c *CredentialsAlgType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New(InvalidAlgType, "Cannot unmarshal credentials algorithm type")
	}

	*c = CredentialsAlgType(str)
	if *c != EDDSA {
		return errors.New(InvalidTokenType, "%s is not a valid CredentialsAlgType", str)
	}
	return nil
}

// MarshalGQL marshals a CredentailsAlgType enum to string
func (c CredentialsAlgType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

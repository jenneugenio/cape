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
	// EDDSA algorithm type
	EDDSA CredentialsAlgType = "EDDSA"
)

// String returns the string represented by the enum value
func (c *CredentialsAlgType) String() string {
	return string(*c)
}

// Validate checks to see if the service type is valid
func (c *CredentialsAlgType) Validate() error {
	switch *c {
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

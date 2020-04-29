package primitives

import (
	"fmt"
	"io"
	"strconv"

	errors "github.com/capeprivacy/cape/partyerrors"
)

// TokenType is an enum holding the category of sessions
type TokenType string

var (
	// EmailLogin is the session type used during the login flow
	Login TokenType = "LOGIN"
	// Authenticated is the session type used on normal API calls
	Authenticated TokenType = "AUTHENTICATED"
)

// Validate checks to see if the service type is valid
func (t *TokenType) Validate() error {
	switch *t {
	case Login:
		return nil
	case Authenticated:
		return nil
	default:
		return errors.New(InvalidTokenType, "%s is not a valid TokenType", t)
	}
}

// String returns the string represented by the enum value
func (t *TokenType) String() string {
	return string(*t)
}

// UnmarshalGQL unmarshals a string in the TokenType enum
func (t *TokenType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New(InvalidTokenType, "Cannot unmarshal token type")
	}

	*t = TokenType(str)
	return t.Validate()
}

// MarshalGQL marshals a CredentailsAlgType enum to string
func (t TokenType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}

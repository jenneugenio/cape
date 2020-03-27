package primitives

import (
	"fmt"
	"io"
	"strconv"

	"github.com/badoux/checkmail"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

// Email represents a valid email for use within Cape
type Email string

// NewEmail validates that the string is a valid label before returning an email
func NewEmail(in string) (Email, error) {
	e := Email(in)
	err := e.Validate()
	return e, err
}

// Validate returns an error if the contents of the label are invalid
func (e Email) Validate() error {
	err := checkmail.ValidateFormat(string(e))
	if err != nil {
		return errors.New(InvalidEmail, "A valid email address must be provided")
	}

	return nil
}

// String returns the string representation of the label
func (e Email) String() string {
	return string(e)
}

func (e *Email) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return errors.New(InvalidEmail, "Cannot unmarshall provided ID")
	}

	email, err := NewEmail(s)
	if err != nil {
		return err
	}

	*e = email
	return nil
}

func (e Email) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
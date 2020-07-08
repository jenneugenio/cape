package primitives

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/badoux/checkmail"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// EmailType is an enum holding representing a user
type EmailType string

var (
	// UserEmail represents a user email
	UserEmail EmailType = "user"
)

// Email represents a valid email for use within Cape
type Email struct {
	Email string
	Type  EmailType
}

// NewEmail validates that the string is a valid label before returning an email
func NewEmail(in string) (Email, error) {
	typ := UserEmail
	e := Email{
		Email: in,
		Type:  typ,
	}

	return e, e.Validate()
}

// Validate returns an error if the contents of the label are invalid
func (e Email) Validate() error {
	s := e.String()

	err := checkmail.ValidateFormat(s)
	if err != nil {
		return errors.New(InvalidEmail, fmt.Sprintf("A valid email address must be provided %s", s))
	}

	return nil
}

// String returns the string representation of the label
func (e Email) String() string {
	return e.Email
}

// SetType sets the email type
func (e *Email) SetType(typ EmailType) {
	e.Type = typ
}

// MarshalJSON implements the JSON.Marshaller interface
func (e Email) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(e.String())), nil
}

// UnmarshalJSON implements the JSON.Unmarshaller interface
func (e *Email) UnmarshalJSON(b []byte) error {
	emailStr := ""
	err := json.Unmarshal(b, &emailStr)
	if err != nil {
		return err
	}

	email, err := NewEmail(emailStr)
	if err != nil {
		return err
	}

	*e = email

	return nil
}

func (e *Email) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return errors.New(InvalidEmail, "Cannot unmarshall provided Email")
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

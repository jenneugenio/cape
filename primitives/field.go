package primitives

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
	"regexp"
)

// Field represents a column in a database
type Field string

var fieldRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

// Validate the field
func (f Field) Validate() error {
	if !fieldRegex.MatchString(string(f)) {
		msg := "field must start with a letter, and then only contain letters, numbers, or underscores"
		return errors.New(InvalidFieldCause, msg)
	}

	return nil
}

// NewField validates and returns a new field
func NewField(in string) (Field, error) {
	field := Field(in)
	return field, field.Validate()
}
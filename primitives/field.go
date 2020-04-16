package primitives

import (
	errors "github.com/capeprivacy/cape/partyerrors"
	"regexp"
)

// Field represents a column in a database
type Field string

const Star Field = "*"

var fieldRegex = regexp.MustCompile(`^(\*|[a-zA-Z]+[a-zA-Z0-9_]*)$`)

// Validate the field
func (f Field) Validate() error {
	if !fieldRegex.MatchString(string(f)) {
		msg := "field must start with a letter, and then only contain letters, numbers, or underscores, or *"
		return errors.New(InvalidFieldCause, msg)
	}

	return nil
}

// String turns the field into a string
func (f Field) String() string {
	return string(f)
}

// NewField validates and returns a new field
func NewField(in string) (Field, error) {
	field := Field(in)
	return field, field.Validate()
}

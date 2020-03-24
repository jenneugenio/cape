package primitives

import (
	"regexp"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

var labelRegex = regexp.MustCompile("^[a-z][a-z/-//]{3,64}$")

// Label represents a uri-safe identifier name for an entity within the Cape
// ecosystem. Labels are generally unique.
type Label string

// NewLabel validates that the string is a valid label before returning a label
func NewLabel(in string) (Label, error) {
	v := Label(in)
	err := v.Validate()
	return v, err
}

// Validate returns an error if the contents of the label are invalid
func (l Label) Validate() error {
	if !labelRegex.MatchString(string(l)) {
		msg := "Labels must only contain alphabetical a-z, -, or /. They must start with a-z and be between 4 and 64 characters in length."
		return errors.New(InvalidLabelCause, msg)
	}

	return nil
}

// String returns the string representation of the label
func (l Label) String() string {
	return string(l)
}

package primitives

import (
	"regexp"

	errors "github.com/capeprivacy/cape/partyerrors"
)

var labelRegex = regexp.MustCompile("^[a-z0-9][a-z0-9/-]{3,64}$")

// Label represents a uri-safe identifier name for an entity within the Cape
// ecosystem. Labels are generally unique.
type Label string

// NewLabel validates that the string is a valid label before returning a label
func NewLabel(in string) (Label, error) {
	v := Label(in)
	return v, v.Validate()
}

// Validate returns an error if the contents of the label are invalid
func (l Label) Validate() error {
	if !labelRegex.MatchString(string(l)) {
		msg := "Labels must only contain 0-9, a-z, or -. They must start with a letter and be between 4 and 64 characters in length."
		return errors.New(InvalidLabelCause, msg)
	}

	return nil
}

// String returns the string representation of the label
func (l Label) String() string {
	return string(l)
}

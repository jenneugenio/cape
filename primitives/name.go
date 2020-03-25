package primitives

import (
	"regexp"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

var nameRegex = regexp.MustCompile("^.{2,64}$")

// Name represents a users name in Cape
type Name string

// NewName validates that the string is a valid name before returning a Name
func NewName(in string) (Name, error) {
	n := Name(in)
	err := n.Validate()
	return n, err
}

// Validate returns an error if the name is not valid
func (n Name) Validate() error {
	if !nameRegex.MatchString(string(n)) {
		msg := "Names must only contain alphabetical characters and be between 2-64 characters in length."
		return errors.New(InvalidLabelCause, msg)
	}

	return nil
}

// String returns the string representation of the name
func (n Name) String() string {
	return string(n)
}

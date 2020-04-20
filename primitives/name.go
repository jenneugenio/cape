package primitives

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	errors "github.com/capeprivacy/cape/partyerrors"
)

var nameRegex = regexp.MustCompile("^.{2,64}$")

// Name represents a users name in Cape
type Name string

// NewName validates that the string is a valid name before returning a Name
func NewName(in string) (Name, error) {
	n := Name(in)
	return n, n.Validate()
}

// Validate returns an error if the name is not valid
func (n Name) Validate() error {
	if !nameRegex.MatchString(string(n)) {
		msg := "Names must only contain alphabetical characters and be between 2-64 characters in length."
		return errors.New(InvalidNameCause, msg)
	}

	return nil
}

// String returns the string representation of the name
func (n Name) String() string {
	return string(n)
}

func (n *Name) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return errors.New(InvalidNameCause, "Cannot unmarshall provided Name")
	}

	name, err := NewName(s)
	if err != nil {
		return err
	}

	*n = name

	return nil
}

func (n Name) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(n.String()))
}

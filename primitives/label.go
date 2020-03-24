package primitives

import (
	"encoding/json"
	"regexp"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

var labelRegex = regexp.MustCompile("^[a-z][a-z/-]{3,64}$")

// Label represents a uri-safe identifier name for an entity within the Cape
// ecosystem. Labels are generally unique.
type Label string

// UnmarshalJSON unmarshals the given byte stream of valid json into this struct
func (l Label) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &l)
	if err != nil {
		return err
	}

	return l.Validate()
}

// MarshalJSON marshals this struct into valid json
func (l Label) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(l))
}

// Validate returns an error if the contents of the label are invalid
func (l Label) Validate() error {
	if !labelRegex.MatchString(string(l)) {
		msg := "Labels must only contain alphabetical characters and dashes. They cannot be between 4 and 64 characters."
		return errors.New(InvalidLabelCause, msg)
	}

	return nil
}

// String returns the string representation of the label
func (l Label) String() string {
	return string(l)
}

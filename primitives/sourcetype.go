package primitives

import (
	"fmt"
	"io"
	"strconv"

	errors "github.com/capeprivacy/cape/partyerrors"
)

// SourceType enum holding the different types of supported data sources
type SourceType string

var (
	// PostgresType represents a Postgres database source
	PostgresType SourceType = "postgres"
)

// NewSourceType returns a source type for the given string or an invalid type
func NewSourceType(in string) (SourceType, error) {
	t := SourceType(in)
	return t, t.Validate()
}

// String returns the source type as a string
func (s SourceType) String() string {
	return string(s)
}

// Validate returns whether or not the SourceType is a valid
func (s SourceType) Validate() error {
	switch s {
	case PostgresType:
		return nil
	default:
		return errors.New(InvalidSourceType, "Invalid source type provided")
	}
}

// UnmarshalGQL unmarshals a string into the SourceType enum
func (s *SourceType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New(InvalidSourceType, "Invalid source type provided")
	}

	*s = SourceType(str)
	return s.Validate()
}

// MarshalGQL marshals a SourceType enum to a string
func (s SourceType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(s.String()))
}

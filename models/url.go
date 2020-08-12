package models

import (
	"fmt"
	"io"
	"net/url"
	"strconv"

	errors "github.com/capeprivacy/cape/partyerrors"
)

// NewURL parses the given string and returns a URL if the
// given URL is a valid coordinator url. If it's not an error is returned.
func NewURL(in string) (*URL, error) {
	u, err := url.Parse(in)
	if err != nil {
		return nil, err
	}

	c := &URL{URL: u}

	return c, c.Validate()
}

// URL contains a url to a Cape coordinator
type URL struct {
	*url.URL
}

// Validate returns an error if the coordinator url is _not_ valid
func (u *URL) Validate() error {
	if u.URL == nil {
		return errors.New(InvalidURLCause, "Missing url")
	}

	if u.URL.Scheme != "http" && u.URL.Scheme != "https" {
		return errors.New(InvalidURLCause, "Invalid scheme, must be http or https")
	}

	if u.URL.Host == "" {
		return errors.New(InvalidURLCause, "A host must be provided")
	}

	return nil
}

// Copy returns a copy of the URL
func (u *URL) Copy() (*URL, error) {
	return NewURL(u.URL.String())
}

// MarshalJSON implements the JSON.Marshaller interface
func (u *URL) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(u.URL.String())), nil
}

// UnmarshalJSON implements the JSON.Unmarshaller interface
func (u *URL) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != byte('"') || b[len(b)-1] != byte('"') {
		return errors.New(InvalidURLCause, "Invalid json provided")
	}

	out, err := url.Parse(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}

	u.URL = out
	return u.Validate()
}

// UnmarshalGQL unmarshalls a string from GraphQL into the URL
func (u *URL) UnmarshalGQL(v interface{}) error {
	switch s := v.(type) {
	case string:
		t, err := url.Parse(s)
		if err != nil {
			return err
		}

		u.URL = t
		return u.Validate()
	default:
		return errors.New(InvalidURLCause, "Invalid URL value provided, expected a string, got %T", s)
	}
}

// MarshalGQL marshals a URL to a strong for GraphQL
func (u URL) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(u.String()))
}

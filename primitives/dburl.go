package primitives

import (
	"fmt"
	"io"
	"net/url"
	"strconv"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

// NewDBURL parses the given string and returns a database url.
func NewDBURL(in string) (*DBURL, error) {
	u, err := url.Parse(in)
	if err != nil {
		return nil, err
	}

	d := &DBURL{URL: u}
	err = d.Validate()
	if err != nil {
		return nil, err
	}

	return d, nil
}

// DBURLFromURL returns a DBURL from a net/url.URL
func DBURLFromURL(u *url.URL) (*DBURL, error) {
	d := &DBURL{URL: u}
	err := d.Validate()
	if err != nil {
		return nil, err
	}

	return d, nil
}

// DBURL contains a url for a database
type DBURL struct {
	*url.URL
}

// Validate returns an error if the uri is not a valid database uri
func (d *DBURL) Validate() error {
	if d.URL == nil {
		return errors.New(InvalidDBURLCause, "Missing db url")
	}

	if d.URL.Scheme != "postgres" {
		return errors.New(InvalidDBURLCause, "Invalid scheme, only postgres is supported")
	}

	if d.URL.Host == "" {
		return errors.New(InvalidDBURLCause, "A host must be provided")
	}

	if d.URL.Path == "" || d.URL.Path == "/" {
		return errors.New(InvalidDBURLCause, "A database must be provided")
	}

	return nil
}

// ToURL returns the underlying url.URL
func (d *DBURL) ToURL() *url.URL {
	return d.URL
}

// SetPassword sets the password
func (d *DBURL) SetPassword(pw string) {
	d.User = url.UserPassword(d.User.Username(), pw)
}

// Copy creates a copy of this DBURL
func (d *DBURL) Copy() (*DBURL, error) {
	return NewDBURL(d.URL.String())
}

// MarshalJSON implements the JSON.Marshaller interface
func (d *DBURL) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(d.URL.String())), nil
}

// UnmarshalJSON implements the JSON.Unmarshaller interface
func (d *DBURL) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != byte('"') || b[len(b)-1] != byte('"') {
		return errors.New(InvalidURLCause, "Invalid json provided")
	}

	u, err := url.Parse(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}

	d.URL = u
	return d.Validate()
}

// UnmarshalGQL impements the interface required to marshal this type to GraphQL
func (d *DBURL) UnmarshalGQL(v interface{}) error {
	switch s := v.(type) {
	case string:
		u, err := url.Parse(s)
		if err != nil {
			return err
		}

		d.URL = u
		return d.Validate()
	default:
		return errors.New(InvalidDBURLCause, "Invalid DBURL value provided, expected a string, got %T", s)
	}
}

// MarshalGQL implements the interface required to unmarshal this type from GraphQL
func (d DBURL) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(d.URL.String()))
}

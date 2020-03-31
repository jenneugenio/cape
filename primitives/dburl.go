package primitives

import (
	"net/url"

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

// MarshalJSON implements the JSON.Marshaller interface
func (d *DBURL) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.URL.String() + "\""), nil
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

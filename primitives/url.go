package primitives

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

// NewURL parses the given string and returns a URL if the
// given URL is a valid controller url. If it's not an error is returned.
func NewURL(in string) (*URL, error) {
	u, err := url.Parse(in)
	if err != nil {
		return nil, err
	}

	c := &URL{URL: u}
	err = c.Validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewURLFromStdLib creates a new url from the std lib url
// type from net/url
func NewURLFromStdLib(u *url.URL) (*URL, error) {
	c := &URL{URL: u}
	err := c.Validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// URL contains a url to a Cape controller
type URL struct {
	*url.URL
}

// Validate returns an error if the controller url is _not_ valid
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
	return []byte("\"" + u.URL.String() + "\""), nil
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

// UnmarshalURL converts a url.URL into a string for usage in graphQL
func MarshalURL(u url.URL) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprint(w, strconv.Quote(u.String()))
	})
}

// UnmarshalURL converts a string into a url.URL for usage in graphQL
func UnmarshalURL(v interface{}) (url.URL, error) {
	switch v := v.(type) {
	case string:
		u, err := url.Parse(v)

		if err != nil {
			return url.URL{}, err
		}

		return *u, nil
	case map[string]interface{}:
		x, err := json.Marshal(v)
		if err != nil {
			return url.URL{}, err
		}

		var resp url.URL
		err = json.Unmarshal(x, &resp)
		if err != nil {
			return url.URL{}, err
		}

		return resp, nil
	default:
		return url.URL{}, fmt.Errorf("%T is not a string", v)
	}
}

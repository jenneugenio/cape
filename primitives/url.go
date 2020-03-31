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
func (c *URL) Validate() error {
	if c.URL == nil {
		return errors.New(InvalidURLCause, "Missing url")
	}

	if c.URL.Scheme != "http" && c.URL.Scheme != "https" {
		return errors.New(InvalidURLCause, "Invalid scheme, must be http or https")
	}

	if c.URL.Host == "" {
		return errors.New(InvalidURLCause, "A host must be provided")
	}

	return nil
}

// MarshalJSON implements the JSON.Marshaller interface
func (c *URL) MarshalJSON() ([]byte, error) {
	return []byte("\"" + c.URL.String() + "\""), nil
}

// UnmarshalJSON implements the JSON.Unmarshaller interface
func (c *URL) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != byte('"') || b[len(b)-1] != byte('"') {
		return errors.New(InvalidURLCause, "Invalid json provided")
	}

	u, err := url.Parse(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}

	c.URL = u
	return c.Validate()
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

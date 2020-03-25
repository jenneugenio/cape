package primitives

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"io"
	"net/url"
	"strconv"
)

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

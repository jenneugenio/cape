package primitives

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/manifoldco/go-base64"
)

// MarshalBase64Value marshals a base64 value for graphql
func MarshalBase64Value(v base64.Value) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprint(w, strconv.Quote(v.String()))
	})
}

// UnmarshalBase64Value unmarshal a base64 value from graphql
func UnmarshalBase64Value(v interface{}) (base64.Value, error) {
	switch v := v.(type) {
	case string:
		b, err := base64.NewFromString(v)
		fmt.Println("DOESN'T WORK BISH", b, err, v)
		if err != nil {
			return nil, err
		}
		return *b, nil
	default:
		return nil, fmt.Errorf("%T is not a string", v)
	}
}

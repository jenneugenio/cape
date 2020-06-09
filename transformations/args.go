package transformations

import (
	"encoding/json"

	errors "github.com/capeprivacy/cape/partyerrors"
)

// Args represents the arguments to be passed into
// a transformation
type Args map[string]interface{}

// LookupFloat64 looks up a float64 number on the args. If not present,
// returns 0, false and nil. If present, attempts to convert it into float64
// or json.Number. If neither of these then returns an error.
func (a Args) LookupFloat64(key string) (float64, bool, error) {
	val, ok := a[key]
	if !ok {
		return 0, false, nil
	}

	switch v := val.(type) {
	case float64:
		return v, true, nil
	case json.Number:
		float, err := v.Float64()
		if err != nil {
			return 0, false, err
		}
		return float, true, nil
	}

	return 0, false, errors.New(UnsupportedType, "Unsupported type for %q: found %T expected float64", key, val)
}

package sources

import (
	"encoding/json"
	"strconv"
	"time"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

// PostgresEncode encodes data out of postgres into Strings
func PostgresEncode(values []interface{}) ([][]byte, error) {
	outVals := make([][]byte, len(values))
	for i, val := range values {
		var outVal string
		switch v := val.(type) {
		case int64:
			outVal = strconv.FormatInt(v, 10)
		case int32:
			outVal = strconv.FormatInt(int64(v), 10)
		case time.Time:
			outVal = v.Format(time.RFC3339Nano)
		case float64:
			outVal = strconv.FormatFloat(v, 'E', -1, 64)
		case string:
			outVal = v
		default:
			return nil, errors.New(UnknownFieldType, "Unknown type %T", v)
		}

		by, err := json.Marshal(outVal)
		if err != nil {
			return nil, err
		}

		outVals[i] = by
	}

	return outVals, nil
}

package sources

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/dropoutlabs/cape/connector/proto"
	errors "github.com/dropoutlabs/cape/partyerrors"
)

// Record contains the values decoded by our decoder
type Record struct {
	values []interface{}
}

// NewRecord decodes the incoming bytes given the schema and returns a Record
func NewRecord(schema *proto.Schema, fields []*proto.Field) (*Record, error) {
	values, err := Decode(schema, fields)
	if err != nil {
		return nil, err
	}

	return &Record{
		values: values,
	}, nil
}

// Values returns the underlying values
func (r *Record) Values() []interface{} {
	return r.values
}

// ToStrings encodes the underlying types to
// their string format
func (r *Record) ToStrings() ([]string, error) {
	strs := make([]string, len(r.values))
	for i, val := range r.values {
		var outVal string
		switch t := val.(type) {
		case int64:
			outVal = strconv.FormatInt(t, 10)
		case int32:
			outVal = strconv.FormatInt(int64(t), 10)
		case int16:
			outVal = strconv.FormatInt(int64(t), 10)
		case float64:
			outVal = strconv.FormatFloat(t, 'f', -1, 64)
		case float32:
			outVal = strconv.FormatFloat(float64(t), 'f', -1, 32)
		case time.Time:
			outVal = t.Format(time.RFC3339Nano)
		case string:
			outVal = t
		case bool:
			outVal = strconv.FormatBool(t)
		case []byte:
			dst := make([]byte, hex.EncodedLen(len(t)))
			hex.Encode(dst, t)
			outVal = string(dst)
		default:
			return nil, errors.New(UnknownFieldType, "Unknown type %T", t)
		}

		strs[i] = outVal
	}

	return strs, nil
}

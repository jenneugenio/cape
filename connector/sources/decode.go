package sources

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/dropoutlabs/cape/connector/proto"
)

// Decode decodes a byte stream from a data connector
func Decode(schema *proto.Schema, values [][]byte) ([]interface{}, error) {
	outVals := make([]interface{}, len(values))
	for i, val := range values {
		var tmpVal string

		err := json.Unmarshal(val, &tmpVal)
		if err != nil {
			return nil, err
		}

		var outVal interface{}
		switch schema.Fields[i].Field {
		case proto.FieldType_BIGINT:
			v, err := strconv.ParseInt(tmpVal, 10, 64)
			if err != nil {
				return nil, err
			}
			outVal = v
		case proto.FieldType_INT:
			v, err := strconv.ParseInt(tmpVal, 10, 32)
			if err != nil {
				return nil, err
			}
			outVal = int32(v)
		case proto.FieldType_TIMESTAMP:
			t, err := time.Parse(time.RFC3339Nano, tmpVal)
			if err != nil {
				return nil, err
			}

			outVal = t
		case proto.FieldType_DOUBLE:
			v, err := strconv.ParseFloat(tmpVal, 64)
			if err != nil {
				return nil, err
			}
			outVal = v
		case proto.FieldType_TEXT:
			outVal = tmpVal
		}

		outVals[i] = outVal
	}

	return outVals, nil
}

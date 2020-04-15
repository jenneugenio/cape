package sources

import (
	"github.com/golang/protobuf/ptypes"

	"github.com/dropoutlabs/cape/connector/proto"
)

// Decode decodes a byte stream from a data connector
func Decode(schema *proto.Schema, values []*proto.Field) ([]interface{}, error) {
	outVals := make([]interface{}, len(values))
	for i, val := range values {
		switch schema.Fields[i].Field {
		case proto.FieldType_BIGINT:
			outVals[i] = val.GetInt64()
		case proto.FieldType_INT:
			outVals[i] = val.GetInt32()
		case proto.FieldType_TIMESTAMP:
			ts, err := ptypes.Timestamp(val.GetTimestamp())
			if err != nil {
				return nil, err
			}
			outVals[i] = ts
		case proto.FieldType_DOUBLE:
			outVals[i] = val.GetDouble()
		case proto.FieldType_TEXT:
			outVals[i] = val.GetString_()
		}
	}

	return outVals, nil
}

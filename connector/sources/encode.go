package sources

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	pb "github.com/dropoutlabs/cape/connector/proto"
	errors "github.com/dropoutlabs/cape/partyerrors"
)

// PostgresEncode encodes data out of postgres into Strings
func PostgresEncode(values []interface{}) ([]*pb.Field, error) {
	outVals := make([]*pb.Field, len(values))
	for i, val := range values {
		switch v := val.(type) {
		case int64:
			outVals[i] = &pb.Field{
				Value: &pb.Field_Int64{Int64: v},
			}
		case int32:
			outVals[i] = &pb.Field{
				Value: &pb.Field_Int32{Int32: v},
			}
		case time.Time:
			ts, err := ptypes.TimestampProto(v)
			if err != nil {
				return nil, err
			}

			outVals[i] = &pb.Field{
				Value: &pb.Field_Timestamp{Timestamp: ts},
			}
		case float64:
			outVals[i] = &pb.Field{
				Value: &pb.Field_Double{Double: v},
			}
		case string:
			outVals[i] = &pb.Field{
				Value: &pb.Field_String_{String_: v},
			}
		default:
			return nil, errors.New(UnknownFieldType, "Unknown type %T", v)
		}
	}

	return outVals, nil
}

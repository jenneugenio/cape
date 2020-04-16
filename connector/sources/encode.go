package sources

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	pb "github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// PostgresEncode encodes data out of postgres into Strings
func PostgresEncode(values []interface{}) ([]*pb.Field, error) {
	outVals := make([]*pb.Field, len(values))
	for i, val := range values {
		outVals[i] = &pb.Field{}
		switch v := val.(type) {
		case int64:
			outVals[i].Value = &pb.Field_Int64{Int64: v}
		case int32:
			outVals[i].Value = &pb.Field_Int32{Int32: v}
		case int16:
			outVals[i].Value = &pb.Field_Int32{Int32: int32(v)}
		case time.Time:
			ts, err := ptypes.TimestampProto(v)
			if err != nil {
				return nil, err
			}

			outVals[i].Value = &pb.Field_Timestamp{Timestamp: ts}
		case float64:
			outVals[i].Value = &pb.Field_Double{Double: v}
		case float32:
			outVals[i].Value = &pb.Field_Float{Float: v}
		case string:
			outVals[i].Value = &pb.Field_String_{String_: v}
		case bool:
			outVals[i].Value = &pb.Field_Bool{Bool: v}
		case []byte:
			outVals[i].Value = &pb.Field_Bytes{Bytes: v}
		default:
			return nil, errors.New(UnknownFieldType, "Unknown type %T", v)
		}
	}

	return outVals, nil
}

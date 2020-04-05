package sources

import (
	"encoding/json"
	"testing"
	"time"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/connector/proto"
)

func TestRecordsToStrings(t *testing.T) {
	gm.RegisterTestingT(t)

	schema := &proto.Schema{
		DataSource: "transactions",
		Target:     "transactions",
		Type:       proto.RecordType_DOCUMENT,
		Fields: []*proto.Field{
			{
				Field: proto.FieldType_BIGINT,
			},
			{
				Field: proto.FieldType_INT,
			},
			{
				Field: proto.FieldType_DOUBLE,
			},
			{
				Field: proto.FieldType_TEXT,
			},
			{
				Field: proto.FieldType_TIMESTAMP,
			},
		},
	}

	var input []interface{}
	input = append(input, int64(12345677777))
	input = append(input, int32(42))
	input = append(input, float64(1000.1000))
	input = append(input, "MASTERCARD")
	input = append(input, time.Now())

	data, err := PostgresEncode(input)
	gm.Expect(err).To(gm.BeNil())

	record, err := NewRecord(schema, data)
	gm.Expect(err).To(gm.BeNil())

	strs := record.ToStrings()
	for i, str := range strs {
		var expected string
		err = json.Unmarshal(data[i], &expected)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(str).To(gm.Equal(expected))
	}
}

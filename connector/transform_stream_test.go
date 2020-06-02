package connector

import (
	"testing"

	gm "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/primitives"
)

func TestTransformStream(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("stream runs transforms", func(t *testing.T) {
		schema := &proto.Schema{
			Fields: []*proto.FieldInfo{
				{
					Field: proto.FieldType_BIGINT,
					Name:  "my-field",
					Size:  8,
				},
			},
		}

		transform := []*primitives.Transformation{
			{
				Field:    "my-field",
				Function: "plusOne",
				Args:     nil,
			},
		}

		backingStream := &testStream{}
		stream, err := NewTransformStream(backingStream, schema, transform)
		gm.Expect(err).To(gm.BeNil())

		field := &proto.Field{Value: &proto.Field_Int64{Int64: 64}}
		record := &proto.Record{Fields: []*proto.Field{field}, Schema: schema}

		expectedField := &proto.Field{Value: &proto.Field_Int64{Int64: 65}}
		expectedRecord := &proto.Record{Fields: []*proto.Field{expectedField}, Schema: schema}

		err = stream.Send(record)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(backingStream.Buffer[0]).To(gm.Equal(expectedRecord))
	})

	t.Run("transform with where", func(t *testing.T) {
		schema := &proto.Schema{
			Fields: []*proto.FieldInfo{
				{
					Field: proto.FieldType_BIGINT,
					Name:  "my_field",
					Size:  8,
				},
			},
		}

		transform := []*primitives.Transformation{
			{
				Field:    "my_field",
				Function: "plusOne",
				Args:     nil,
				Where:    "my_field == 64",
			},
		}

		backingStream := &testStream{}
		stream, err := NewTransformStream(backingStream, schema, transform)
		gm.Expect(err).To(gm.BeNil())

		field := &proto.Field{Value: &proto.Field_Int64{Int64: 64}}
		record := &proto.Record{Fields: []*proto.Field{field}, Schema: schema}

		expectedField := &proto.Field{Value: &proto.Field_Int64{Int64: 64}}
		expectedRecord := &proto.Record{Fields: []*proto.Field{expectedField}, Schema: schema}

		err = stream.Send(record)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(backingStream.Buffer[0]).To(gm.Equal(expectedRecord))
	})

	newErrorTests := []struct {
		name      string
		schema    *proto.Schema
		transform []*primitives.Transformation
		errStr    string
	}{
		{
			name: "transform on unsupported type",
			schema: &proto.Schema{
				Fields: []*proto.FieldInfo{
					{
						Field: proto.FieldType_VARCHAR,
						Name:  "my-field",
						Size:  8,
					},
				},
			},
			transform: []*primitives.Transformation{
				{
					Field:    "my-field",
					Function: "plusOne",
					Args:     nil,
				},
			},
			errStr: "unsupported_type: Attempted to call plusOne transform on an unsupported type VARCHAR",
		},
		{
			name: "transform on non-existent field",
			schema: &proto.Schema{
				Fields: []*proto.FieldInfo{
					{
						Field: proto.FieldType_VARCHAR,
						Name:  "my-field",
						Size:  8,
					},
				},
				Target: "cool-target",
			},
			transform: []*primitives.Transformation{
				{
					Field:    "non-existent-field",
					Function: "plusOne",
					Args:     nil,
				},
			},
			errStr: "field_not_found: Could not find field non-existent-field for target cool-target",
		},
		{
			name: "transform with where that has syntax errors",
			schema: &proto.Schema{
				Fields: []*proto.FieldInfo{
					{
						Field: proto.FieldType_BIGINT,
						Name:  "my_field",
						Size:  8,
					},
				},
				Target: "cool-target",
			},
			transform: []*primitives.Transformation{
				{
					Field:    "my_field",
					Function: "plusOne",
					Args:     nil,
					Where:    "\"my_field == '0'\"",
				},
			},
			errStr: "Cannot transition token types from STRING [my_field == ] to NUMERIC [0]",
		},
	}

	for _, test := range newErrorTests {
		t.Run(test.name, func(t *testing.T) {
			backingStream := &testStream{}
			_, err := NewTransformStream(backingStream, test.schema, test.transform)
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(err.Error()).To(gm.Equal(test.errStr))
		})
	}

	t.Run("transform with where that is invalid", func(t *testing.T) {
		schema := &proto.Schema{
			DataSource: "cool-source",
			Fields: []*proto.FieldInfo{
				{
					Field: proto.FieldType_BIGINT,
					Name:  "my_field",
					Size:  8,
				},
			},
		}

		transform := []*primitives.Transformation{
			{
				Field:    "my_field",
				Function: "plusOne",
				Args:     nil,
				Where:    "hehe == 64",
			},
		}

		backingStream := &testStream{}
		stream, err := NewTransformStream(backingStream, schema, transform)
		gm.Expect(err).To(gm.BeNil())

		field := &proto.Field{Value: &proto.Field_Int64{Int64: 64}}
		record := &proto.Record{Fields: []*proto.Field{field}, Schema: schema}

		err = stream.Send(record)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("field_not_found: Could not evaluate where clause because 'hehe' is not a field in cool-source"))
	})
}

func TestFieldToInterface(t *testing.T) {
	ts := &timestamppb.Timestamp{Seconds: 343433, Nanos: 343433}
	tests := []struct {
		name  string
		field *proto.Field
	}{
		{
			name: "bool", field: &proto.Field{Value: &proto.Field_Bool{Bool: false}},
		},
		{
			name: "double", field: &proto.Field{Value: &proto.Field_Double{Double: 1.0}},
		},
		{
			name: "float", field: &proto.Field{Value: &proto.Field_Float{Float: 1.0}},
		},
		{
			name: "int32", field: &proto.Field{Value: &proto.Field_Int32{Int32: 1}},
		},
		{
			name: "int64", field: &proto.Field{Value: &proto.Field_Int64{Int64: 2}},
		},
		{
			name: "timestamp", field: &proto.Field{Value: &proto.Field_Timestamp{Timestamp: ts}},
		},
		{
			name: "bytes", field: &proto.Field{Value: &proto.Field_Bytes{Bytes: []byte("hey")}},
		},
		{
			name: "string", field: &proto.Field{Value: &proto.Field_String_{String_: "hey"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := fieldToInterface(test.field)
			gm.Expect(err).To(gm.BeNil())
		})
	}
}

package connector

import (
	"testing"

	"github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/primitives"
	gm "github.com/onsi/gomega"
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
		record := &proto.Record{Fields: []*proto.Field{field}}

		expectedField := &proto.Field{Value: &proto.Field_Int64{Int64: 65}}
		expectedRecord := &proto.Record{Fields: []*proto.Field{expectedField}}

		err = stream.Send(record)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(backingStream.Buffer[0]).To(gm.Equal(expectedRecord))
	})

	t.Run("transform on unsupported type", func(t *testing.T) {
		schema := &proto.Schema{
			Fields: []*proto.FieldInfo{
				{
					Field: proto.FieldType_VARCHAR,
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
		_, err := NewTransformStream(backingStream, schema, transform)
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("unsupported_type: Attempted to call plusOne transform on an unsupported type VARCHAR"))
	})

	t.Run("transform on non-existent field", func(t *testing.T) {
		schema := &proto.Schema{
			Fields: []*proto.FieldInfo{
				{
					Field: proto.FieldType_VARCHAR,
					Name:  "my-field",
					Size:  8,
				},
			},
			Target: "cool-target",
		}

		transform := []*primitives.Transformation{
			{
				Field:    "non-existent-field",
				Function: "plusOne",
				Args:     nil,
			},
		}

		backingStream := &testStream{}
		_, err := NewTransformStream(backingStream, schema, transform)
		gm.Expect(err).NotTo(gm.BeNil())
		gm.Expect(err.Error()).To(gm.Equal("field_not_found: Could not find field non-existent-field for target cool-target"))
	})
}

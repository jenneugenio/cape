package transformations

import (
	"testing"

	"github.com/capeprivacy/cape/connector/proto"
	gm "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestConditionals(t *testing.T) {
	gm.RegisterTestingT(t)

	schema := &proto.Schema{
		Fields: []*proto.FieldInfo{
			{
				Field: proto.FieldType_BIGINT,
				Name:  "hey",
				Size:  8,
			},
		},
	}

	transform, _ := NewPlusOneTransform("hey")

	tests := []struct {
		name     string
		inputVal int64
	}{
		{
			name:     "test does transform",
			inputVal: 0,
		},
		{
			name:     "test doesn't transform",
			inputVal: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := NewConditionalTransformation("hey == 0", transform)
			gm.Expect(err).To(gm.BeNil())

			inputField := &proto.Field{Value: &proto.Field_Int64{Int64: test.inputVal}}
			expectedOutputField := &proto.Field{Value: &proto.Field_Int64{Int64: 1}}
			record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: schema}
			expectedRecord := &proto.Record{Fields: []*proto.Field{expectedOutputField}, Schema: schema}

			err = c.Transform(schema, record)
			gm.Expect(err).To(gm.BeNil())

			gm.Expect(record).To(gm.Equal(expectedRecord))
		})
	}

	t.Run("new error", func(t *testing.T) {
		_, err := NewConditionalTransformation("\"hey == '0'\"", transform)
		gm.Expect(err).NotTo(gm.BeNil())
	})

	errorTests := []struct {
		name       string
		expression Condition
		schema     *proto.Schema
	}{
		{
			name:       "bad schema",
			expression: "hey == 0",
			schema: &proto.Schema{
				Fields: []*proto.FieldInfo{
					{
						Field: proto.FieldType_BIGINT,
						Name:  "hello",
						Size:  8,
					},
				},
			},
		},
		{
			name:       "wrong return type",
			expression: "hey + 5",
			schema: &proto.Schema{
				Fields: []*proto.FieldInfo{
					{
						Field: proto.FieldType_BIGINT,
						Name:  "hey",
						Size:  8,
					},
				},
			},
		},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			c, err := NewConditionalTransformation(test.expression, transform)
			gm.Expect(err).To(gm.BeNil())

			inputField := &proto.Field{Value: &proto.Field_Int64{Int64: 0}}
			record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: schema}

			err = c.Transform(test.schema, record)
			gm.Expect(err).ToNot(gm.BeNil())
		})
	}
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

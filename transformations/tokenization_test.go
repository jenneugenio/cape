package transformations

import (
	"testing"

	"github.com/capeprivacy/cape/connector/proto"
	gm "github.com/onsi/gomega"
)

var tokSchema = &proto.Schema{
	Fields: []*proto.FieldInfo{
		{
			Field: proto.FieldType_TEXT,
			Name:  "name",
			Size:  8,
		},
	},
}

func TestTokenizationArgs(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewTokenizationTransform("name")
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{})
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{
		"maxSize": 1,
	})
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{
		"maxSize": -1,
	})
	gm.Expect(err).NotTo(gm.BeNil())
}

func TestTokenizationString(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewTokenizationTransform("name")
	gm.Expect(err).To(gm.BeNil())
	args := Args{}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	transform.(*TokenizationTransform).key = []byte("secret_key")

	inputField := &proto.Field{Value: &proto.Field_String_{String_: "Jack"}}
	expectedToken := "81a7c769227edceaca2ed2bd320f87a5fbf504ef064b3dca8f2a9ed00125723f"
	expectedOutputField := &proto.Field{Value: &proto.Field_String_{String_: expectedToken}}
	record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: tokSchema}
	expectedRecord := &proto.Record{Fields: []*proto.Field{expectedOutputField}, Schema: tokSchema}

	err = transform.Transform(tokSchema, record)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(record).To(gm.Equal(expectedRecord))
}

func TestTokenizationByte(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewTokenizationTransform("name")
	gm.Expect(err).To(gm.BeNil())
	args := Args{}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	transform.(*TokenizationTransform).key = []byte("secret_key")

	inputField := &proto.Field{Value: &proto.Field_Bytes{Bytes: []byte("Jack")}}
	expectedToken := []byte("81a7c769227edceaca2ed2bd320f87a5fbf504ef064b3dca8f2a9ed00125723f")
	expectedOutputField := &proto.Field{Value: &proto.Field_Bytes{Bytes: expectedToken}}
	record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: tokSchema}
	expectedRecord := &proto.Record{Fields: []*proto.Field{expectedOutputField}, Schema: tokSchema}

	err = transform.Transform(tokSchema, record)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(record).To(gm.Equal(expectedRecord))
}

func TestTokenizationWithSize(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewTokenizationTransform("name")
	gm.Expect(err).To(gm.BeNil())
	args := Args{
		"maxSize": 10,
	}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	transform.(*TokenizationTransform).key = []byte("secret_key")

	inputField := &proto.Field{Value: &proto.Field_String_{String_: "Jack"}}
	expectedToken := "81a7c76922"
	expectedOutputField := &proto.Field{Value: &proto.Field_String_{String_: expectedToken}}
	record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: tokSchema}
	expectedRecord := &proto.Record{Fields: []*proto.Field{expectedOutputField}, Schema: tokSchema}

	err = transform.Transform(tokSchema, record)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(record).To(gm.Equal(expectedRecord))
}

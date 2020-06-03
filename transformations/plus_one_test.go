package transformations

import (
	"testing"

	"github.com/capeprivacy/cape/connector/proto"
	gm "github.com/onsi/gomega"
)

func TestPlusOne(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewPlusOneTransform("income")
	gm.Expect(err).To(gm.BeNil())
	var args Args = nil

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	inputField := &proto.Field{Value: &proto.Field_Double{Double: 64.5}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Double{Double: 65.5}}
	record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: schema}
	expectedRecord := &proto.Record{Fields: []*proto.Field{expectedOutputField}, Schema: schema}

	err = transform.Transform(schema, record)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(record).To(gm.Equal(expectedRecord))
}

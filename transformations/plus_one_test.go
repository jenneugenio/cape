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

	actualOutputField, err := transform.Transform(inputField)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(actualOutputField).To(gm.Equal(expectedOutputField))
}

package transformations

import (
	"github.com/capeprivacy/cape/connector/proto"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestRoundingArgs(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewRoundingTransform("income")
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{})
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{
		"roundingType": "roundToEven",
	})
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{
		"precision": 1,
	})
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{
		"roundingType": "roundToEven",
		"precision":    1,
	})
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{
		"roundingType": "WRONG",
	})
	gm.Expect(err).NotTo(gm.BeNil())

	err = transform.Validate(Args{
		"precision": 1.5,
	})
	gm.Expect(err).NotTo(gm.BeNil())

	err = transform.Validate(Args{
		"precision": -1,
	})
	gm.Expect(err).NotTo(gm.BeNil())
}

func TestRoundingDouble(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewRoundingTransform("income")
	gm.Expect(err).To(gm.BeNil())
	args := Args{
		"roundingType": "roundToEven",
		"precision":    1,
	}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	inputField := &proto.Field{Value: &proto.Field_Double{Double: 64.54}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Double{Double: 64.5}}

	actualOutputField, err := transform.Transform(inputField)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(actualOutputField).To(gm.Equal(expectedOutputField))
}

func TestRoundingFloat(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewRoundingTransform("income")
	gm.Expect(err).To(gm.BeNil())
	var args Args

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	inputField := &proto.Field{Value: &proto.Field_Float{Float: 64.5}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Float{Float: 64}}

	actualOutputField, err := transform.Transform(inputField)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(actualOutputField).To(gm.Equal(expectedOutputField))
}

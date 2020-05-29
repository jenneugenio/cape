package transformations

import (
	"testing"

	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
	gm "github.com/onsi/gomega"
)

func TestPerturbationArgs(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewPerturbationTransform("income")
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{})
	gm.Expect(errors.CausedBy(err, MissingArgument)).To(gm.BeTrue())

	err = transform.Validate(Args{
		"min": 10.0,
	})
	gm.Expect(errors.CausedBy(err, MissingArgument)).To(gm.BeTrue())

	err = transform.Validate(Args{
		"min":  -10.0,
		"max":  10.0,
		"seed": int64(1234),
	})
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{
		"seed": 1234.38,
	})
	gm.Expect(err).NotTo(gm.BeNil())
}

func TestPerturbationInt64(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewPerturbationTransform("income")
	gm.Expect(err).To(gm.BeNil())

	args := Args{
		"min":  -10.,
		"max":  10.,
		"seed": int64(1234),
	}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	inputField := &proto.Field{Value: &proto.Field_Int64{Int64: 100}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Int64{Int64: 94}}

	actualOutputField, err := transform.Transform(inputField)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(actualOutputField).To(gm.Equal(expectedOutputField))
}

func TestPerturbationInt32(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewPerturbationTransform("income")
	gm.Expect(err).To(gm.BeNil())

	args := Args{
		"min":  -10.,
		"max":  10.,
		"seed": int64(3241),
	}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	inputField := &proto.Field{Value: &proto.Field_Int32{Int32: 100}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Int32{Int32: 101}}

	actualOutputField, err := transform.Transform(inputField)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(actualOutputField).To(gm.Equal(expectedOutputField))
}

func TestPerturbationDouble(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewPerturbationTransform("income")
	gm.Expect(err).To(gm.BeNil())

	args := Args{
		"min":  -10.,
		"max":  10.,
		"seed": int64(4354),
	}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	inputField := &proto.Field{Value: &proto.Field_Double{Double: 100}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Double{Double: 93.5}}

	actualOutputField, err := transform.Transform(inputField)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(actualOutputField.GetDouble()).To(gm.BeNumerically("~", expectedOutputField.GetDouble(), 0.1))
}

func TestPerturbationFloat(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewPerturbationTransform("income")
	gm.Expect(err).To(gm.BeNil())

	args := Args{
		"min":  -10.,
		"max":  10.,
		"seed": int64(9876),
	}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	inputField := &proto.Field{Value: &proto.Field_Float{Float: 100}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Float{Float: 107.1}}

	actualOutputField, err := transform.Transform(inputField)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(actualOutputField.GetFloat()).To(gm.BeNumerically("~", expectedOutputField.GetFloat(), 0.1))
}

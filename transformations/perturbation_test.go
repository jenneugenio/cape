package transformations

import (
	"math/rand"
	"testing"

	"github.com/capeprivacy/cape/connector/proto"
	errors "github.com/capeprivacy/cape/partyerrors"
	gm "github.com/onsi/gomega"
)

var schema = &proto.Schema{
	Fields: []*proto.FieldInfo{
		{
			Field: proto.FieldType_BIGINT,
			Name:  "income",
			Size:  8,
		},
	},
}

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
		"min": -10.0,
		"max": 10.0,
	})
	gm.Expect(err).To(gm.BeNil())

	err = transform.Validate(Args{
		"min": "10.0",
	})
	gm.Expect(errors.CausedBy(err, UnsupportedType)).To(gm.BeTrue())
}

func TestPerturbationInt64(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewPerturbationTransform("income")
	gm.Expect(err).To(gm.BeNil())

	args := Args{
		"min": -10.,
		"max": 10.,
	}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	transform.(*PerturbationTransform).sourceSeeded = rand.NewSource(1234)

	inputField := &proto.Field{Value: &proto.Field_Int64{Int64: 100}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Int64{Int64: 94}}
	record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: schema}
	expectedRecord := &proto.Record{Fields: []*proto.Field{expectedOutputField}, Schema: schema}

	err = transform.Transform(schema, record)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(record).To(gm.Equal(expectedRecord))
}

func TestPerturbationInt32(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewPerturbationTransform("income")
	gm.Expect(err).To(gm.BeNil())

	args := Args{
		"min": -10.,
		"max": 10.,
	}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	transform.(*PerturbationTransform).sourceSeeded = rand.NewSource(3241)

	inputField := &proto.Field{Value: &proto.Field_Int32{Int32: 100}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Int32{Int32: 101}}
	record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: schema}
	expectedRecord := &proto.Record{Fields: []*proto.Field{expectedOutputField}, Schema: schema}

	err = transform.Transform(schema, record)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(record).To(gm.Equal(expectedRecord))
}

func TestPerturbationDouble(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewPerturbationTransform("income")
	gm.Expect(err).To(gm.BeNil())

	args := Args{
		"min": -10.,
		"max": 10.,
	}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	transform.(*PerturbationTransform).sourceSeeded = rand.NewSource(4354)

	inputField := &proto.Field{Value: &proto.Field_Double{Double: 100}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Double{Double: 93.5}}
	record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: schema}

	err = transform.Transform(schema, record)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(record.Fields[0].GetDouble()).To(gm.BeNumerically("~", expectedOutputField.GetDouble(), 0.1))
}

func TestPerturbationFloat(t *testing.T) {
	gm.RegisterTestingT(t)

	transform, err := NewPerturbationTransform("income")
	gm.Expect(err).To(gm.BeNil())

	args := Args{
		"min": -10.,
		"max": 10.,
	}

	err = transform.Validate(args)
	gm.Expect(err).To(gm.BeNil())

	err = transform.Initialize(args)
	gm.Expect(err).To(gm.BeNil())

	transform.(*PerturbationTransform).sourceSeeded = rand.NewSource(9876)

	inputField := &proto.Field{Value: &proto.Field_Float{Float: 100}}
	expectedOutputField := &proto.Field{Value: &proto.Field_Float{Float: 107.1}}
	record := &proto.Record{Fields: []*proto.Field{inputField}, Schema: schema}

	err = transform.Transform(schema, record)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(record.Fields[0].GetFloat()).To(gm.BeNumerically("~", expectedOutputField.GetFloat(), 0.1))
}
